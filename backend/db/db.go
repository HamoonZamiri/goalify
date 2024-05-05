package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New(dbname string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("user=postgres dbname=%s sslmode=disable", dbname)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
