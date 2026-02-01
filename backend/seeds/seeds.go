// Package seeds provides database seeding functionality for static data.
// Seeds are idempotent and can be safely run on every app startup.
package seeds

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Run executes all seed functions in order.
// Returns error if any seed fails.
func Run(ctx context.Context, pool *pgxpool.Pool) error {
	seeders := []struct {
		fn   func(context.Context, *pgxpool.Pool) error
		name string
	}{
		{SeedLevels, "levels"},
	}

	for _, s := range seeders {
		slog.Info("Seeding " + s.name + "...")
		if err := s.fn(ctx, pool); err != nil {
			return fmt.Errorf("seed %s: %w", s.name, err)
		}
	}

	slog.Info("All seeds completed successfully")
	return nil
}
