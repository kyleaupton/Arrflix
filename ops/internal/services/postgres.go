package services

import (
	"time"

	"github.com/kyleaupton/snaggle/ops/internal/config"
)

// PostgresService implements the Service interface for PostgreSQL
type PostgresService struct {
	config *config.Config
}

// NewPostgres creates a new PostgreSQL service
func NewPostgres(cfg *config.Config) *PostgresService {
	return &PostgresService{config: cfg}
}

func (p *PostgresService) Name() string {
	return "snaggle-postgres"
}

func (p *PostgresService) Image() string {
	return "snaggle-postgres:latest"
}

func (p *PostgresService) Env() map[string]string {
	return map[string]string{
		"POSTGRES_DB":       p.config.PostgresDB,
		"POSTGRES_USER":     p.config.PostgresUser,
		"POSTGRES_PASSWORD": p.config.PostgresPassword,
		"PGDATA":            "/var/lib/postgresql/data",
	}
}

func (p *PostgresService) Ports() []PortMapping {
	return []PortMapping{
		{Host: "5432", Container: "5432", Protocol: "tcp"},
	}
}

func (p *PostgresService) Volumes() []VolumeMount {
	return []VolumeMount{
		{Source: "snaggle_pg_data", Target: "/var/lib/postgresql/data", Type: "volume"},
	}
}

func (p *PostgresService) Networks() []string {
	return []string{p.config.NetworkName}
}

func (p *PostgresService) DependsOn() []string {
	return []string{} // PostgreSQL has no dependencies
}

func (p *PostgresService) HealthCheck() *HealthCheckConfig {
	return &HealthCheckConfig{
		Test:     []string{"CMD-SHELL", "pg_isready -U " + p.config.PostgresUser},
		Interval: 10 * time.Second,
		Timeout:  5 * time.Second,
		Retries:  5,
	}
}

func (p *PostgresService) Labels() map[string]string {
	return map[string]string{
		"snaggle.managed": "true",
		"snaggle.service": p.Name(),
		"snaggle.type":    "postgres",
	}
}

func (p *PostgresService) BuildInfo() *BuildInfo {
	if p.config.RuntimeMode == "dev" {
		return &BuildInfo{
			Dockerfile: "ops/images/Dockerfile.postgres",
			Context:    "/host", // Build context is the mounted host directory
		}
	}
	return nil // Production images come from registry
}
