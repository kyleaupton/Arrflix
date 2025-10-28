package services

import (
	"time"

	"github.com/kyleaupton/snaggle/ops/internal/config"
)

// APIDevService implements the Service interface for the Snaggle API in dev mode
type APIDevService struct {
	config *config.Config
}

// NewAPIDev creates a new API dev service
func NewAPIDev(cfg *config.Config) *APIDevService {
	return &APIDevService{config: cfg}
}

func (a *APIDevService) Name() string {
	return "snaggle-api-dev"
}

func (a *APIDevService) Image() string {
	return "golang:1.24"
}

func (a *APIDevService) Env() map[string]string {
	return map[string]string{
		"DATABASE_URL":     a.config.DatabaseURL,
		"PORT":             "8080",
		"JWT_SECRET":       a.config.JWTSecret,
		"TMDB_API_KEY":     a.config.TmdbAPIKey,
		"PROWLARR_API_KEY": a.config.ProwlarrAPIKey,
		"SSE_ALLOW_ORIGIN": "*",
	}
}

func (a *APIDevService) Ports() []PortMapping {
	return []PortMapping{
		{Host: "8080", Container: "8080", Protocol: "tcp"},
	}
}

func (a *APIDevService) Volumes() []VolumeMount {
	return []VolumeMount{
		{Source: "/host/backend", Target: "/app/backend", Type: "bind"},
	}
}

func (a *APIDevService) Networks() []string {
	return []string{a.config.NetworkName}
}

func (a *APIDevService) DependsOn() []string {
	return []string{"snaggle-postgres", "snaggle-prowlarr"}
}

func (a *APIDevService) HealthCheck() *HealthCheckConfig {
	return &HealthCheckConfig{
		Test:     []string{"CMD-SHELL", "curl -f http://localhost:8080/api/v1/health || exit 1"},
		Interval: 10 * time.Second,
		Timeout:  5 * time.Second,
		Retries:  5,
	}
}

func (a *APIDevService) Labels() map[string]string {
	return map[string]string{
		"snaggle.managed": "true",
		"snaggle.service": a.Name(),
		"snaggle.type":    "api-dev",
	}
}

func (a *APIDevService) Command() []string {
	return []string{"sh", "-c", "cd /app/backend && go run ./cmd/api"}
}
