package ports

import (
	"context"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	ReadByName(ctx context.Context, name string) (domain.Product, error)
	Update(ctx context.Context, product domain.Product) error
	Delete(ctx context.Context, productID primitive.ObjectID) error
}
