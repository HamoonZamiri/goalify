// Package db handles database connection and sqlc generated files
package db

import (
	"context"
	"fmt"
	"goalify/pkg/options"
	"math"

	sqlcdb "goalify/internal/db/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPoolWithConnString(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.New(ctx, config.ConnString())
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func NewPgx(dbname, user, password, host string) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=disable",
		user,
		password,
		host,
		dbname,
	)
	return NewPgxPoolWithConnString(context.Background(), connStr)
}

func UUIDToPgxUUID(uuid uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: uuid, Valid: true}
}

func StringToPgxText(str string) pgtype.Text {
	return pgtype.Text{String: str, Valid: true}
}

func AnyToString(v any) (str string, ok bool) {
	if v == nil {
		return "", false
	}
	str, ok = v.(string)
	return str, ok
}

func AnyToInt(v any) (i int, ok bool) {
	if v == nil {
		return 0, false
	}
	i, ok = v.(int)
	return i, ok
}

func IntToPgxInt4(i int) (pgtype.Int4, error) {
	if i > math.MaxInt32 || i < math.MinInt32 {
		return pgtype.Int4{Valid: false}, fmt.Errorf("integer overflow: %d exceeds int32 range", i)
	}
	return pgtype.Int4{Int32: int32(i), Valid: true}, nil
}

func StringToGoalStatus(status string) (sqlcdb.NullGoalStatus, error) {
	goalStatus := sqlcdb.GoalStatus(status)
	if goalStatus != sqlcdb.GoalStatusComplete && goalStatus != sqlcdb.GoalStatusNotComplete {
		return sqlcdb.NullGoalStatus{Valid: false}, fmt.Errorf("invalid goal status: %s", status)
	}
	return sqlcdb.NullGoalStatus{
		GoalStatus: goalStatus,
		Valid:      true,
	}, nil
}

func OptionStringToPgxText(opt options.Option[string]) pgtype.Text {
	if opt.IsPresent() {
		return StringToPgxText(opt.ValueOrZero())
	}
	return pgtype.Text{}
}

func OptionUUIDToPgxUUID(opt options.Option[uuid.UUID]) pgtype.UUID {
	if opt.IsPresent() {
		return UUIDToPgxUUID(opt.ValueOrZero())
	}
	return pgtype.UUID{}
}

func OptionIntToPgxInt4(opt options.Option[int]) (pgtype.Int4, error) {
	if opt.IsPresent() {
		return IntToPgxInt4(opt.ValueOrZero())
	}
	return pgtype.Int4{}, nil
}
