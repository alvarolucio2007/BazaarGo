package logger

import (
	"log/slog"
	"os"
)

// NewSlog creates a new [*slog.Logger] JSON Handler.
func NewSlog() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
