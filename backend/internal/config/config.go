package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kyleaupton/snaggle/backend/internal/logger"
)

type Config struct {
	Env           string // dev|prod
	Port          string // API internal port, default 8080
	DatabaseURL   string
	CORSOrigin    string // used for dev SPA
	JWTSecret     string // HMAC secret for JWT signing
	TmdbAPIKey    string // TMDB API key
	JackettAPIKey string // Jackett API key
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func Load(log *logger.Logger) Config {
	// Best effort to load .env file
	godotenv.Load()

	config := Config{
		Env:           envOr("ENV", "dev"),
		Port:          envOr("PORT", "8080"),
		DatabaseURL:   envOr("DATABASE_URL", "postgres://snaggle:snaggle@127.0.0.1:5432/snaggle?sslmode=disable"),
		CORSOrigin:    envOr("SSE_ALLOW_ORIGIN", "*"),
		JWTSecret:     envOr("JWT_SECRET", "dev-insecure-change-me"),
		TmdbAPIKey:    envOr("TMDB_API_KEY", ""),
		JackettAPIKey: envOr("JACKETT_API_KEY", ""),
	}

	log.Debug().Interface("config", config).Msg("config")

	return config
}
