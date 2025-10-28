package services

import (
	"fmt"

	"github.com/kyleaupton/snaggle/ops/internal/config"
)

// NewDynamicService creates a dynamic service from a ServiceInstance
func NewDynamicService(cfg *config.Config, instance *ServiceInstance) (Service, error) {
	switch instance.Type {
	case "qbittorrent":
		return NewQBittorrent(cfg, instance), nil
	case "transmission":
		return NewTransmission(cfg, instance), nil
	default:
		return nil, fmt.Errorf("unsupported service type: %s", instance.Type)
	}
}
