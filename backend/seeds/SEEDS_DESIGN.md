# Seeds Package Design

## Overview

The `seeds/` package provides idempotent database seeding for static/reference data that needs to be versioned in Git and synchronized to the database on app startup. This is **not** for test data or mock data - this is for production static data like levels, eras, stories, and other game configuration.

## Why Seeds vs Migrations?

**Migrations** are for schema changes (CREATE TABLE, ALTER COLUMN, etc.)

**Seeds** are for static data that:
- Changes frequently during development
- Should be versioned as "current state" not "change history"
- Needs to be easily editable (JSON/CSV > SQL)
- Can be safely re-run on every deploy (idempotent upserts)

## Architecture

### Directory Structure

```
backend/seeds/
├── data/                        # Embedded data files
│   ├── levels.json             # Level progression config
│   ├── eras.json               # Historical eras/time periods
│   └── stories/                # Story content by era
│       ├── ww1-spy/
│       │   ├── nodes.json      # Story beats/scenes
│       │   └── choices.json    # Player choices & branches
│       ├── ww2-resistance/
│       │   ├── nodes.json
│       │   └── choices.json
│       └── ...
├── seeds.go                    # Main orchestrator
├── levels.go                   # Level seeding logic
├── levels_test.go             # Validates levels.json format
├── eras.go                    # Era seeding logic
├── eras_test.go               # Validates eras.json format
├── stories.go                 # Story seeding logic
└── stories_test.go            # Validates story data format
```

### Integration with App Startup

Seeds run automatically on app startup (in `cmd/app/app.go`):

```go
func Run(ctx context.Context) error {
	// ... setup pgxPool, run migrations ...
	
	// Run seeds after migrations
	if err := seeds.Run(ctx, pgxPool); err != nil {
		return fmt.Errorf("failed to seed database: %w", err)
	}
	
	// ... start HTTP server ...
}
```

## Data Models

### Required Schema Changes

Add these tables via migration:

```sql
-- Eras/time periods (WW1, WW2, etc.)
CREATE TABLE eras (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	slug VARCHAR(100) NOT NULL UNIQUE,  -- "ww1-spy", "ww2-resistance"
	title VARCHAR(255) NOT NULL,        -- "The Great War - Behind Enemy Lines"
	year INT NOT NULL,                  -- 1916
	description TEXT,
	unlock_level INT REFERENCES levels(id),
	image_url VARCHAR(255),
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW()
);

-- Story nodes (scenes/beats in a story)
CREATE TABLE story_nodes (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	slug VARCHAR(100) NOT NULL,         -- "mission-briefing"
	era_id UUID REFERENCES eras(id) ON DELETE CASCADE,
	title TEXT NOT NULL,
	content TEXT NOT NULL,              -- Main story text
	order_index INT NOT NULL,           -- Position in story flow
	parent_node_id UUID REFERENCES story_nodes(id), -- For branching
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	UNIQUE(era_id, slug)
);

-- Story choices (player decisions)
CREATE TABLE story_choices (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	slug VARCHAR(100) NOT NULL,         -- "accept-mission"
	node_id UUID REFERENCES story_nodes(id) ON DELETE CASCADE,
	choice_text TEXT NOT NULL,          -- "Accept the mission immediately"
	next_node_id UUID REFERENCES story_nodes(id),
	skill_requirement JSONB,            -- {"stealth": 5, "charisma": 3}
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	UNIQUE(node_id, slug)
);

-- User progress through stories
CREATE TABLE user_story_progress (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	era_id UUID REFERENCES eras(id) ON DELETE CASCADE,
	current_node_id UUID REFERENCES story_nodes(id),
	choices_made JSONB,                 -- Track path taken
	completed BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	UNIQUE(user_id, era_id)
);

-- User skills (for unlock conditions)
CREATE TABLE user_skills (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	skill_name VARCHAR(50) NOT NULL,    -- "stealth", "combat", "charisma"
	skill_level INT DEFAULT 1,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	UNIQUE(user_id, skill_name)
);
```

### sqlc Queries Needed

Add to `backend/db/query.sql`:

```sql
-- name: UpsertLevel :exec
INSERT INTO levels (id, level_up_xp, cash_reward, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
	level_up_xp = EXCLUDED.level_up_xp,
	cash_reward = EXCLUDED.cash_reward,
	updated_at = NOW();

-- name: UpsertEra :exec
INSERT INTO eras (slug, title, year, description, unlock_level, image_url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
ON CONFLICT (slug) DO UPDATE SET
	title = EXCLUDED.title,
	year = EXCLUDED.year,
	description = EXCLUDED.description,
	unlock_level = EXCLUDED.unlock_level,
	image_url = EXCLUDED.image_url,
	updated_at = NOW();

-- name: GetEraBySlug :one
SELECT * FROM eras WHERE slug = $1;

-- name: UpsertStoryNode :exec
INSERT INTO story_nodes (slug, era_id, title, content, order_index, parent_node_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
ON CONFLICT (era_id, slug) DO UPDATE SET
	title = EXCLUDED.title,
	content = EXCLUDED.content,
	order_index = EXCLUDED.order_index,
	parent_node_id = EXCLUDED.parent_node_id,
	updated_at = NOW();

-- name: GetStoryNodeBySlug :one
SELECT * FROM story_nodes WHERE era_id = $1 AND slug = $2;

-- name: UpsertStoryChoice :exec
INSERT INTO story_choices (slug, node_id, choice_text, next_node_id, skill_requirement, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
ON CONFLICT (node_id, slug) DO UPDATE SET
	choice_text = EXCLUDED.choice_text,
	next_node_id = EXCLUDED.next_node_id,
	skill_requirement = EXCLUDED.skill_requirement,
	updated_at = NOW();
```

