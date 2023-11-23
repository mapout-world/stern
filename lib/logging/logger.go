package logging

import (
	"context"
	"io"
	"log/slog"
	"strings"

	"github.com/mapout-world/stern/lib/config"
)

var LogLevel = new(slog.LevelVar)

func NewLogger(ctx context.Context, w io.Writer) *slog.Logger {
	SetLevel(config.Get("logging", "level").String("info"))

	logger := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: LogLevel,
	}))

	config.NewWatcher(logger, "logging", "level").Watch(ctx, func(val config.Value) {
		SetLevel(val.String("info"))
	})

	return logger
}

func SetLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		LogLevel.Set(slog.LevelDebug)
	case "info":
		LogLevel.Set(slog.LevelInfo)
	case "warn":
		LogLevel.Set(slog.LevelWarn)
	case "error":
		LogLevel.Set(slog.LevelError)
	}
}
