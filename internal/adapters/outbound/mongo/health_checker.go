package mongo

import (
	"context"
	"time"

	"github.com/alvarolucio2007/bazaargo/internal/ports"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ (ports.HealthChecker) = (*HealthChecker)(nil)

type HealthChecker struct {
	mdb *mongo.Client
}

// NewHealthChecker creates a new [HealthChecker].
func NewHealthChecker(mdb *mongo.Client) *HealthChecker {
	return &HealthChecker{mdb: mdb}
}

// Health implements [ports.HealthChecker].
// Health checks whether mongoDB is healthy and properly connected.
//
// Returns an error if the [mongo.Client.Ping] health ping fails.
func (h *HealthChecker) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return h.mdb.Ping(ctx, nil)
}
