package http

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kyleaupton/snaggle/backend/internal/config"
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
func NewServer(cfg config.Config, log *logger.Logger, pool *pgxpool.Pool, services *service.Services, repo *repo.Repository) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Handlers
	auth := handlers.NewAuth(cfg, log, pool, services)
	health := handlers.NewHealth()
	indexers := handlers.NewIndexers(services)
	libraries := handlers.NewLibraries(services)
	media := handlers.NewMedia(services)
	nameTemplates := handlers.NewNameTemplates(services)
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
	indexers.RegisterProtected(protected)
	libraries.RegisterProtected(protected)
	media.RegisterProtected(protected)
	nameTemplates.RegisterProtected(protected)
	rails.RegisterProtected(protected)
	settings.RegisterProtected(protected)

	return e
}
