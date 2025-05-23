package stores

import (
	"context"
	"goalify/db"
	"goalify/testsetup"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	dbConn      *sqlx.DB
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

	dbConn, err = db.NewWithConnString(connStr)
	if err != nil {
		panic(err)
	}

	cStore = NewChestStore(dbConn)
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	setup(ctx)
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("Failed to terminate container: %s", err)
		}
	}()

	code := m.Run()
	os.Exit(code)
}

func TestCreateChest(t *testing.T) {
	t.Parallel()

	chest, err := cStore.CreateChest("bronze", "bronze chest", 100)
	assert.Nil(t, err)
	assert.NotNil(t, chest)
	assert.Equal(t, 100, chest.Price)
}
