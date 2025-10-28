package services

import (
	"fmt"
	"time"

	"github.com/kyleaupton/snaggle/ops/internal/config"
)

// QBittorrentService implements the Service interface for qBittorrent
type QBittorrentService struct {
	instance *ServiceInstance
	config   *config.Config
}

// NewQBittorrent creates a new qBittorrent service
func NewQBittorrent(cfg *config.Config, instance *ServiceInstance) *QBittorrentService {
	return &QBittorrentService{
		instance: instance,
		config:   cfg,
	}
}

func (q *QBittorrentService) Name() string {
	return fmt.Sprintf("snaggle-qbittorrent-%s", q.instance.Name)
}

func (q *QBittorrentService) Image() string {
	return "linuxserver/qbittorrent:latest"
}

func (q *QBittorrentService) Env() map[string]string {
	env := map[string]string{
		"PUID": "1000",
		"PGID": "1000",
	}

	// Add webui port if specified
	if port, ok := q.instance.Config["webui_port"].(string); ok {
		env["WEBUI_PORT"] = port
	}

	return env
}

func (q *QBittorrentService) Ports() []PortMapping {
	ports := []PortMapping{
		{Host: "6881", Container: "6881", Protocol: "tcp"},
		{Host: "6881", Container: "6881", Protocol: "udp"},
	}

	// Add webui port if specified
	if port, ok := q.instance.Config["webui_port"].(string); ok {
		ports = append(ports, PortMapping{
			Host:      port,
			Container: port,
			Protocol:  "tcp",
		})
	}

	return ports
}

func (q *QBittorrentService) Volumes() []VolumeMount {
	volumes := []VolumeMount{}

	// Add config path if specified
	if configPath, ok := q.instance.Config["config_path"].(string); ok {
		volumes = append(volumes, VolumeMount{
			Source: configPath,
			Target: "/config",
			Type:   "bind",
		})
	}

	// Add downloads path if specified
	if downloadsPath, ok := q.instance.Config["downloads_path"].(string); ok {
		volumes = append(volumes, VolumeMount{
			Source: downloadsPath,
			Target: "/downloads",
			Type:   "bind",
		})
	}

	return volumes
}

func (q *QBittorrentService) Networks() []string {
	return []string{q.config.NetworkName}
}

func (q *QBittorrentService) DependsOn() []string {
	return []string{} // qBittorrent has no dependencies
}

func (q *QBittorrentService) HealthCheck() *HealthCheckConfig {
	webuiPort := "8080" // Default port
	if port, ok := q.instance.Config["webui_port"].(string); ok {
		webuiPort = port
	}

	return &HealthCheckConfig{
		Test:     []string{"CMD-SHELL", fmt.Sprintf("curl -f http://localhost:%s || exit 1", webuiPort)},
		Interval: 30 * time.Second,
		Timeout:  10 * time.Second,
		Retries:  3,
	}
}

func (q *QBittorrentService) Labels() map[string]string {
	return map[string]string{
		"snaggle.managed":  "true",
		"snaggle.service":  q.Name(),
		"snaggle.type":     "qbittorrent",
		"snaggle.instance": q.instance.Name,
	}
}
