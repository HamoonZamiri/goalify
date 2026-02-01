package seeds

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	db "goalify/internal/db/generated"
)

//go:embed data/levels.json
var levelsJSON []byte

// LevelSeed represents a single level's data in the JSON file.
type LevelSeed struct {
	ID         int32 `json:"id"`
	LevelUpXP  int32 `json:"xp"`
	CashReward int32 `json:"cash"`
}

// loadLevels parses the embedded levels.json file.
func loadLevels() ([]LevelSeed, error) {
	var levels []LevelSeed
	if err := json.Unmarshal(levelsJSON, &levels); err != nil {
		return nil, fmt.Errorf("parse levels.json: %w", err)
	}
	return levels, nil
}

// SeedLevels upserts all levels from levels.json into the database.
// Uses a transaction to ensure atomicity - all levels are seeded or none are.
func SeedLevels(ctx context.Context, pool *pgxpool.Pool) error {
	levels, err := loadLevels()
	if err != nil {
		return err
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("failed to rollback transaction", "error", err)
		}
	}()

	queries := db.New(tx)

	for _, lvl := range levels {
		if err := queries.UpsertLevel(ctx, db.UpsertLevelParams{
			ID:         lvl.ID,
			LevelUpXp:  lvl.LevelUpXP,
			CashReward: lvl.CashReward,
		}); err != nil {
			return fmt.Errorf("upsert level %d: %w", lvl.ID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
