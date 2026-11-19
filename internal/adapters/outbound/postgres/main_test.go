package postgres_test

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testDB     *sql.DB
	testLogger *slog.Logger
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	testLogger = logger.NewStdSlog()
	var testContainer *postgres.PostgresContainer
	var err error
	testDB, testContainer, err = setupTestDB()
	if err != nil {
		log.Fatalf("unable to create DB and testcontainers %v", err)
	}
	if testDB == nil {
		log.Fatalf("testDB is nil")
	}
	if testContainer == nil {
		log.Fatalf("testDB is nil")
	}

	code := m.Run()
	if err := testContainer.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}
	if err := testDB.Close(); err != nil {
		log.Printf("failed to terminate DB: %v", err)
	}
	os.Exit(code)
}

func setupTestDB() (*sql.DB, *postgres.PostgresContainer, error) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"pgvector/pgvector:pg18-trixie",
		postgres.WithDatabase("scholarly_test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, nil, err
	}
	connStr, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, nil, err
	}

	if err := runMigrations(connStr); err != nil {
		return nil, nil, err
	}

	return db, pgContainer, nil
}

func runMigrations(connStr string) error {
	m, err := migrate.New("file:../../../../database/migrations/", connStr)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		return err
	}
	return nil
}
