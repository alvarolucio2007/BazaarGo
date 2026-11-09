package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"
)

var protectedOperations = map[string]bool{
	"CreateProduct":    true,
	"UpdateProduct":    true,
	"DeleteProductsId": true,
}

type userKey string

const userCtx userKey = "user"

func (h *ServerHandler) AuthStrictMiddleware(next StrictHandlerFunc, operationID string) StrictHandlerFunc {
	if !protectedOperations[operationID] {
		return next
	}

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			return nil, domain.ErrUnauthorized
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return nil, domain.ErrUnauthorized
		}

		pasetoToken, err := h.authenticator.VerifyToken(parts[1])
		if err != nil {
			return nil, domain.ErrUnauthorized
		}
		user, err := h.userService.Read(ctx, pasetoToken.UserID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return nil, domain.ErrUnauthorized
			}
			return nil, err
		}
		ctx = context.WithValue(ctx, userCtx, user)
		return next(ctx, w, r, request)
	}
}
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			next.ServeHTTP(ww, r)
			span := trace.SpanContextFromContext(r.Context())
			reqLogger := logger
			if span.HasTraceID() {
				reqLogger = reqLogger.With("trace_id", span.TraceID().String())
			}
			reqLogger.InfoContext(r.Context(), "request completed",
				"status", ww.Status(),
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(t1),
			)
		}
		return http.HandlerFunc(fn)
	}
}
