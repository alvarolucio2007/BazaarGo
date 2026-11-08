package ports

import "context"

type RecentlyViewedRepository interface {
	Add(ctx context.Context, clientID, productID string) error
	Get(ctx context.Context, clientID string) ([]string, error)
}
