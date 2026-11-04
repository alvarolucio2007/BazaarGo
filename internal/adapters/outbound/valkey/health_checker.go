package valkey

import (
	"context"
	"time"

	"github.com/alvarolucio2007/bazaargo/internal/ports"
	"github.com/redis/go-redis/v9"
)

var _ (ports.HealthChecker) = (*HealthChecker)(nil)

type HealthChecker struct {
	rdb *redis.Client
}

// NewHealthChecker creates a new [HealthChecker].
func NewHealthChecker(rdb *redis.Client) *HealthChecker {
	return &HealthChecker{rdb: rdb}
}

// Health implements [ports.HealthChecker].
// Health checks whether the valkey cache is healthy and properly connected.
//
// Returns an error if the [redis.Client.Ping] health ping fails.
func (h *HealthChecker) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return h.rdb.Ping(ctx).Err()
}
