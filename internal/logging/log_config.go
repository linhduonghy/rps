package logging

import (
	"log/slog"
	"os"
	"rps/config"
)

func New(cfg config.LogConfig) *slog.Logger {
	level := parseLevel(cfg.Level)

	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: level,
	}
	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)

	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug

	case "warn":
		return slog.LevelWarn

	case "errors":
		return slog.LevelError

	default:
		return slog.LevelInfo
	}
}
