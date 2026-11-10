package service

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
	"golang.org/x/sync/errgroup"
)

type HealthService struct {
	pg    ports.UserRepository
	mongo ports.ProductRepository
	cache ports.Cache[domain.Product]
}

func NewHealthService(pg ports.UserRepository, mongo ports.ProductRepository, cache ports.Cache[domain.Product]) *HealthService {
	return &HealthService{pg: pg, mongo: mongo, cache: cache}
}

// Check performs a concurrent health check on all external dependencies
// (PostgreSQL, MongoDB, and Valkey/Redis).
//
// It uses an errgroup to ping all services in parallel, applying a strict
// 3-second timeout to prevent the endpoint from hanging if a database is down.
//
// Returns the individual health status of each dependency in the order:
// (PostgreSQL, MongoDB, Valkey), along with an error if any of the pings fail.
// Atomic booleans are used to ensure thread-safe reads and writes across goroutines.
func (s *HealthService) Check(ctx context.Context) (bool, bool, bool, error) {
	ctx, span := tracer.Start(ctx, "HealthService.Check")
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var pgOk, mongoOk, cacheOk atomic.Bool

	var g errgroup.Group

	g.Go(func() error {
		err := s.pg.Ping(ctx)
		if err == nil {
			pgOk.Store(true)
		}
		return err
	})

	g.Go(func() error {
		err := s.mongo.Ping(ctx)
		if err == nil {
			mongoOk.Store(true)
		}
		return err
	})

	g.Go(func() error {
		err := s.cache.Ping(ctx)
		if err == nil {
			cacheOk.Store(true)
		}
		return err
	})

	err := g.Wait()

	return pgOk.Load(), mongoOk.Load(), cacheOk.Load(), err
}
