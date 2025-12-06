package config

import (
	"os"
	"strconv"
	"strings"
)

// Config aggregates all configuration for the application.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

// AppConfig contains Fiber specific settings.
type AppConfig struct {
	Name string
	Port string
	Env  string
}

// DatabaseConfig captures everything required to boot GORM.
type DatabaseConfig struct {
	Driver      string
	DSN         string
	AutoMigrate bool
}

// Load builds a Config struct with sane defaults overridden by environment variables.
func Load() Config {
	return Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "grtblog-api"),
			Port: getEnv("APP_PORT", "8080"),
			Env:  strings.ToLower(getEnv("APP_ENV", "development")),
		},
		Database: DatabaseConfig{
			Driver:      strings.ToLower(getEnv("DB_DRIVER", "sqlite")),
			DSN:         getEnv("DB_DSN", "./storage/grtblog.db"),
			AutoMigrate: getEnvAsBool("DB_AUTO_MIGRATE", true),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return boolVal
}
