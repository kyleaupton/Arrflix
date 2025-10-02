package http

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kyleaupton/snaggle/backend/internal/config"
	"github.com/kyleaupton/snaggle/backend/internal/http/handlers"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
	"github.com/kyleaupton/snaggle/backend/internal/service"

	// "github.com/kyleaupton/snaggle/backend/internal/jobs"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func NewServer(cfg config.Config, log zerolog.Logger, pool *pgxpool.Pool) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Global middleware
	// e.Use(middleware.Recover())
	// e.Use(middleware.RequestID())
	// e.Use(middleware.Gzip())
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{cfg.CORSOrigin, "*"},
	// 	AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE, echo.OPTIONS},
	// }))

    // Layers
    r := repo.New(pool)
    services := service.New(r, service.WithJWTSecret(cfg.JWTSecret))

    // Handlers
    health := handlers.NewHealth()
    auth := handlers.NewAuth(cfg, log, pool, services)

    // Routes
    api := e.Group("/api")
    v1 := api.Group("/v1")

    e.GET("/health", health.Health)
    v1.POST("/auth/login", auth.Login)
    v1.GET("/auth/me", auth.Me)

	return e
}
