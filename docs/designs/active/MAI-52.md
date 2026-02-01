# Feature: Seeds Package for Static Data

**Status:** review

**Ticket:** MAI-52

**Branch:** hamoondev/mai-52-refactor-run-migrations-and-seed-directory-with-goose-on-app

**Created:** 2026-01-28

Status definitions:
  - draft: Planning/designing, not yet approved for implementation
  - approved: Design approved, ready to start implementation
  - in-progress: Actively being worked on
  - review: Implementation complete, undergoing self-review or awaiting PR review
  - done: Merged, acceptance criteria met, moved to completed/

## Goal
Refactor app startup to run migrations (via goose Go SDK) and seeds (via new `backend/seeds/` package) in code, eliminating the need for separate migration commands and enabling idempotent static data syncing from JSON.

## Context
Currently, levels are seeded via SQL in migration files. This approach has limitations:
- Hard to edit (SQL vs JSON)
- Not easily reviewable (schema changes mixed with data changes)
- Can't be re-run safely (need manual conflict resolution)

Since we're running on a single Hetzner node with Coolify, the simplest approach is to run migrations and seeds in code on app startup. This provides:
- Version-controlled static data in Git
- Easy editing (JSON > SQL)
- Automatic synchronization on deploy
- Type-safe using sqlc
- Idempotent and transactional operations

## Out of Scope
- Story features (eras, story_nodes, story_choices tables)
- Database modeling for new story-related tables
- Web-based content editor
- Seed versioning/rollback mechanisms
- Dry-run or selective seeding modes

This doc focuses solely on the seeds package infrastructure for the **existing levels table**.

## Approach

### Directory Structure
```
backend/seeds/
├── data/
│   └── levels.json          # Level progression config (embedded)
├── seeds.go                 # Main orchestrator
├── levels.go                # Level seeding logic
└── levels_test.go           # Validates levels.json format
```

### Implementation Pattern

1. **Embedded Data Files** - Use `//go:embed` to bake JSON into binary
2. **Idempotent Upserts** - Use existing `UpsertLevel` sqlc query with `ON CONFLICT ... DO UPDATE`
3. **Fail-Fast** - App won't start if JSON parsing or seeding fails
4. **Transactional** - Each seed function runs in a single transaction
5. **Testable** - Validate JSON structure and business rules in tests

### Integration

Both migrations and seeds run in `cmd/app/app.go` on startup using the goose Go SDK:

```go
import (
    "github.com/pressly/goose/v3"
    "goalify/internal/db/migrations" // embed.FS for migrations
    "goalify/seeds"
)

func Run(ctx context.Context) error {
    // ... setup pgxPool ...

    // Run migrations using goose SDK
    goose.SetBaseFS(migrations.EmbedMigrations)
    if err := goose.Up(pgxPool.Config().ConnConfig.Database, "postgres"); err != nil {
        return fmt.Errorf("failed to run migrations: %w", err)
    }

    // Run seeds after migrations
    if err := seeds.Run(ctx, pgxPool); err != nil {
        return fmt.Errorf("failed to seed database: %w", err)
    }

    // ... start HTTP server ...
}
```

**Why run migrations in Go?**
- Eliminates separate goose CLI installation in Dockerfile
- Single binary deployment (migrations embedded via `//go:embed`)
- Consistent startup flow: migrations → seeds → server
- Simpler for Hetzner/Coolify single-node deployment

**Dockerfile changes:**
- Remove goose CLI installation/commands
- Migrations are embedded in binary (already copying migration files)
- Single `CMD ["./goalify"]` runs everything

### Data Format (levels.json)
```json
[
  {
    "id": 1,
    "xp": 100,
    "cash": 100
  },
  {
    "id": 2,
    "xp": 150,
    "cash": 110
  }
]
```

### Required sqlc Query

The seed functions use sqlc-generated code for type safety and consistency. Add this query to `backend/db/query.sql`:

```sql
-- name: UpsertLevel :exec
INSERT INTO levels (id, level_up_xp, cash_reward, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
    level_up_xp = EXCLUDED.level_up_xp,
    cash_reward = EXCLUDED.cash_reward,
    updated_at = NOW();
```

**How it works:**
- `EXCLUDED` is a PostgreSQL keyword referencing the values from the attempted INSERT
- Makes the operation idempotent - can run on every app startup
- `ON CONFLICT (id)` triggers when level with that ID already exists
- Updates existing row with new values from JSON

**Workflow:**
1. Write SQL query in `backend/db/query.sql`
2. Run `make sqlc` to generate Go code
3. Seed functions call `queries.UpsertLevel(ctx, db.UpsertLevelParams{...})`

### Key Files

**seeds.go** - Orchestrator that runs all seed functions in order
```go
func Run(ctx context.Context, pool *pgxpool.Pool) error {
    seeders := []struct {
        name string
        fn   func(context.Context, *pgxpool.Pool) error
    }{
        {"levels", SeedLevels},
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
```

