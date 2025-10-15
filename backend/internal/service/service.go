package service

import (
	"github.com/kyleaupton/snaggle/backend/internal/config"
	"github.com/kyleaupton/snaggle/backend/internal/logger"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
)

type Services struct {
	Auth      *AuthService
	Indexer   *IndexerService
	Libraries *LibrariesService
	Media     *MediaService
	Rails     *RailsService
	Scanner   *ScannerService
	Settings  *SettingsService
	Tmdb      *TmdbService
}

func New(r *repo.Repository, l *logger.Logger, c *config.Config, opts ...Option) *Services {
	cfg := &cfg{}
	for _, o := range opts {
		o.apply(cfg)
	}

	tmdb := NewTmdbService(r, l)

	return &Services{
		Auth:      NewAuthService(r, cfg),
		Indexer:   NewIndexerService(r, l, c),
		Libraries: NewLibrariesService(r),
		Media:     NewMediaService(r, l, tmdb),
		Rails:     NewRailsService(r, tmdb),
		Scanner:   NewScannerService(r, l, tmdb),
		Settings:  NewSettingsService(r),
		Tmdb:      tmdb,
	}
}

type cfg struct {
	jwtSecret string
}

type Option interface{ apply(*cfg) }

type withJWT string

func (w withJWT) apply(c *cfg) { c.jwtSecret = string(w) }

func WithJWTSecret(secret string) Option { return withJWT(secret) }
