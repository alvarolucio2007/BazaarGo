package paseto

import (
	"errors"
	"time"
	"uuid"

	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
)

var (
	ErrExpiredToken    = errors.New("paseto: token is expired")
	ErrInvalidToken    = errors.New("paseto: token is invalid")
	ErrInvalidUserID   = errors.New("paseto: userID must be greater than zero")
	ErrInvalidDuration = errors.New("paseto: token duration must be positive")
)

// NewClaims creates a new [ports.Claims].
//
// Returns [ErrInvalidUserID] if userID is less than or equal to zero.
// Returns [ErrInvalidDuration] if duration is negative.
func NewClaims(email string, userID int64, duration time.Duration) (ports.Claims, error) {
	if userID <= 0 {
		return ports.Claims{}, ErrInvalidUserID
	}
	if duration < 0 {
		return ports.Claims{}, ErrInvalidDuration
	}
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
