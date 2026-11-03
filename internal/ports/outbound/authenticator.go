package ports

import (
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
