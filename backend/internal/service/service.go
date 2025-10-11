package service

import (
	"github.com/kyleaupton/snaggle/backend/internal/logger"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
)

type Services struct {
	Auth      *AuthService
	Libraries *LibrariesService
	Media     *MediaService
	Scanner   *ScannerService
	Settings  *SettingsService
}

func New(r *repo.Repository, l *logger.Logger, opts ...Option) *Services {
	cfg := &config{}
	for _, o := range opts {
		o.apply(cfg)
	}

	return &Services{
		Auth:      NewAuthService(r, cfg),
		Libraries: NewLibrariesService(r),
		Media:     NewMediaService(r),
		Scanner:   NewScannerService(r, l),
		Settings:  NewSettingsService(r),
	}
}

type config struct {
	jwtSecret string
}

type Option interface{ apply(*config) }

type withJWT string

func (w withJWT) apply(c *config) { c.jwtSecret = string(w) }

func WithJWTSecret(secret string) Option { return withJWT(secret) }
