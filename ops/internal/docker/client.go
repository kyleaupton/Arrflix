package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
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
	if err != nil && !client.IsErrNotFound(err) {
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

	// Create container
	containerConfig := &container.Config{
		Image:  svc.Image(),
		Env:    envMapToSlice(svc.Env()),
		Labels: svc.Labels(),
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
		hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%s:%s:%s", vol.Source, vol.Target, vol.Type))
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
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Start container
	if err := c.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

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
		return nil, fmt.Errorf("container not found")
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

// envMapToSlice converts a map of environment variables to a slice
func envMapToSlice(env map[string]string) []string {
	var result []string
	for key, value := range env {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}
	return result
}
