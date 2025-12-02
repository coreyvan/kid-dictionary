package config

import (
	"log/slog"
	"os"
)

func NewLogger(cfg Config) *slog.Logger {
	var handler slog.Handler
	if cfg.PrettyLog {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: LogLevelFromString(cfg.LogLevel),
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: LogLevelFromString(cfg.LogLevel),
		})
	}
	return slog.New(handler)
}

func LogLevelFromString(levelStr string) slog.Level {
	switch levelStr {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
