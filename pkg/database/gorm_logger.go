// file: internal/pkg/database/gorm_logger.go
package database

import (
	"backend/pkg/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm/logger"
)

func NewSlogGormLogger(cfg *config.GormLoggerConfig, l *slog.Logger) (logger.Interface, error) {
	var level logger.LogLevel
	switch cfg.Level {
	case "silent":
		level = logger.Silent
	case "error":
		level = logger.Error
	case "warn":
		level = logger.Warn
	case "info":
		level = logger.Info
	default:
		return nil, fmt.Errorf("invalid log level for gorm logger: %s", cfg.Level)
	}

	return &slogGormLogger{
		logLevel:      level,
		slowThreshold: time.Duration(cfg.SlowThreshold) * time.Millisecond,
		logger:        l,
	}, nil
}

type slogGormLogger struct {
	logLevel      logger.LogLevel
	slowThreshold time.Duration
	logger        *slog.Logger
}

func (l *slogGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

func (l *slogGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.InfoContext(ctx, msg, data...)
}

func (l *slogGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WarnContext(ctx, msg, data...)
}

func (l *slogGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.ErrorContext(ctx, msg, data...)
}

func (l *slogGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// if log level is silent, return
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fieldLogger := l.logger.With(
		"sql", sql,
		"duration", elapsed.String(),
		"rows", rows,
	)

	if elapsed > l.slowThreshold {
		fieldLogger.WarnContext(ctx, "GORM Slow Query")
	} else if l.logLevel <= logger.Info {
		fieldLogger.DebugContext(ctx, "GORM Trace")
	}
}
