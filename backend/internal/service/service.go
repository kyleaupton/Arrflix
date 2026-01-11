package service

import (
	"github.com/kyleaupton/snaggle/backend/internal/config"
	"github.com/kyleaupton/snaggle/backend/internal/logger"
	"github.com/kyleaupton/snaggle/backend/internal/policy"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
)

type Services struct {
	Auth               *AuthService
	Downloaders        *DownloadersService
	DownloadCandidates *DownloadCandidatesService
	DownloadJobs       *DownloadJobsService
	Feed               *FeedService
	Import             *ImportService
	Indexer            *IndexerService
	Libraries          *LibrariesService
	Media              *MediaService
	NameTemplates      *NameTemplatesService
	Policies           *PoliciesService
	Scanner            *ScannerService
	Settings           *SettingsService
	Tmdb               *TmdbService
	Users              *UsersService
}

func New(r *repo.Repository, l *logger.Logger, c *config.Config, opts ...Option) *Services {
	cfg := &cfg{}
	for _, o := range opts {
		o.apply(cfg)
	}

	tmdb := NewTmdbService(r, l)
	indexer := NewIndexerService(r, l, c)
	media := NewMediaService(r, l, tmdb)
	policies := NewPoliciesService(r, l)
	policyEngine := policy.NewEngine(r, l)

	return &Services{
		Auth:               NewAuthService(r, cfg),
		Downloaders:        NewDownloadersService(r),
		DownloadCandidates: NewDownloadCandidatesService(r, l, indexer, media, policyEngine),
		DownloadJobs:       NewDownloadJobsService(r),
		Feed:               NewFeedService(r, l, tmdb),
		Import:             NewImportService(r, l),
		Indexer:            indexer,
		Libraries:          NewLibrariesService(r),
		Media:              media,
		NameTemplates:      NewNameTemplatesService(r),
		Policies:           policies,
		Scanner:            NewScannerService(r, l, tmdb),
		Settings:           NewSettingsService(r),
		Tmdb:               tmdb,
		Users:              NewUsersService(r),
	}
}

type cfg struct {
	jwtSecret string
}

type Option interface{ apply(*cfg) }

type withJWT string

func (w withJWT) apply(c *cfg) { c.jwtSecret = string(w) }

func WithJWTSecret(secret string) Option { return withJWT(secret) }

func coalesce(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}
