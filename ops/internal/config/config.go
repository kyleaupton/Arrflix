package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
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

// generateRandomKey generates a cryptographically secure random key
func generateRandomKey(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to a deterministic but still random-looking key if crypto/rand fails
		return "generated-fallback-key-" + hex.EncodeToString([]byte(time.Now().String()))[:length]
	}
	return hex.EncodeToString(bytes)
}

func Load() Config {
	godotenv.Load()

	reconcileInterval := 10 * time.Second
	if interval := os.Getenv("RECONCILE_INTERVAL"); interval != "" {
		if parsed, err := time.ParseDuration(interval); err == nil {
			reconcileInterval = parsed
		}
	}

	// Generate a random Prowlarr API key if not provided
	prowlarrAPIKey := os.Getenv("PROWLARR_API_KEY")
	if prowlarrAPIKey == "" {
		prowlarrAPIKey = generateRandomKey(32) // 32 bytes = 64 hex characters
		log.Printf("Generated random Prowlarr API key: %s", prowlarrAPIKey)
	}

	return Config{
		RuntimeMode:       envOr("RUNTIME_MODE", "prod"),
		DatabaseURL:       envOr("DATABASE_URL", "postgres://snaggle:snaggle@snaggle-postgres:5432/snaggle?sslmode=disable"),
		PostgresDB:        envOr("POSTGRES_DB", "snaggle"),
		PostgresUser:      envOr("POSTGRES_USER", "snaggle"),
		PostgresPassword:  envOr("POSTGRES_PASSWORD", "snaggle"),
		JWTSecret:         envOr("JWT_SECRET", "todo-insecure-change-me"),
		TmdbAPIKey:        envOr("TMDB_API_KEY", ""),
		ProwlarrAPIKey:    prowlarrAPIKey,
		ReconcileInterval: reconcileInterval,
		NetworkName:       envOr("NETWORK_NAME", "snaggle-network"),
	}
}