## Data Formats

### levels.json

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
	},
	{
		"id": 5,
		"xp": 500,
		"cash": 200
	}
]
```

**Progression curve recommendations**:
- **Milestone-based**: Big jumps at story unlock levels (5, 10, 15, etc.)
- **Exponential**: Later levels require significantly more XP (avoids endgame grind)
- **Tiered**: Different eras have different XP curves

Example progression:
- Levels 1-4: 100, 150, 200, 250 (quick tutorial)
- Level 5: 500 (first story unlock - milestone)
- Levels 6-9: 300, 350, 400, 450
- Level 10: 800 (second story unlock)
- And so on...

### eras.json

```json
[
	{
		"slug": "ww1-spy",
		"title": "The Great War - Behind Enemy Lines",
		"year": 1916,
		"description": "Navigate the trenches of WWI as an American spy infiltrating German command.",
		"unlock_level": 5,
		"image_url": "/images/eras/ww1-spy.jpg"
	},
	{
		"slug": "ww2-resistance",
		"title": "WWII - French Resistance",
		"year": 1943,
		"description": "Join the French Resistance and sabotage Nazi operations in occupied Paris.",
		"unlock_level": 10,
		"image_url": "/images/eras/ww2-resistance.jpg"
	}
]
```

### stories/ww1-spy/nodes.json

```json
[
	{
		"slug": "mission-briefing",
		"title": "The Mission Briefing",
		"content": "Colonel Hayes slides a dossier across the desk. 'Intelligence suggests the Germans are planning a major offensive. We need someone behind enemy lines.'",
		"order_index": 1,
		"parent_node_slug": null
	},
	{
		"slug": "accept-mission",
		"title": "Into the Trenches",
		"content": "You accept without hesitation. Within hours, you're on a supply truck headed for the front lines, disguised as a German officer.",
		"order_index": 2,
		"parent_node_slug": "mission-briefing"
	},
	{
		"slug": "question-extraction",
		"title": "Planning the Exit",
		"content": "You press Hayes for details on extraction. He reluctantly reveals a network of safe houses, but warns the route is dangerous.",
		"order_index": 2,
		"parent_node_slug": "mission-briefing"
	}
]
```

### stories/ww1-spy/choices.json

```json
[
	{
		"slug": "accept-immediately",
		"node_slug": "mission-briefing",
		"choice_text": "Accept the mission immediately",
		"next_node_slug": "accept-mission",
		"skill_requirement": null
	},
	{
		"slug": "ask-extraction",
		"node_slug": "mission-briefing",
		"choice_text": "Ask about extraction plans (Intelligence 3+)",
		"next_node_slug": "question-extraction",
		"skill_requirement": {
			"intelligence": 3
		}
	}
]
```

## Implementation Pattern

### seeds.go (Orchestrator)

```go
package seeds

