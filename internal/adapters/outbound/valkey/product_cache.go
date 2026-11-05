package valkey

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
	"github.com/redis/go-redis/v9"
)

var _ ports.Cache[domain.Product] = (*ProductCache)(nil)

type ProductCache struct {
	rdb *redis.Client
}

// Create implements [ports.Cache].
// It marshals a [domain.Product] into JSON and stores it in Redis with a predefined TTL.
//
// Returns an error if JSON marshaling fails or if the Redis operation cannot be completed.
func (p *ProductCache) Create(ctx context.Context, item domain.Product) error {
	cacheKey := fmt.Sprintf("product:%s", item.ID)

	jsonData, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("cache: create product error: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()

	return p.rdb.Set(ctx, cacheKey, jsonData, CacheExpTime).Err()
}

// Read implements [ports.Cache].
// It fetches and unmarshals a [domain.Product] from Redis based on its ID.
//
// Returns [ports.ErrCacheMiss] (or your domain equivalent) if the product is not found in the cache.
// Returns an error if the Redis operation fails or if the stored data is malformed.
func (p *ProductCache) Read(ctx context.Context, itemID string) (domain.Product, error) {
	cacheKey := fmt.Sprintf("product:%s", itemID)

	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()

	data, err := p.rdb.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return domain.Product{}, ErrCacheMiss
	} else if err != nil {
		return domain.Product{}, fmt.Errorf("cache: read product error: %w", err)
	}

	var product domain.Product
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		return domain.Product{}, fmt.Errorf("cache: unmarshal product error: %w", err)
	}

	return product, nil
}

// Delete implements [ports.Cache].
// It removes a product from the Redis cache by its ID.
//
// Returns an error if the Redis operation cannot be completed.
func (p *ProductCache) Delete(ctx context.Context, itemID string) error {
	cacheKey := fmt.Sprintf("product:%s", itemID)

	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()

	if err := p.rdb.Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("cache: delete product error: %w", err)
	}

	return nil
}
