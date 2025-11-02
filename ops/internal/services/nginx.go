package services

import (
	"time"

	"github.com/kyleaupton/snaggle/ops/internal/config"
)

// NginxService implements the Service interface for Nginx
type NginxService struct {
	config *config.Config
}

// NewNginx creates a new Nginx service
func NewNginx(cfg *config.Config) *NginxService {
	return &NginxService{config: cfg}
}

func (n *NginxService) Name() string {
	return "snaggle-nginx"
}

func (n *NginxService) Image() string {
	return "snaggle-nginx:latest"
}

func (n *NginxService) Env() map[string]string {
	return map[string]string{}
}

func (n *NginxService) Ports() []PortMapping {
	return []PortMapping{
		{Host: "8484", Container: "80", Protocol: "tcp"},
	}
}

func (n *NginxService) Volumes() []VolumeMount {
	return []VolumeMount{}
}

func (n *NginxService) Networks() []string {
	return []string{n.config.NetworkName}
}

func (n *NginxService) DependsOn() []string {
	return []string{"snaggle-api"} // Nginx depends on API
}

func (n *NginxService) HealthCheck() *HealthCheckConfig {
	return &HealthCheckConfig{
		Test:     []string{"CMD-SHELL", "curl -f http://localhost:80 || exit 1"},
		Interval: 10 * time.Second,
		Timeout:  5 * time.Second,
		Retries:  5,
	}
}

func (n *NginxService) Labels() map[string]string {
	return map[string]string{
		"snaggle.managed": "true",
		"snaggle.service": n.Name(),
		"snaggle.type":    "nginx",
	}
}

func (n *NginxService) BuildInfo() *BuildInfo {
	// Production nginx images come from registry
	return nil
}
