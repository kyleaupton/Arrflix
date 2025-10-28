package services

import (
	"time"

	"github.com/kyleaupton/snaggle/ops/internal/config"
)

// NginxDevService implements the Service interface for Nginx in dev mode
type NginxDevService struct {
	config *config.Config
}

// NewNginxDev creates a new Nginx dev service
func NewNginxDev(cfg *config.Config) *NginxDevService {
	return &NginxDevService{config: cfg}
}

func (n *NginxDevService) Name() string {
	return "snaggle-nginx-dev"
}

func (n *NginxDevService) Image() string {
	return "nginx:alpine"
}

func (n *NginxDevService) Env() map[string]string {
	return map[string]string{}
}

func (n *NginxDevService) Ports() []PortMapping {
	return []PortMapping{
		{Host: "8484", Container: "80", Protocol: "tcp"},
	}
}

func (n *NginxDevService) Volumes() []VolumeMount {
	return []VolumeMount{}
}

func (n *NginxDevService) Networks() []string {
	return []string{n.config.NetworkName}
}

func (n *NginxDevService) DependsOn() []string {
	return []string{"snaggle-api-dev", "snaggle-web-dev"} // Nginx depends on both API and web dev servers
}

func (n *NginxDevService) HealthCheck() *HealthCheckConfig {
	return &HealthCheckConfig{
		Test:     []string{"CMD-SHELL", "curl -f http://localhost:80 || exit 1"},
		Interval: 10 * time.Second,
		Timeout:  5 * time.Second,
		Retries:  5,
	}
}

func (n *NginxDevService) Labels() map[string]string {
	return map[string]string{
		"snaggle.managed": "true",
		"snaggle.service": n.Name(),
		"snaggle.type":    "nginx-dev",
	}
}
