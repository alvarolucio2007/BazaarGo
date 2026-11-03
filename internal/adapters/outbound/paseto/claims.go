package paseto

import (
	"errors"
	"time"
	"uuid"

	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
)

var (
	ErrExpiredToken = errors.New("paseto: token is expired")
	ErrInvalidToken = errors.New("paseto: token is invalid")
)

func NewClaims(email string, userID int64, duration time.Duration) (ports.Claims, error) {
	tokenID := uuid.New()
	timeNow := time.Now()
	payload := ports.Claims{
		ID:        tokenID,
		Email:     email,
		UserID:    userID,
		IssuedAt:  timeNow,
		ExpiresAt: timeNow.Add(duration),
	}
	return payload, nil
}
