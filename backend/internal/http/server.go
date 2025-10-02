package http

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kyleaupton/snaggle/backend/internal/config"
	"github.com/kyleaupton/snaggle/backend/internal/http/handlers"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
	"github.com/kyleaupton/snaggle/backend/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

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

	// Routes
	api := e.Group("/api")
	v1 := api.Group("/v1")

	health.Register(e)
	auth.Register(v1)

	return e
}
