package valkey

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const CacheQueryTimeout = 500 * time.Millisecond
const CacheExpTime = time.Hour

var ErrCacheMiss = errors.New("cache miss")

// NewConn creates a new Valkey connection.
//
// Returns an error if the [redis.Client.Ping] health ping fails.
func NewConn(address, password string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:                  address,
		Password:              password,
		DB:                    db,
		Protocol:              2,
		ContextTimeoutEnabled: true,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, err
	}
	return rdb, nil
}
