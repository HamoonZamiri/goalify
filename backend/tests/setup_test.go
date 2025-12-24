package tests

import (
	"context"
	"goalify/cmd/app"
	"goalify/internal/config"
	"goalify/internal/testsetup"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	sqlcdb "goalify/internal/db/generated"
)

const BaseURL = "http://localhost:8080"

var (
	pgxPool       *pgxpool.Pool
	queries       *sqlcdb.Queries
	pgContainer   *postgres.PostgresContainer
	cancelServer  context.CancelFunc
	serverStopped chan struct{}
)

func setup(ctx context.Context) {
	var err error

	pgContainer, err = testsetup.GetPgContainer()
	if err != nil {
		panic(err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	// Set environment and database connection string before getting config
	config.SetEnv(config.ENV, string(config.LocalTest))
	config.SetEnv(config.TestDBConnString, connStr)

	// Reset config singleton to pick up test environment variables
	config.ResetForTesting()

	pgxPool, err = testsetup.GetPgxPool()
	if err != nil {
		panic(err)
	}

	queries = sqlcdb.New(pgxPool)

	// Create cancellable context for server
	serverCtx, cancel := context.WithCancel(ctx)
	cancelServer = cancel
	serverStopped = make(chan struct{})

	// Start server in goroutine with graceful shutdown
	go func() {
		defer close(serverStopped)
		if err := app.Run(serverCtx); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	setup(ctx)

	code := m.Run()

	if cancelServer != nil {
		log.Println("Shutting down test server...")
		cancelServer()

		select {
		case <-serverStopped:
			log.Println("Test server stopped gracefully")
		case <-time.After(5 * time.Second):
			log.Println("Test server shutdown timed out")
		}
	}

	if pgContainer != nil {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("Failed to terminate container: %s", err)
		}
	}

	os.Exit(code)
}
