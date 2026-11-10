package valkey

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"time"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/codes"
)

var _ ports.Cache[domain.Product] = (*ProductCache)(nil)

// ProductCache implements the [ports.Cache] interface for Redis (Valkey).
type ProductCache struct {
	rdb *redis.Client
}

// NewProductCache creates a new instance of ProductCache.
func NewProductCache(rdb *redis.Client) *ProductCache {
	return &ProductCache{rdb: rdb}
}

// Create implements [ports.Cache].
// It marshals a [domain.Product] into JSON and stores it in Redis with a predefined TTL.
//
// Returns an error if JSON marshaling fails or if the Redis operation cannot be completed.
func (p *ProductCache) Create(ctx context.Context, item domain.Product) error {
	ctx, span := tracer.Start(ctx, "ProductCache.Create")
	defer span.End()

	cacheKey := fmt.Sprintf("product:%s", item.ID)

	jsonData, err := json.Marshal(item)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("cache: marshal product error: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()

	if err := p.rdb.Set(ctx, cacheKey, jsonData, CacheExpTime).Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("cache: set product error: %w", err)
	}
	return nil
}

// Read implements [ports.Cache].
// It fetches and unmarshals a [domain.Product] from Redis based on its ID.
//
// Returns [domain.ErrNotFound] if the product is not found in the cache.
// Returns an error if either the Redis operation fails or if the stored data is malformed.
func (p *ProductCache) Read(ctx context.Context, itemID string) (domain.Product, error) {
	ctx, span := tracer.Start(ctx, "ProductCache.Read")
	defer span.End()

	cacheKey := fmt.Sprintf("product:%s", itemID)

	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()

	data, err := p.rdb.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return domain.Product{}, domain.ErrNotFound
	} else if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.Product{}, fmt.Errorf("cache: get product error: %w", err)
	}

	var product domain.Product
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.Product{}, fmt.Errorf("cache: unmarshal product error: %w", err)
	}

	return product, nil
}

// Delete implements [ports.Cache].
// It removes a product from the Redis cache by its ID.
//
// Returns an error if the Redis operation cannot be completed.
func (p *ProductCache) Delete(ctx context.Context, itemID string) error {
	ctx, span := tracer.Start(ctx, "ProductCache.Delete")
	defer span.End()

	cacheKey := fmt.Sprintf("product:%s", itemID)

	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()

	if err := p.rdb.Del(ctx, cacheKey).Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("cache: delete product error: %w", err)
	}

	return nil
}

// Ping implements [ports.Cache].
// Ping checks whether the valkey cache is healthy and properly connected.
//
// Returns an error if the [redis.Client.Ping] health ping fails.
func (p *ProductCache) Ping(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "ProductCache.Ping")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return p.rdb.Ping(ctx).Err()
}
