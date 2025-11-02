package services

import (
	"time"

	"github.com/kyleaupton/snaggle/ops/internal/config"
)

// WebDevService implements the Service interface for the Vite dev server
type WebDevService struct {
	config *config.Config
}

// NewWebDev creates a new web dev service
func NewWebDev(cfg *config.Config) *WebDevService {
	return &WebDevService{config: cfg}
}

func (w *WebDevService) Name() string {
	return "snaggle-web-dev"
}

func (w *WebDevService) Image() string {
	return "node:20"
}

func (w *WebDevService) Env() map[string]string {
	return map[string]string{
		"NODE_ENV": "development",
	}
}

func (w *WebDevService) Ports() []PortMapping {
	return []PortMapping{
		{Host: "5173", Container: "5173", Protocol: "tcp"},
	}
}

func (w *WebDevService) Volumes() []VolumeMount {
	return []VolumeMount{
		{Source: "/host/web", Target: "/web", Type: "bind"},
		{Source: "web_node_modules", Target: "/web/node_modules", Type: "volume"},
	}
}

func (w *WebDevService) Networks() []string {
	return []string{w.config.NetworkName}
}

func (w *WebDevService) DependsOn() []string {
	return []string{} // Vite dev server has no dependencies
}

func (w *WebDevService) HealthCheck() *HealthCheckConfig {
	return &HealthCheckConfig{
		Test:     []string{"CMD-SHELL", "curl -f http://localhost:5173 || exit 1"},
		Interval: 10 * time.Second,
		Timeout:  5 * time.Second,
		Retries:  5,
	}
}

func (w *WebDevService) Labels() map[string]string {
	return map[string]string{
		"snaggle.managed": "true",
		"snaggle.service": w.Name(),
		"snaggle.type":    "web-dev",
	}
}

func (w *WebDevService) Command() []string {
	return []string{"sh", "-c", "cd /web && npm run dev -- --host 0.0.0.0 --port 5173"}
}

func (w *WebDevService) BuildInfo() *BuildInfo {
	return &BuildInfo{
		Dockerfile: "ops/images/Dockerfile.web-dev",
		Context:    "/host", // Build context is the mounted host directory
	}
}
