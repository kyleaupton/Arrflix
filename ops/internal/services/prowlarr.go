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
	return "linuxserver/prowlarr:latest"
}

func (p *ProwlarrService) Env() map[string]string {
	return map[string]string{
		"PUID": "1000",
		"PGID": "1000",
		// PostgreSQL database configuration
		"PROWLARR__POSTGRES__HOST":     "snaggle-postgres",
		"PROWLARR__POSTGRES__PORT":     "5432",
		"PROWLARR__POSTGRES__USER":     "prowlarr",
		"PROWLARR__POSTGRES__PASSWORD": "prowlarrpw",
		"PROWLARR__POSTGRES__DATABASE": "prowlarr",
		"PROWLARR__AUTH__APIKEY":       p.config.ProwlarrAPIKey,
		"PROWLARR__SERVER__PORT":       "9696",
		"PROWLARR__AUTH__METHOD":       "External",
		// "PROWLARR__AUTH__ENABLED":      "false",
		// "PROWLARR__AUTH__REQUIRED":     "false",
	}
}

func (p *ProwlarrService) Ports() []PortMapping {
	return []PortMapping{
		{Host: "9697", Container: "9696", Protocol: "tcp"}, // testing port
	} // No host port mapping - internal network only
}

func (p *ProwlarrService) Volumes() []VolumeMount {
	return []VolumeMount{
		// {Source: "snaggle_prowlarr_data", Target: "/var/lib/prowlarr", Type: "volume"},
	}
}

func (p *ProwlarrService) Networks() []string {
	return []string{p.config.NetworkName}
}

func (p *ProwlarrService) DependsOn() []string {
	return []string{"snaggle-postgres"} // Prowlarr depends on PostgreSQL
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

func (p *ProwlarrService) BuildInfo() *BuildInfo {
	// Prowlarr uses official image from registry
	return nil
}
