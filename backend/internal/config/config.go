package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string // dev|prod
	Port        string // API internal port, default 8080
	DatabaseURL string
	CORSOrigin  string // used for dev SPA
	JWTSecret   string // HMAC secret for JWT signing
	TmdbAPIKey  string // TMDB API key
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return Config{
		Env:         envOr("ENV", "dev"),
		Port:        envOr("PORT", "8080"),
		DatabaseURL: envOr("DATABASE_URL", "postgres://snaggle:snaggle@127.0.0.1:5432/snaggle?sslmode=disable"),
		CORSOrigin:  envOr("SSE_ALLOW_ORIGIN", "*"),
		JWTSecret:   envOr("JWT_SECRET", "dev-insecure-change-me"),
		TmdbAPIKey:  envOr("TMDB_API_KEY", ""),
	}
}
