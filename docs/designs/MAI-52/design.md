# Feature: Seeds Package for Static Data

**Status:** done

**Ticket:** MAI-52

**Branch:** hamoondev/mai-52-refactor-run-migrations-and-seed-directory-with-goose-on-app

**Created:** 2026-01-28

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

```sql
-- name: UpsertLevel :exec
INSERT INTO levels (id, level_up_xp, cash_reward, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
    level_up_xp = EXCLUDED.level_up_xp,
    cash_reward = EXCLUDED.cash_reward,
    updated_at = NOW();
```

## Tasks
- [x] Add `UpsertLevel` sqlc query to `backend/db/query.sql`
- [x] Run `make sqlc` to generate Go code
- [x] Refactor `cmd/app/app.go` to run migrations using goose Go SDK
- [x] Create `backend/seeds/` directory structure
- [x] Generate `data/levels.json` from current database state
- [x] Implement `seeds.go` orchestrator
- [x] Implement `levels.go` with `loadLevels()` and `SeedLevels()`
- [x] Implement `levels_test.go` validation tests
- [x] Integrate seeds into `cmd/app/app.go` after migrations
- [x] Update Dockerfile to remove migrations COPY (embedded via go:embed)
- [x] Verify migrations + seeds run successfully on app startup locally

## Acceptance Criteria
- [x] Migrations run automatically on app startup using goose Go SDK
- [x] Seeds run automatically after migrations complete
- [x] `levels.json` successfully loads and validates
- [x] Levels are upserted to database (no duplicates on re-run)
- [x] Tests pass: `make test` validates JSON structure
- [x] App fails to start if JSON is malformed (fail-fast)
- [x] Seed operations are transactional (rollback on error)
- [x] Logs show migration and seeding progress
- [x] Dockerfile no longer requires separate goose CLI installation
- [x] Single `./goalify` binary runs migrations, seeds, and server

## Decisions Log

### 2026-02-01
- **Kept existing progression curve**: Used current database state (linear +10 XP/cash per level) from migration for levels.json
- **Migrations already using SDK**: `internal/db/migrate.go` was already using goose Go SDK with embedded migrations, so no goose CLI installation needed
- **Query organization**: Created new `internal/db/queries/levels.sql` for level-related queries following existing pattern
- **Removed migrations COPY from Dockerfile**: Since migrations and seeds are embedded via `//go:embed`, no need to copy directories at runtime

## Session Log

### 2026-02-01 - Implementation Complete
- Created `backend/internal/db/queries/levels.sql` with `UpsertLevel` and `GetAllLevels` queries
- Generated sqlc code via `make generate`
- Created `backend/seeds/` package structure:
  - `seeds.go`: Orchestrator with `Run()` function
  - `levels.go`: `loadLevels()` and `SeedLevels()` with embedded JSON
  - `levels_test.go`: Validation tests for JSON structure, progression, and count
  - `data/levels.json`: 100 levels exported from current database
- Integrated seeds into `cmd/app/app.go` after migrations
- Verified locally via `make jqdev` — migrations and seeds run successfully, app starts on port 8080
- All seeds tests pass (3/3), all unit tests pass
- PR #72 merged
