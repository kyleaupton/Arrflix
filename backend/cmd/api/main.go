package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/kyleaupton/Arrflix/internal/config"
	"github.com/kyleaupton/Arrflix/internal/db"
	"github.com/kyleaupton/Arrflix/internal/downloader"
	"github.com/kyleaupton/Arrflix/internal/downloader/qbittorrent"
	"github.com/kyleaupton/Arrflix/internal/http"
	"github.com/kyleaupton/Arrflix/internal/jobs/downloadjobs"
	"github.com/kyleaupton/Arrflix/internal/logger"
	"github.com/kyleaupton/Arrflix/internal/repo"
	"github.com/kyleaupton/Arrflix/internal/service"
	"github.com/kyleaupton/Arrflix/internal/sse"
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

	// In-process SSE broker
	broker := sse.NewBroker()

	// Downloader Manager
	downloaderRegistry := downloader.NewRegistry()
	qbittorrent.Register(downloaderRegistry)
	downloaderManager := downloader.NewManager(downloaderRegistry, repo, logg)

	// Initialize downloader manager (loads all enabled downloaders)
	ctx := context.Background()
	if err := downloaderManager.Initialize(ctx); err != nil {
		logg.Error().Err(err).Msg("failed to initialize downloader manager")
		// Don't fatal - allow server to start even if downloaders fail
	}

	// HTTP
	e := http.NewServer(cfg, logg, pool, services, repo, downloaderManager, broker)
	go func() {
		logg.Info().Str("port", cfg.Port).Msg("http listen")
		if err := e.Start(":" + cfg.Port); err != nil {
			log.Println("server stopped:", err)
		}
	}()

	// Download job worker
	workerCtx, workerCancel := context.WithCancel(context.Background())
	jobWorker := downloadjobs.New(repo, downloaderManager, services.Import, logg, broker)
	go jobWorker.Run(workerCtx)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop()
	workerCancel()

	shCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = e.Shutdown(shCtx)
	logg.Info().Msg("bye")
}
