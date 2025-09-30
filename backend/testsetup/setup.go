package testsetup

import (
	"context"
	"goalify/config"
	"goalify/db"
	"path/filepath"
	"runtime"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	pgContainer *postgres.PostgresContainer
	pgxPool     *pgxpool.Pool
)

func setupPgContainer() error {
	configService := config.GetConfig()
	ctx := context.Background()
	if pgContainer != nil {
		return nil
	}

	var err error
	pgContainer, err = postgres.Run(ctx, "docker.io/postgres:16-alpine",
		postgres.WithDatabase(configService.TestDBName),
		postgres.WithUsername(configService.TestDBUser),
		postgres.WithPassword(configService.TestDBPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return err
	}

	pgxPool, err = db.NewPgxPoolWithConnString(ctx, connStr)
	if err != nil {
		panic(err)
	}

	// using goose run migrations from db/migrations
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	migrationDir := filepath.Join(basepath, "../db/migrations")
	err = goose.UpContext(ctx, stdlib.OpenDBFromPool(pgxPool), migrationDir)
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

func GetPgxPool() (*pgxpool.Pool, error) {
	if pgxPool == nil {
		err := setupPgContainer()
		if err != nil {
			return nil, err
		}
	}
	return pgxPool, nil
}