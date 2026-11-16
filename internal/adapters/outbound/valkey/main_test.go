package valkey_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testRDB *redis.Client
)

func cleanupTest(t *testing.T) {
	t.Cleanup(func() {
		testRDB.FlushAll(t.Context())
	})
}

func TestMain(m *testing.M) {
	var rct *testcontainers.DockerContainer
	var err error
	testRDB, rct, err = setupTestValkey()
	if err != nil {
		log.Fatalf("unable to create rdb and testcontainer %v", err)
	}
	if testRDB == nil {
		log.Fatalf("rdb is nil")
	}
	if rct == nil {
		log.Fatalf("rct is nil")
	}
	code := m.Run()
	if err := testcontainers.TerminateContainer(rct); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
	if err := testRDB.Close(); err != nil {
		log.Printf("failed to close rdb: %v", err)
	}
	os.Exit(code)
}
func setupTestValkey() (*redis.Client, *testcontainers.DockerContainer, error) {
	ctx := context.Background()
	valkeyC, err := testcontainers.Run(
		ctx, "valkey/valkey:9-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("6379/tcp"),
			wait.ForLog("Ready to accept connections"),
		),
	)
	if err != nil {
		return nil, nil, err
	}
	endpoint, err := valkeyC.Endpoint(ctx, "")
	if err != nil {
		return nil, nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})
	return client, valkeyC, nil
}
