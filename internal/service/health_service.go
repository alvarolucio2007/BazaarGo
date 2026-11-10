package service

import (
	"context"

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

func (s *HealthService) Check(ctx context.Context) (bool, bool, bool, error) {
	var pgOk, mongoOk, cacheOk bool

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		err := s.pg.Ping(gCtx)
		if err == nil {
			pgOk = true
		}
		return err
	})
	g.Go(func() error {
		err := s.mongo.Ping(gCtx)
		if err == nil {
			mongoOk = true
		}
		return err
	})
	g.Go(func() error {
		err := s.cache.Ping(gCtx)
		if err == nil {
			cacheOk = true
		}
		return err
	})

	err := g.Wait()
	return pgOk, mongoOk, cacheOk, err
}
