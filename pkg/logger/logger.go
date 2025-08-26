package logger

import (
	"backend/pkg/config"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func InitLogger(cfg *config.LoggerConfig) (*slog.Logger, error) {
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return nil, fmt.Errorf("invalid log level: %s", cfg.Level)
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
	}

	if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	timestamp := time.Now().Format("20060102150405")

	logFile, err := os.OpenFile(
		filepath.Join(cfg.Dir, "app."+timestamp+".log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(logFile, opts))

	return logger, nil
}
