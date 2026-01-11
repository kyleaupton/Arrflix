package internal

import (
	"github.com/kyleaupton/Arrflix/internal/config"
	"github.com/kyleaupton/Arrflix/internal/db"
	"github.com/kyleaupton/Arrflix/internal/logger"
	"github.com/kyleaupton/Arrflix/internal/repo"
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
