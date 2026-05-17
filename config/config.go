package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr    string
	DatabaseDSN string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		HTTPAddr:    resolveHTTPAddr(),
		DatabaseDSN: resolveDatabaseDSN(),
	}

	if cfg.DatabaseDSN == "" {
		return Config{}, fmt.Errorf("DATABASE_DSN or DATABASE_URL is required")
	}
	return cfg, nil
}

func resolveHTTPAddr() string {
	if addr := strings.TrimSpace(os.Getenv("HTTP_ADDR")); addr != "" {
		return addr
	}
	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		if strings.Contains(port, ":") && !strings.HasPrefix(port, ":") {
			return port
		}
		if strings.HasPrefix(port, ":") {
			return port
		}
		return ":" + port
	}
	return ":8080"
}

func resolveDatabaseDSN() string {
	if dsn := strings.TrimSpace(os.Getenv("DATABASE_DSN")); dsn != "" {
		return dsn
	}
	if url := strings.TrimSpace(os.Getenv("DATABASE_URL")); url != "" {
		return url
	}
	return "postgres://postgres:postgres@localhost:5432/subsservice?sslmode=disable"
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
