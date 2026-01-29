package service

import (
	"github.com/kyleaupton/arrflix/internal/config"
	"github.com/kyleaupton/arrflix/internal/logger"
	"github.com/kyleaupton/arrflix/internal/policy"
	"github.com/kyleaupton/arrflix/internal/repo"
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
	Setup              *SetupService
	Tmdb               *TmdbService
	UnmatchedFiles     *UnmatchedFilesService
	Users              *UsersService
	Version            *VersionService
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
	users := NewUsersService(r)

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
		Setup:              NewSetupService(r, users),
		Tmdb:               tmdb,
		UnmatchedFiles:     NewUnmatchedFilesService(r, l, tmdb),
		Users:              users,
		Version:            NewVersionService(r, l),
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
