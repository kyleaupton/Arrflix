package config

import (
	"os"
)

type Config struct {
	Env         string // dev|prod
	Port        string // API internal port, default 8080
	DatabaseURL string
	CORSOrigin  string // used for dev SPA
    JWTSecret   string // HMAC secret for JWT signing
}

func envOr(k, d string) string { if v := os.Getenv(k); v != "" { return v }; return d }

func Load() Config {
	return Config{
		Env:         envOr("ENV", "dev"),
		Port:        envOr("PORT", "8080"),
		DatabaseURL: envOr("DATABASE_URL", "postgres://snaggle:snaggle@127.0.0.1:5432/snaggle?sslmode=disable"),
		CORSOrigin:  envOr("SSE_ALLOW_ORIGIN", "*"),
        JWTSecret:   envOr("JWT_SECRET", "dev-insecure-change-me"),
	}
}