import (
	"context"
	"fmt"
	"log/slog"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

// Run executes all seed functions in order
func Run(ctx context.Context, pool *pgxpool.Pool) error {
	seeders := []struct {
		name string
		fn   func(context.Context, *pgxpool.Pool) error
	}{
		{"levels", SeedLevels},
		{"eras", SeedEras},
		{"stories", SeedStories},
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

### levels.go (Example Implementation)

```go
package seeds

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	
	"goalify/internal/db/generated"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

### levels_test.go (Validation)

```go
package seeds

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLevelsJSONValid(t *testing.T) {
	levels, err := loadLevels()
	require.NoError(t, err, "levels.json should parse without errors")
	require.Len(t, levels, 100, "should have exactly 100 levels")
	
	// Validate business rules
	for i, lvl := range levels {
		assert.Equal(t, int32(i+1), lvl.ID, "level IDs should be sequential starting at 1")
		assert.Greater(t, lvl.LevelUpXP, int32(0), "level_up_xp must be positive")
		assert.Greater(t, lvl.CashReward, int32(0), "cash_reward must be positive")
	}
	
	// Ensure progression feels right
	assert.Greater(t, levels[4].LevelUpXP, levels[0].LevelUpXP, "level 5 should require more XP than level 1")
	assert.Greater(t, levels[99].LevelUpXP, levels[4].LevelUpXP*2, "endgame should be significantly harder")
}
```

## Key Design Principles

### 1. Idempotency

All seed operations use `ON CONFLICT ... DO UPDATE` so they can be run repeatedly without errors or duplicate data. This means:
- Safe to run on every app startup
- Safe to re-run after editing JSON files
- No "has this been seeded?" checks needed

### 2. Fail-Fast

If JSON parsing fails or data validation fails, the app won't start. This is intentional:
- Prevents serving broken/stale data
- Catches issues in CI/tests before production
- Forces data quality

### 3. Transactional

Each seed function runs in a transaction:
- All-or-nothing semantics
- Rollback on any error
- Database stays consistent

### 4. Testable

Every seed file has a corresponding test that validates:
- JSON parses correctly
- Required fields are present
- Business rules are satisfied (e.g., levels are sequential)
- Data relationships are valid (e.g., era unlock_level exists in levels table)

### 5. Use sqlc

Leverage existing sqlc queries instead of raw SQL:
- Type safety
- Consistent with rest of codebase
- Easy to test and mock

## Workflow

### Development Workflow

1. Edit JSON file (e.g., add level 101 to `levels.json`)
2. Run tests: `make test` (validates JSON format)
3. Restart app: `make dev`
4. Seeds run automatically on startup
5. Changes reflected in database

### Adding New Static Data

1. Create migration for new table (if needed)
2. Add sqlc queries for upsert operations
3. Create JSON file in `seeds/data/`
4. Create `seeds/{feature}.go` with load + seed functions
5. Create `seeds/{feature}_test.go` with validation
6. Add to `seeds.go` orchestrator
7. Run `make dev` - seeds run automatically

### Story Content Workflow

For non-engineers (writers, designers):

1. Edit story JSON files directly
2. Commit to Git
3. Deploy triggers automatic seed on startup
4. Content live in production

**Future enhancement**: Build web-based story editor so non-technical team members can create/edit stories through a UI.

## Performance Considerations

### Embedded Files

Using `//go:embed` means:
- Data files are baked into the binary at compile time
- No file I/O at runtime (faster, more reliable)
- Single binary deployment (no need to ship JSON separately)

### Batch Operations

Seed functions should use:
- Prepared statements (where applicable)
- Single transaction per seed type
- Avoid N+1 queries

### Startup Time

Seeds add ~100-500ms to app startup depending on data volume:
- Acceptable for development
- Acceptable for production deploys
- If this becomes an issue, add `-skip-seeds` flag for local dev

## Error Handling

Seeds should:
- Return descriptive errors with context (which file, which record)
- Fail fast (don't continue if one seed fails)
- Log progress (`slog.Info("Seeding levels...")`)
- Rollback transactions on error

Example:
```go
if err := queries.UpsertLevel(ctx, params); err != nil {
	return fmt.Errorf("upsert level %d from levels.json: %w", lvl.ID, err)
}
```

## Future Enhancements

### 1. Seed Versioning

Track which seed version was last applied:
```sql
CREATE TABLE seed_versions (
	name VARCHAR(100) PRIMARY KEY,
	version VARCHAR(50) NOT NULL,
	applied_at TIMESTAMP DEFAULT NOW()
);
```

Skip seeding if version unchanged (optimization).

### 2. Seed Rollback

Support for rolling back seeds (rare, but useful for emergency fixes):
```go
func RollbackLevels(ctx context.Context, pool *pgxpool.Pool) error
```

### 3. Web-Based Story Editor

Admin panel for creating/editing stories without touching JSON:
- WYSIWYG story editor
- Visual branching diagram
- Skill requirement UI
- Export to JSON for version control

### 4. Dry-Run Mode

Preview what would be seeded without committing:
```bash
make seed-dry-run
```

### 5. Selective Seeding

Seed only specific resources:
```bash
make seed SEEDS=levels,eras
```

## Migration from Current System

Current state: Levels are seeded via SQL in migration file.

Migration plan:
1. Create `seeds/` package with structure above
2. Generate `levels.json` from current database state:
   ```sql
   SELECT json_agg(json_build_object('id', id, 'xp', level_up_xp, 'cash', cash_reward))
   FROM levels ORDER BY id;
   ```
3. Implement `seeds/levels.go` with upsert logic
4. Test: verify idempotent re-runs work
5. Remove SQL seed logic from migration (or keep as fallback)
6. Deploy: seeds run on startup, overwrite existing data (safe due to ON CONFLICT)
7. Expand to eras, stories, etc.

## Summary

The `seeds/` package provides:
- ✅ Version-controlled static data in Git
- ✅ Easy editing (JSON > SQL)
- ✅ Automatic synchronization on deploy
- ✅ Type-safe using sqlc
- ✅ Testable and validated
- ✅ Idempotent and transactional
- ✅ Fail-fast error handling

This approach scales well for RPG-style game content (levels, eras, stories, items, etc.) while keeping the data maintainable and reviewable.
