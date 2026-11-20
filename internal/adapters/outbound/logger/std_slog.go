package logger

import (
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

// NewSlog creates a new [slog.Logger] OTel Handler.
func NewSlog(name string) *slog.Logger {
	handler := otelslog.NewHandler(name)
	return slog.New(handler)
}

func NewStdSlog() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
