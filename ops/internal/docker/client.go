package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/kyleaupton/snaggle/ops/internal/services"
)

// Client wraps the Docker client with helper methods
type Client struct {
	cli *client.Client
}

// ErrContainerNotFound is returned when a container is not found
var ErrContainerNotFound = errors.New("container not found")

// NewClient creates a new Docker client wrapper
func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}
	return &Client{cli: cli}, nil
}

// ContainerStatus represents the current state of a container
type ContainerStatus struct {
	ID      string
	Name    string
	Status  string
	Image   string
	Labels  map[string]string
	Created time.Time
}

// EnsureNetwork creates a network if it doesn't exist
func (c *Client) EnsureNetwork(ctx context.Context, name string) error {
	_, err := c.cli.NetworkInspect(ctx, name, types.NetworkInspectOptions{})
	if err != nil {
		if client.IsErrNotFound(err) {
			_, err = c.cli.NetworkCreate(ctx, name, types.NetworkCreate{
				Driver: "bridge",
				Labels: map[string]string{
					"snaggle.managed": "true",
				},
			})
			if err != nil {
				return fmt.Errorf("failed to create network %s: %w", name, err)
			}
		} else {
			return fmt.Errorf("failed to inspect network %s: %w", name, err)
		}
	}
	return nil
}

// EnsureVolume creates a volume if it doesn't exist
func (c *Client) EnsureVolume(ctx context.Context, name string) error {
	_, err := c.cli.VolumeInspect(ctx, name)
	if err != nil {
		if client.IsErrNotFound(err) {
			_, err = c.cli.VolumeCreate(ctx, volume.CreateOptions{
				Name: name,
				Labels: map[string]string{
					"snaggle.managed": "true",
				},
			})
			if err != nil {
				return fmt.Errorf("failed to create volume %s: %w", name, err)
			}
		} else {
			return fmt.Errorf("failed to inspect volume %s: %w", name, err)
		}
	}
	return nil
}

// EnsureContainer creates and starts a container if it doesn't exist or needs updating
func (c *Client) EnsureContainer(ctx context.Context, svc services.Service) error {
	// Check if container exists
	existing, err := c.GetContainerStatus(ctx, svc.Name())
	if err != nil && !errors.Is(err, ErrContainerNotFound) {
		return fmt.Errorf("failed to check container status: %w", err)
	}

	// If container exists, check if it needs recreation
	if existing != nil {
		// For now, always recreate if config hash differs
		// TODO: Implement proper config hash comparison
		if err := c.StopContainer(ctx, svc.Name()); err != nil {
			return fmt.Errorf("failed to stop existing container: %w", err)
		}
		if err := c.RemoveContainer(ctx, svc.Name()); err != nil {
			return fmt.Errorf("failed to remove existing container: %w", err)
		}
	}

	// Handle image preparation (build or pull)
	if buildInfo := svc.BuildInfo(); buildInfo != nil {
		// Build image locally
		log.Printf("Building image %s for container %s", svc.Image(), svc.Name())
		if err := c.BuildImage(ctx, svc.Image(), buildInfo.Dockerfile, buildInfo.Context); err != nil {
			log.Printf("Failed to build image %s: %v", svc.Image(), err)
			return fmt.Errorf("failed to build image %s: %w", svc.Image(), err)
		}
	} else {
		// Pull image from registry
		log.Printf("Pulling image %s for container %s", svc.Image(), svc.Name())
		if err := c.pullImageIfNeeded(ctx, svc.Image()); err != nil {
			log.Printf("Failed to pull image %s: %v", svc.Image(), err)
			return fmt.Errorf("failed to pull image %s: %w", svc.Image(), err)
		}
	}

	// Create container
	log.Printf("Creating container %s with image %s", svc.Name(), svc.Image())
	containerConfig := &container.Config{
		Image:  svc.Image(),
		Env:    envMapToSlice(svc.Env()),
		Labels: svc.Labels(),
	}

	// Add command if service provides one
	if cmdSvc, ok := svc.(interface{ Command() []string }); ok {
		containerConfig.Cmd = cmdSvc.Command()
	}

	// Add health check if defined
	if hc := svc.HealthCheck(); hc != nil {
		containerConfig.Healthcheck = &container.HealthConfig{
			Test:     hc.Test,
			Interval: hc.Interval,
			Timeout:  hc.Timeout,
			Retries:  hc.Retries,
		}
	}

	// Create host config
	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
		PortBindings: make(map[nat.Port][]nat.PortBinding),
	}

	// Add port bindings
	for _, port := range svc.Ports() {
		portKey := nat.Port(port.Container + "/" + port.Protocol)
		hostConfig.PortBindings[portKey] = []nat.PortBinding{
			{HostIP: "0.0.0.0", HostPort: port.Host},
		}
	}

	// Add volume mounts
	for _, vol := range svc.Volumes() {
		if vol.Type == "volume" {
			// For named volumes, just use source:target format
			hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%s:%s", vol.Source, vol.Target))
		} else if vol.Type == "bind" {
			// For bind mounts, use source:target format (Docker auto-detects bind mounts)
			hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%s:%s", vol.Source, vol.Target))
		} else {
			// For other types (tmpfs, etc.), use source:target:mode format
			hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%s:%s:%s", vol.Source, vol.Target, vol.Type))
		}
	}

	// Create network config
	networkingConfig := &network.NetworkingConfig{
		EndpointsConfig: make(map[string]*network.EndpointSettings),
	}
	for _, netName := range svc.Networks() {
		networkingConfig.EndpointsConfig[netName] = &network.EndpointSettings{}
	}

	// Create container
	resp, err := c.cli.ContainerCreate(ctx, containerConfig, hostConfig, networkingConfig, nil, svc.Name())
	if err != nil {
		log.Printf("Failed to create container %s: %v", svc.Name(), err)
		return fmt.Errorf("failed to create container: %w", err)
	}
	log.Printf("Successfully created container %s with ID %s", svc.Name(), resp.ID)

	// Start container
	if err := c.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Printf("Failed to start container %s: %v", svc.Name(), err)
		return fmt.Errorf("failed to start container: %w", err)
	}
	log.Printf("Successfully started container %s", svc.Name())

	return nil
}

