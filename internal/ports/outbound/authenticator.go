package ports

import (
	"time"
	"uuid"
)

type Authenticator interface {
	CreateToken(userID int64) (string, error)
	VerifyToken(token string) (int64, error)
}
type Claims struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	UserID    int64     `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
