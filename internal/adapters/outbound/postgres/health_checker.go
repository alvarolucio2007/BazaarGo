package postgres

import (
	"context"
	"database/sql"
	"time"
)

type HealthChecker struct {
	db *sql.DB
}

// NewHealthChecker creates a new [HealthChecker].
func NewHealthChecker(db *sql.DB) *HealthChecker {
	return &HealthChecker{db: db}
}

// Health checks whether the postgres DB is healthy and properly connected.
//
// Returns an error if the [sql.DB.PingContext] health ping fails.
func (h *HealthChecker) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return h.db.PingContext(ctx)
}
