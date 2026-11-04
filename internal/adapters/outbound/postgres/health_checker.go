package postgres

import "database/sql"

type HealthChecker struct {
	db *sql.DB
}

// NewHealthChecker creates a new [HealthChecker].
func NewHealthChecker(db *sql.DB) *HealthChecker {
	return &HealthChecker{db: db}
}
