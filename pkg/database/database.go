package database

import (
	"backend/pkg/config"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(cfg *config.DatabaseConfig, l *slog.Logger) (*gorm.DB, error) {
	gormLogger, err := NewSlogGormLogger(&cfg.GormLogger, l)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
