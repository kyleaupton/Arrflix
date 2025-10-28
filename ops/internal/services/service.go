package services

import (
	"time"
)

// Service defines the interface for all managed services
type Service interface {
	Name() string
	Image() string
	Env() map[string]string
	Ports() []PortMapping
	Volumes() []VolumeMount
	Networks() []string
	DependsOn() []string
	HealthCheck() *HealthCheckConfig
	Labels() map[string]string
}

// PortMapping defines port mappings for containers
type PortMapping struct {
	Host      string
	Container string
	Protocol  string // tcp/udp
}

// VolumeMount defines volume mounts for containers
type VolumeMount struct {
	Source string // volume name or host path
	Target string // container path
	Type   string // volume, bind, tmpfs
}

// HealthCheckConfig defines health check configuration
type HealthCheckConfig struct {
	Test     []string
	Interval time.Duration
	Timeout  time.Duration
	Retries  int
}

// ServiceInstance represents a dynamic service from the database
type ServiceInstance struct {
	ID      string
	Name    string
	Type    string
	Enabled bool
	Config  map[string]interface{} // parsed from JSONB
}
