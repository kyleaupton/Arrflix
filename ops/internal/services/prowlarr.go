package services

import (
	"time"

	"github.com/kyleaupton/snaggle/ops/internal/config"
)

// ProwlarrService implements the Service interface for Prowlarr
type ProwlarrService struct {
	config *config.Config
}

// NewProwlarr creates a new Prowlarr service
func NewProwlarr(cfg *config.Config) *ProwlarrService {
	return &ProwlarrService{config: cfg}
}

func (p *ProwlarrService) Name() string {
	return "snaggle-prowlarr"
}

func (p *ProwlarrService) Image() string {
	return "snaggle-prowlarr:latest"
}

func (p *ProwlarrService) Env() map[string]string {
	return map[string]string{
		"PUID": "1000",
		"PGID": "1000",
	}
}

func (p *ProwlarrService) Ports() []PortMapping {
	return []PortMapping{
		{Host: "9696", Container: "9696", Protocol: "tcp"},
	}
}

func (p *ProwlarrService) Volumes() []VolumeMount {
	return []VolumeMount{
		{Source: "snaggle_prowlarr_data", Target: "/var/lib/prowlarr", Type: "volume"},
	}
}

func (p *ProwlarrService) Networks() []string {
	return []string{p.config.NetworkName}
}

func (p *ProwlarrService) DependsOn() []string {
	return []string{} // Prowlarr has no dependencies
}

func (p *ProwlarrService) HealthCheck() *HealthCheckConfig {
	return &HealthCheckConfig{
		Test:     []string{"CMD-SHELL", "curl -f http://localhost:9696/api/v1/system/status || exit 1"},
		Interval: 10 * time.Second,
		Timeout:  5 * time.Second,
		Retries:  5,
	}
}

func (p *ProwlarrService) Labels() map[string]string {
	return map[string]string{
		"snaggle.managed": "true",
		"snaggle.service": p.Name(),
		"snaggle.type":    "prowlarr",
	}
}
