package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
)

// New returns a configured *gorm.DB instance for the selected driver.
func New(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dialector, err := createDialector(cfg)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return db, nil
}

func createDialector(cfg config.DatabaseConfig) (gorm.Dialector, error) {
	switch cfg.Driver {
	case "postgres":
		return postgres.Open(cfg.DSN), nil
	case "sqlite":
		return sqlite.Open(cfg.DSN), nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}
