package http

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kyleaupton/snaggle/backend/internal/config"
	"github.com/kyleaupton/snaggle/backend/internal/http/handlers"

	// "github.com/kyleaupton/snaggle/backend/internal/http/handlers"
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

	// Handlers
	health := handlers.NewHealth()
	// queue := handlers.NewQueue(log, jr)
	// media := handlers.NewMedia(log, pool)

	// Routes
	// api := e.Group("/api")
	// v1 := api.Group("/v1")

	e.GET("/health", health.Health)
	// v1.GET("/queue", queue.List)
	// v1.POST("/request", queue.Enqueue)
	// v1.GET("/events", queue.Events) // SSE (optional)
	// v1.GET("/media/:id", media.Get)
	// v1.POST("/media", media.Create)

	return e
}
