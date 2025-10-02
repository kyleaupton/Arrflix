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
)

func main() {
	// Load config
	cfg := config.Load()

	// Logger
	logg := logger.New(cfg.Env)

	// DB
	pool, err := db.Open(cfg.DatabaseURL)
	if err != nil { logg.Fatal().Err(err).Msg("open db") }
	defer pool.Close()

	// Migrations (run on startup; idempotent)
	// if err := db.ApplyMigrations(pool, "db/migrations"); err != nil {
	// 	logg.Fatal().Err(err).Msg("migrate")
	// }

	// HTTP
	e := http.NewServer(cfg, logg, pool)
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
