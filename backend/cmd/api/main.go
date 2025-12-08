package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/kyleaupton/snaggle/backend/internal/config"
	"github.com/kyleaupton/snaggle/backend/internal/db"
	"github.com/kyleaupton/snaggle/backend/internal/http"
	"github.com/kyleaupton/snaggle/backend/internal/logger"
	"github.com/kyleaupton/snaggle/backend/internal/model"
	"github.com/kyleaupton/snaggle/backend/internal/policy"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
	"github.com/kyleaupton/snaggle/backend/internal/service"
)

func main() {
	// Logger
	logg := logger.New(true)

	// Load config
	cfg := config.Load(logg)

	// DB
	pool, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		logg.Fatal().Err(err).Msg("open db")
	}
	defer pool.Close()

	// Migrations (run on startup; idempotent, using embedded files)
	if err := db.ApplyMigrations(cfg.DatabaseURL); err != nil {
		logg.Fatal().Err(err).Msg("migrate")
	}

	// Repo
	repo := repo.New(pool)

	// Services
	services := service.New(repo, logg, &cfg, service.WithJWTSecret(cfg.JWTSecret))

	// testing
	engine := policy.NewEngine(repo)

	plan, err := engine.Evaluate(context.Background(), "https://example.com/torrent.torrent", model.TorrentMetadata{
		Size:       1000,
		Seeders:    10,
		Peers:      10,
		Title:      "Test Torrent",
		Tracker:    "Test Tracker",
		TrackerID:  "1234567890",
		Categories: []string{"Test Category"},
	})
	if err != nil {
		logg.Fatal().Err(err).Msg("evaluate plan")
	}

	logg.Info().Msgf("plan: %+v", plan)

	// HTTP
	e := http.NewServer(cfg, logg, pool, services, repo)
	go func() {
		logg.Info().Str("port", cfg.Port).Msg("http listen")
		if err := e.Start(":" + cfg.Port); err != nil {
			log.Println("server stopped:", err)
		}
	}()

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop()

	shCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = e.Shutdown(shCtx)
	logg.Info().Msg("bye")
}
