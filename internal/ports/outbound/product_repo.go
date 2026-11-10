package ports

import (
	"context"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	ReadByID(ctx context.Context, productID string) (domain.Product, error)
	ListByName(ctx context.Context, name string) ([]domain.Product, error)
	Update(ctx context.Context, product domain.Product) error
	Delete(ctx context.Context, productID string) error

	Ping(ctx context.Context) error
}
