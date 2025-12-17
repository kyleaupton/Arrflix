package http

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kyleaupton/snaggle/backend/internal/config"
	"github.com/kyleaupton/snaggle/backend/internal/downloader"
	"github.com/kyleaupton/snaggle/backend/internal/http/handlers"
	"github.com/kyleaupton/snaggle/backend/internal/http/middlewares"
	"github.com/kyleaupton/snaggle/backend/internal/logger"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
	"github.com/kyleaupton/snaggle/backend/internal/service"

	_ "github.com/kyleaupton/snaggle/backend/internal/http/docs"

	"github.com/labstack/echo/v4"
)

// @title		Snaggle API
// @version		0.0.1
// @BasePath	/api
func NewServer(cfg config.Config, log *logger.Logger, pool *pgxpool.Pool, services *service.Services, repo *repo.Repository, downloaderManager *downloader.Manager) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Handlers
	auth := handlers.NewAuth(cfg, log, pool, services)
	downloadCandidates := handlers.NewDownloadCandidates(services)
	downloaders := handlers.NewDownloaders(services)
	health := handlers.NewHealth()
	indexers := handlers.NewIndexers(services)
	libraries := handlers.NewLibraries(services)
	media := handlers.NewMedia(services)
	nameTemplates := handlers.NewNameTemplates(services)
	policies := handlers.NewPolicies(services)
	rails := handlers.NewRails(services)
	settings := handlers.NewSettings(services)

	api := e.Group("/api")
	v1 := api.Group("/v1")
	protected := v1.Group("", middlewares.JWT(cfg.JWTSecret))

	// Public routes
	auth.RegisterPublic(v1)
	health.RegisterPublic(e)

	// Protected routes
	auth.RegisterProtected(protected)
	downloadCandidates.RegisterProtected(protected)
	downloaders.RegisterProtected(protected)
	indexers.RegisterProtected(protected)
	libraries.RegisterProtected(protected)
	media.RegisterProtected(protected)
	nameTemplates.RegisterProtected(protected)
	policies.RegisterProtected(protected)
	rails.RegisterProtected(protected)
	settings.RegisterProtected(protected)

	// Dev-only routes
	if cfg.Env == "dev" {
		devDownloaderTest := handlers.NewDevDownloaderTest(downloaderManager, repo)
		devDownloaderTest.RegisterDev(e)
	}

	return e
}
