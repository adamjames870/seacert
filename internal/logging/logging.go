package logging

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	platform := os.Getenv("PLATFORM")

	var handler slog.Handler
	if platform == "production" || platform == "test" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return slog.New(handler)
}

func InitDefaultLogger() {
	logger := NewLogger()
	slog.SetDefault(logger)
}
