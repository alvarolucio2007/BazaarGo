package ports

import (
	"errors"
	"time"
	"uuid"
)

type Authenticator interface {
	CreateToken(claims Claims) (string, error)
	VerifyToken(token string) (Claims, error)
}
type Claims struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	UserID    int64     `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

var (
	ErrInvalidToken    = errors.New("paseto: token is invalid")
	ErrInvalidUserID   = errors.New("paseto: userID must be greater than zero")
	ErrInvalidDuration = errors.New("paseto: token duration must be positive")
)

// NewClaims creates new [Claims].
//
// Returns [ErrInvalidUserID] if userID is less than or equal to zero.
// Returns [ErrInvalidDuration] if duration is negative.
func NewClaims(email string, userID int64, duration time.Duration) (Claims, error) {
	if userID <= 0 {
		return Claims{}, ErrInvalidUserID
	}
	if duration < 0 {
		return Claims{}, ErrInvalidDuration
	}
	tokenID := uuid.New()
	timeNow := time.Now()
	payload := Claims{
		ID:        tokenID,
		Email:     email,
		UserID:    userID,
		IssuedAt:  timeNow,
		ExpiresAt: timeNow.Add(duration),
	}
	return payload, nil
}
