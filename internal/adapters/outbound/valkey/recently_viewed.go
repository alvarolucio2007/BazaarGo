package valkey

import (
	"context"
	"fmt"
	"time"

	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
	"github.com/redis/go-redis/v9"
)

var _ ports.RecentlyViewedRepository = (*RecentlyViewedRepo)(nil)

type RecentlyViewedRepo struct {
	rdb *redis.Client
}

func NewRecentlyViewedRepo(rdb *redis.Client) *RecentlyViewedRepo {
	return &RecentlyViewedRepo{rdb}
}

// Add implements [ports.RecentlyViewedRepository].
func (r *RecentlyViewedRepo) Add(ctx context.Context, clientID string, productID string) error {
	key := fmt.Sprintf("recent:%s", clientID)
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()

	pipe := r.rdb.TxPipeline()
	pipe.LRem(ctx, key, 0, productID)
	pipe.LPush(ctx, key, productID)
	pipe.LTrim(ctx, key, 0, 4)
	pipe.Expire(ctx, key, 7*24*time.Hour)
	_, err := pipe.Exec(ctx)
	return err
}

// Get implements [ports.RecentlyViewedRepository].
func (r *RecentlyViewedRepo) Get(ctx context.Context, clientID string) ([]string, error) {
	key := fmt.Sprintf("recent:%s", clientID)
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()

	return r.rdb.LRange(ctx, key, 0, 4).Result()
}
