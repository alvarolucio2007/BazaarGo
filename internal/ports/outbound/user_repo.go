package ports

import (
	"context"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	Read(ctx context.Context, userID int64) (domain.User, error)
	ReadByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, userID int64) error
}
