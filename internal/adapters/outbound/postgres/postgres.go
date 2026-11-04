package postgres

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const PostgresQueryTimeout time.Duration = 5 * time.Second

// NewConn creates a new postgres connection with [*sql.DB].
//
// Returns an error if either the pgx driver fails or cannot connect to the database by the addr.
// Returns an error if the [sql.DB] health ping fails.
func NewConn(addr string, maxOpenConns, maxIdleConns int, maxIdleTime time.Duration) (*sql.DB, error) {
	db, err := sql.Open("pgx", addr)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxIdleTime)
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
