package logger

import (
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

// NewSlog creates a new [slog.Logger] OTel Handler.
func NewSlog() *slog.Logger {
	handler := otelslog.NewHandler("bazaargo-api")
	return slog.New(handler)
}
