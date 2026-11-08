package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
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
