package db

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New(dbname, user, password string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func BuildUpdateQuery(table string, updates map[string]any, id uuid.UUID) (string, []any) {
	setClauses := []string{}
	args := []any{}

	i := 1

	for column, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		i++
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", table, strings.Join(setClauses, ", "), i)
	args = append(args, id)

	return query, args
}

func BuildSelectQuery(table string, filters map[string]any) (string, []any) {
	whereClauses := []string{}
	args := []any{}

	i := 1

	for column, value := range filters {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		i++
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, strings.Join(whereClauses, " AND "))

	return query, args
}