// pullImageIfNeeded pulls an image if it doesn't exist locally
func (c *Client) pullImageIfNeeded(ctx context.Context, imageName string) error {
	// Try to pull the image (Docker will skip if already exists)
	log.Printf("Pulling image %s from registry", imageName)
	reader, err := c.cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer reader.Close()

	// Read the response to completion (this ensures the pull completes)
	buf := make([]byte, 1024)
	for {
		_, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading pull response for %s: %v", imageName, err)
			break
		}
	}

	log.Printf("Successfully pulled image %s", imageName)
	return nil
}

// GetContainerStatus returns the status of a container by name
func (c *Client) GetContainerStatus(ctx context.Context, name string) (*ContainerStatus, error) {
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", name)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	if len(containers) == 0 {
		return nil, ErrContainerNotFound
	}

	container := containers[0]
	return &ContainerStatus{
		ID:      container.ID,
		Name:    container.Names[0][1:], // Remove leading slash
		Status:  container.Status,
		Image:   container.Image,
		Labels:  container.Labels,
		Created: time.Unix(container.Created, 0),
	}, nil
}

// StopContainer stops a running container
func (c *Client) StopContainer(ctx context.Context, name string) error {
	timeout := int(30)
	return c.cli.ContainerStop(ctx, name, container.StopOptions{
		Timeout: &timeout,
	})
}

// RemoveContainer removes a container
func (c *Client) RemoveContainer(ctx context.Context, name string) error {
	return c.cli.ContainerRemove(ctx, name, container.RemoveOptions{
		Force: true,
	})
}

// ListManagedContainers returns all containers managed by Snaggle
func (c *Client) ListManagedContainers(ctx context.Context) ([]*ContainerStatus, error) {
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("label", "snaggle.managed=true")),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list managed containers: %w", err)
	}

	var statuses []*ContainerStatus
	for _, container := range containers {
		statuses = append(statuses, &ContainerStatus{
			ID:      container.ID,
			Name:    container.Names[0][1:], // Remove leading slash
			Status:  container.Status,
			Image:   container.Image,
			Labels:  container.Labels,
			Created: time.Unix(container.Created, 0),
		})
	}

	return statuses, nil
}

// BuildImage builds an image from a Dockerfile
func (c *Client) BuildImage(ctx context.Context, imageName, dockerfilePath, buildContext string) error {
	log.Printf("Building image %s from %s", imageName, dockerfilePath)

	// Create tar archive from build context
	tarReader, err := c.createTarArchive(buildContext)
	if err != nil {
		return fmt.Errorf("failed to create tar archive: %w", err)
	}
	defer tarReader.Close()

	buildOptions := types.ImageBuildOptions{
		Tags:       []string{imageName},
		Dockerfile: dockerfilePath,
		Remove:     true,
	}

	buildResponse, err := c.cli.ImageBuild(ctx, tarReader, buildOptions)
	if err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}
	defer buildResponse.Body.Close()

	// Read the build output
	buf := make([]byte, 1024)
	for {
		_, err := buildResponse.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading build response: %v", err)
			break
		}
	}

	log.Printf("Successfully built image %s", imageName)
	return nil
}

// createTarArchive creates a tar archive from the build context directory
func (c *Client) createTarArchive(buildContext string) (io.ReadCloser, error) {
	var buf bytes.Buffer
	tarWriter := tar.NewWriter(&buf)

	err := filepath.Walk(buildContext, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories
		if filepath.Base(path)[0] == '.' {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		// Make path relative to build context
		relPath, err := filepath.Rel(buildContext, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// Write file content if it's a regular file
		if info.Mode().IsRegular() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := tarWriter.Close(); err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

// envMapToSlice converts a map of environment variables to a slice
func envMapToSlice(env map[string]string) []string {
	var result []string
	for key, value := range env {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}
	return result
}
