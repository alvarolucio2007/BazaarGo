package valkey

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const CacheQueryTimeout = 500 * time.Millisecond
const CacheExpTime = time.Hour

// NewConn creates a new Valkey connection
//
// Returns an error if the [redis.Client.Ping] health ping fails.
func NewConn(address, password string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
		Protocol: 2,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
