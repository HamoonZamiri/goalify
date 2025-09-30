package tests

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"goalify/cmd/app"
	"goalify/config"
	sqlcdb "goalify/db/generated"
	"goalify/testsetup"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const BASE_URL = "http://localhost:8080"

var (
	pgxPool     *pgxpool.Pool
	queries     *sqlcdb.Queries
	pgContainer *postgres.PostgresContainer
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
	config.SetEnv(config.TEST_DB_CONN_STRING, connStr)

	// Reset config singleton to pick up test environment variables
	config.ResetForTesting()

	pgxPool, err = testsetup.GetPgxPool()
	if err != nil {
		panic(err)
	}

	queries = sqlcdb.New(pgxPool)

	go app.Run()
	time.Sleep(50 * time.Millisecond)
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	// start server in a goroutine
	setup(ctx)
	code := m.Run()

	var err error
	defer func() {
		if err = pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("Failed to terminate container: %s", err)
		}
	}()

	os.Exit(code)
}