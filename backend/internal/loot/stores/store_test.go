package stores

import (
	"context"
	"goalify/internal/db"
	"goalify/internal/testsetup"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	sqlcdb "goalify/internal/db/generated"
)

var (
	cStore      LootStore
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

	pgxPool, err := db.NewPgxPoolWithConnString(ctx, connStr)
	if err != nil {
		panic(err)
	}

	queries := sqlcdb.New(pgxPool)
	cStore = NewChestStore(queries)
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	setup(ctx)

	code := m.Run()

	if err := pgContainer.Terminate(ctx); err != nil {
		log.Fatalf("Failed to terminate container: %s", err)
	}

	os.Exit(code)
}

func TestCreateChest(t *testing.T) {
	t.Parallel()

	chest, err := cStore.CreateChest("bronze", "bronze chest", 100)
	assert.Nil(t, err)
	assert.NotNil(t, chest)
	assert.Equal(t, 100, chest.Price)
}
