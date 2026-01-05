package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config aggregates all configuration for the application.
type Config struct {
	App       AppConfig
	Database  DatabaseConfig
	Auth      AuthConfig
	RBAC      RBACConfig
	Turnstile TurnstileConfig
	Redis     RedisConfig
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

// AuthConfig 控制 JWT 签发与校验。
type AuthConfig struct {
	Secret        string
	Issuer        string
	AccessTTL     time.Duration
	DefaultRoles  []string
	OAuthStateTTL time.Duration
}

// RBACConfig 管理 Casbin 相关设置。
type RBACConfig struct {
	ModelPath         string
	AutoReload        bool
	AutoReloadSeconds int
}

// TurnstileConfig 控制 Cloudflare Turnstile 人机校验。
type TurnstileConfig struct {
	Enabled   bool
	Secret    string
	VerifyURL string
	Timeout   time.Duration
}

// RedisConfig 描述 Redis 连接。
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Prefix   string
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
			Driver:      strings.ToLower(getEnv("DB_DRIVER", "postgres")),
			DSN:         getEnv("DB_DSN", "postgres://postgres:postgres@localhost:5432/grtblog?sslmode=disable"),
			AutoMigrate: getEnvAsBool("DB_AUTO_MIGRATE", true),
		},
		Auth: AuthConfig{
			Secret:        getEnv("AUTH_SECRET", "change-me"),
			Issuer:        getEnv("AUTH_ISSUER", "grtblog-api"),
			AccessTTL:     getEnvAsDuration("AUTH_ACCESS_TTL", 7*24*time.Hour),
			DefaultRoles:  getEnvAsSlice("AUTH_DEFAULT_ROLES", []string{"user"}),
			OAuthStateTTL: getEnvAsDuration("AUTH_STATE_TTL", time.Minute*10),
		},
		RBAC: RBACConfig{
			ModelPath:         getEnv("RBAC_MODEL_PATH", "./configs/rbac_model.conf"),
			AutoReload:        getEnvAsBool("RBAC_AUTO_RELOAD", false),
			AutoReloadSeconds: getEnvAsInt("RBAC_AUTO_RELOAD_SECONDS", 30),
		},
		Turnstile: TurnstileConfig{
			Enabled:   getEnvAsBool("TURNSTILE_ENABLED", false),
			Secret:    getEnv("TURNSTILE_SECRET", ""),
			VerifyURL: getEnv("TURNSTILE_VERIFY_URL", "https://challenges.cloudflare.com/turnstile/v0/siteverify"),
			Timeout:   getEnvAsDuration("TURNSTILE_TIMEOUT", 5*time.Second),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			Prefix:   getEnv("REDIS_PREFIX", "grtblog:"),
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

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return d
}

func getEnvAsSlice(key string, fallback []string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parts := strings.Split(value, ",")
	var result []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	if len(result) == 0 {
		return fallback
	}
	return result
}

func getEnvAsInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return i
}
