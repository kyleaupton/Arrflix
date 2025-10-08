package http

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kyleaupton/snaggle/backend/internal/config"
	"github.com/kyleaupton/snaggle/backend/internal/http/handlers"
	"github.com/kyleaupton/snaggle/backend/internal/http/middlewares"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
	"github.com/kyleaupton/snaggle/backend/internal/service"

	_ "github.com/kyleaupton/snaggle/backend/internal/http/docs"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// @title			Snaggle API
// @version		0.0.1
// @BasePath	/api
func NewServer(cfg config.Config, log zerolog.Logger, pool *pgxpool.Pool) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Layers
	r := repo.New(pool)
	services := service.New(r, service.WithJWTSecret(cfg.JWTSecret))

	// Handlers
	health := handlers.NewHealth()
	auth := handlers.NewAuth(cfg, log, pool, services)
	settings := handlers.NewSettings(services)
	libraries := handlers.NewLibraries(services)
	media := handlers.NewMedia(services)

	// Routes
	api := e.Group("/api")
	v1 := api.Group("/v1")
	protected := v1.Group("", middlewares.JWT(cfg.JWTSecret))

	v1.GET("/swagger/*", echoSwagger.WrapHandler)

	health.Register(e)
	auth.RegisterPublic(v1)
	auth.RegisterProtected(protected)
	settings.RegisterProtected(protected)
	libraries.RegisterProtected(protected)
	media.RegisterProtected(protected)

	return e
}
