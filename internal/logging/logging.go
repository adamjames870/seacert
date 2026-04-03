package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type PrettyHandler struct {
	slog.Handler
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()

	timeStr := r.Time.Format("15:04:05")
	msg := r.Message

	levelColor := ""
	resetColor := "\033[0m"

	switch r.Level {
	case slog.LevelDebug:
		level = "DEBUG"
		levelColor = "\033[35m" // Magenta
	case slog.LevelInfo:
		level = "INFO"
		levelColor = "\033[32m" // Green
	case slog.LevelWarn:
		level = "WARN"
		levelColor = "\033[33m" // Yellow
	case slog.LevelError:
		level = "ERROR"
		levelColor = "\033[31m" // Red
	}

	fmt.Printf("[%s] %s%-5s%s %s", timeStr, levelColor, level, resetColor, msg)

	r.Attrs(func(a slog.Attr) bool {
		h.formatAttr(a)
		return true
	})

	fmt.Println()

	return nil
}

func (h *PrettyHandler) formatAttr(a slog.Attr) {
	if a.Value.Kind() == slog.KindGroup {
		fmt.Printf(" %s=[", a.Key)
		for i, attr := range a.Value.Group() {
			if i > 0 {
				fmt.Print(" ")
			}
			h.formatAttr(attr)
		}
		fmt.Print("]")
	} else {
		fmt.Printf(" %s=%v", a.Key, a.Value.Any())
	}
}

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
		opts := &slog.HandlerOptions{
			Level: level,
		}
		handler = &PrettyHandler{
			Handler: slog.NewTextHandler(os.Stdout, opts),
		}
	}

	return slog.New(handler)
}

func InitDefaultLogger() {
	logger := NewLogger()
	slog.SetDefault(logger)
}