**levels.go** - Load JSON, upsert to database
```go
//go:embed data/levels.json
var levelsJSON []byte

type LevelSeed struct {
    ID         int32 `json:"id"`
    LevelUpXP  int32 `json:"xp"`
    CashReward int32 `json:"cash"`
}

func loadLevels() ([]LevelSeed, error) {
    var levels []LevelSeed
    if err := json.Unmarshal(levelsJSON, &levels); err != nil {
        return nil, fmt.Errorf("parse levels.json: %w", err)
    }
    return levels, nil
}

func SeedLevels(ctx context.Context, pool *pgxpool.Pool) error {
    levels, err := loadLevels()
    if err != nil {
        return err
    }

    tx, err := pool.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

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

    return tx.Commit(ctx)
}
```

**levels_test.go** - Validate JSON structure and business rules
```go
func TestLevelsJSONValid(t *testing.T) {
    levels, err := loadLevels()
    require.NoError(t, err)
    require.NotEmpty(t, levels)

    for i, lvl := range levels {
        assert.Equal(t, int32(i+1), lvl.ID, "level IDs should be sequential")
        assert.Greater(t, lvl.LevelUpXP, int32(0))
        assert.Greater(t, lvl.CashReward, int32(0))
    }
}
```

### Progression Curve Recommendations
- **Milestone-based**: Big XP jumps at story unlock levels (5, 10, 15, etc.)
- **Exponential**: Later levels require significantly more XP
- **Example**:
  - Levels 1-4: 100, 150, 200, 250 (tutorial)
  - Level 5: 500 (milestone)
  - Levels 6-9: 300, 350, 400, 450
  - Level 10: 800 (next milestone)

## Tasks
- [x] Add `UpsertLevel` sqlc query to `backend/db/query.sql`
- [x] Run `make sqlc` to generate Go code
- [x] Refactor `cmd/app/app.go` to run migrations using goose Go SDK (already using SDK)
- [x] Create `backend/seeds/` directory structure
- [x] Generate `data/levels.json` from current database state
- [x] Implement `seeds.go` orchestrator
- [x] Implement `levels.go` with `loadLevels()` and `SeedLevels()`
- [x] Implement `levels_test.go` validation tests
- [x] Integrate seeds into `cmd/app/app.go` after migrations
- [x] Update Dockerfile to remove migrations COPY (embedded via go:embed)
- [x] Verify migrations + seeds run successfully on app startup locally
- [ ] Update documentation/README if needed

## Acceptance Criteria
- [x] Migrations run automatically on app startup using goose Go SDK
- [x] Seeds run automatically after migrations complete
- [x] `levels.json` successfully loads and validates
- [x] Levels are upserted to database (no duplicates on re-run)
- [x] Tests pass: `make test` validates JSON structure
- [x] App fails to start if JSON is malformed (fail-fast)
- [x] Seed operations are transactional (rollback on error)
- [x] Logs show migration and seeding progress
- [x] Dockerfile no longer requires separate goose CLI installation (already using SDK)
- [x] Single `./goalify` binary runs migrations, seeds, and server

## Open Questions
- Should we remove the SQL seed logic from migrations or keep as fallback?
- Do we want a `-skip-seeds` or `-skip-migrations` flag for local development?
- What's the target progression curve for levels 1-100?
- Should we handle migration rollback scenarios (goose down), or only support up migrations?

## Decisions Log
<!-- Append-only: key choices made -->

### 2026-02-01
- **Kept existing progression curve**: Used current database state (linear +10 XP/cash per level) from migration for levels.json
- **Migrations already using SDK**: `internal/db/migrate.go` was already using goose Go SDK with embedded migrations, so no goose CLI installation needed
- **Query organization**: Created new `internal/db/queries/levels.sql` for level-related queries following existing pattern
- **Removed migrations COPY from Dockerfile**: Since migrations and seeds are embedded via `//go:embed`, no need to copy directories at runtime

## Session Log
<!-- Append-only: progress updates for context recovery -->

### 2026-02-01 - Implementation Complete
- Created `backend/internal/db/queries/levels.sql` with `UpsertLevel` and `GetAllLevels` queries
- Generated sqlc code via `make generate`
- Created `backend/seeds/` package structure:
  - `seeds.go`: Orchestrator with `Run()` function
  - `levels.go`: `loadLevels()` and `SeedLevels()` with embedded JSON
  - `levels_test.go`: Validation tests for JSON structure, progression, and count
  - `data/levels.json`: 100 levels exported from current database
- Integrated seeds into `cmd/app/app.go` after migrations
- Verified locally via `make jqdev`:
  - Migrations run successfully (no migrations to run - idempotent ✓)
  - Seeds run successfully after migrations
  - All 100 levels seeded
  - App starts and listens on port 8080
- All seeds tests pass (3/3)
- All unit tests pass
- Ready for review and PR
