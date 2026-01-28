package internal

import (
	"github.com/kyleaupton/arrflix/internal/config"
	"github.com/kyleaupton/arrflix/internal/db"
	"github.com/kyleaupton/arrflix/internal/logger"
	"github.com/kyleaupton/arrflix/internal/repo"
)

func GetRepo() *repo.Repository {
	// Logger
	logg := logger.New(true)

	// Load config
	cfg := config.Load(logg)

	// DB
	pool, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		logg.Fatal().Err(err).Msg("open db")
	}

	return repo.New(pool)
}
