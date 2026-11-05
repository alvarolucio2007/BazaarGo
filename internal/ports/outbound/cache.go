package ports

import "context"

type Cache[T any] interface {
	Create(ctx context.Context, item T) error
	Read(ctx context.Context, itemID string) (T, error)
	Delete(ctx context.Context, itemID string) error
}
