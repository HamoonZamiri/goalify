package testsetup

import (
	"context"
	"goalify/config"
	"goalify/db"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	pgContainer *postgres.PostgresContainer
	dbx         *sqlx.DB
)

func setupPgContainer() error {
	ctx := context.Background()
	if pgContainer != nil {
		return nil
	}

	dbName := "goalify_test"
	dbUser := "goalify"
	dbPassword := "goalify"

	var err error
	pgContainer, err = postgres.Run(ctx, "docker.io/postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return err
	}
	os.Setenv(config.TEST_DB_CONN_STRING, connStr)

	dbx, err = db.NewWithConnString(connStr)
	if err != nil {
		panic(err)
	}

	// using goose run migrations from db/migrations
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	migrationDir := filepath.Join(basepath, "../db/migrations")
	err = goose.UpContext(ctx, dbx.DB, migrationDir)
	if err != nil {
		panic(err)
	}
	return err
}

func GetPgContainer() (*postgres.PostgresContainer, error) {
	if pgContainer == nil {
		err := setupPgContainer()
		if err != nil {
			return nil, err
		}
	}
	return pgContainer, nil
}

func GetDbInstance() (*sqlx.DB, error) {
	if dbx == nil {
		err := setupPgContainer()
		if err != nil {
			return nil, err
		}
	}
	return dbx, nil
}
