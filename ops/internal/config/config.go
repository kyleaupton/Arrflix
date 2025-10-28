package config

import (
	"os"
	"time"
)

type Config struct {
	RuntimeMode       string
	DatabaseURL       string
	PostgresDB        string
	PostgresUser      string
	PostgresPassword  string
	JWTSecret         string
	TmdbAPIKey        string
	ProwlarrAPIKey    string
	ReconcileInterval time.Duration
	NetworkName       string
}

func envOr(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func Load() Config {
	reconcileInterval := 10 * time.Second
	if interval := os.Getenv("RECONCILE_INTERVAL"); interval != "" {
		if parsed, err := time.ParseDuration(interval); err == nil {
			reconcileInterval = parsed
		}
	}

	return Config{
		RuntimeMode:       envOr("RUNTIME_MODE", "prod"),
		DatabaseURL:       envOr("DATABASE_URL", "postgres://snaggle:snaggle@snaggle-postgres:5432/snaggle?sslmode=disable"),
		PostgresDB:        envOr("POSTGRES_DB", "snaggle"),
		PostgresUser:      envOr("POSTGRES_USER", "snaggle"),
		PostgresPassword:  envOr("POSTGRES_PASSWORD", "snaggle"),
		JWTSecret:         envOr("JWT_SECRET", "dev-insecure-change-me"),
		TmdbAPIKey:        envOr("TMDB_API_KEY", ""),
		ProwlarrAPIKey:    envOr("PROWLARR_API_KEY", ""),
		ReconcileInterval: reconcileInterval,
		NetworkName:       envOr("NETWORK_NAME", "snaggle-network"),
	}
}
