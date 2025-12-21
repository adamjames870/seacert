package logging

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	platform := os.Getenv("PLATFORM")

	var handler slog.Handler
	var level slog.Level

	switch platform {
	case "production":
		level = slog.LevelInfo
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	case "test":
		level = slog.LevelDebug
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	default:
		// dev or other
		level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}

	return slog.New(handler)
}

func InitDefaultLogger() {
	logger := NewLogger()
	slog.SetDefault(logger)
}
