package services

import (
	"time"

	"github.com/kyleaupton/snaggle/ops/internal/config"
)

// APIService implements the Service interface for the Snaggle API
type APIService struct {
	config *config.Config
}

// NewAPI creates a new API service
func NewAPI(cfg *config.Config) *APIService {
	return &APIService{config: cfg}
}

func (a *APIService) Name() string {
	return "snaggle-api"
}

func (a *APIService) Image() string {
	return "snaggle-api:latest"
}

func (a *APIService) Env() map[string]string {
	return map[string]string{
		"DATABASE_URL":     a.config.DatabaseURL,
		"PORT":             "8080",
		"JWT_SECRET":       a.config.JWTSecret,
		"TMDB_API_KEY":     a.config.TmdbAPIKey,
		"PROWLARR_API_KEY": a.config.ProwlarrAPIKey,
		"SSE_ALLOW_ORIGIN": "*",
	}
}

func (a *APIService) Ports() []PortMapping {
	return []PortMapping{
		{Host: "8080", Container: "8080", Protocol: "tcp"},
	}
}

func (a *APIService) Volumes() []VolumeMount {
	return []VolumeMount{}
}

func (a *APIService) Networks() []string {
	return []string{a.config.NetworkName}
}

func (a *APIService) DependsOn() []string {
	return []string{"snaggle-postgres", "snaggle-prowlarr"}
}

func (a *APIService) HealthCheck() *HealthCheckConfig {
	return &HealthCheckConfig{
		Test:     []string{"CMD-SHELL", "curl -f http://localhost:8080/api/v1/health || exit 1"},
		Interval: 10 * time.Second,
		Timeout:  5 * time.Second,
		Retries:  5,
	}
}

func (a *APIService) Labels() map[string]string {
	return map[string]string{
		"snaggle.managed": "true",
		"snaggle.service": a.Name(),
		"snaggle.type":    "api",
	}
}
