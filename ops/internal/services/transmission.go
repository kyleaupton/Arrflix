package services

import (
	"fmt"
	"time"

	"github.com/kyleaupton/snaggle/ops/internal/config"
)

// TransmissionService implements the Service interface for Transmission
type TransmissionService struct {
	instance *ServiceInstance
	config   *config.Config
}

// NewTransmission creates a new Transmission service
func NewTransmission(cfg *config.Config, instance *ServiceInstance) *TransmissionService {
	return &TransmissionService{
		instance: instance,
		config:   cfg,
	}
}

func (t *TransmissionService) Name() string {
	return fmt.Sprintf("snaggle-transmission-%s", t.instance.Name)
}

func (t *TransmissionService) Image() string {
	return "linuxserver/transmission:latest"
}

func (t *TransmissionService) Env() map[string]string {
	return map[string]string{
		"PUID": "1000",
		"PGID": "1000",
	}
}

func (t *TransmissionService) Ports() []PortMapping {
	return []PortMapping{
		{Host: "9091", Container: "9091", Protocol: "tcp"},
		{Host: "51413", Container: "51413", Protocol: "tcp"},
		{Host: "51413", Container: "51413", Protocol: "udp"},
	}
}

func (t *TransmissionService) Volumes() []VolumeMount {
	volumes := []VolumeMount{}

	// Add config path if specified
	if configPath, ok := t.instance.Config["config_path"].(string); ok {
		volumes = append(volumes, VolumeMount{
			Source: configPath,
			Target: "/config",
			Type:   "bind",
		})
	}

	// Add downloads path if specified
	if downloadsPath, ok := t.instance.Config["downloads_path"].(string); ok {
		volumes = append(volumes, VolumeMount{
			Source: downloadsPath,
			Target: "/downloads",
			Type:   "bind",
		})
	}

	return volumes
}

func (t *TransmissionService) Networks() []string {
	return []string{t.config.NetworkName}
}

func (t *TransmissionService) DependsOn() []string {
	return []string{} // Transmission has no dependencies
}

func (t *TransmissionService) HealthCheck() *HealthCheckConfig {
	return &HealthCheckConfig{
		Test:     []string{"CMD-SHELL", "curl -f http://localhost:9091 || exit 1"},
		Interval: 30 * time.Second,
		Timeout:  10 * time.Second,
		Retries:  3,
	}
}

func (t *TransmissionService) Labels() map[string]string {
	return map[string]string{
		"snaggle.managed":  "true",
		"snaggle.service":  t.Name(),
		"snaggle.type":     "transmission",
		"snaggle.instance": t.instance.Name,
	}
}
