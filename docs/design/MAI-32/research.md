# Research: MAI-32 — Simplify XP to Per-Task Model

**Ticket:** MAI-32

**Date:** 2026-03-12 (reconstructed from commit history)

## Scope

Explored how XP was previously tracked (per-category `xp_per_goal` field) and what needed to change to make every goal completion worth exactly 1 XP.

## Codebase Findings

### Relevant Files

- `backend/internal/db/migrations/20260208215034_drop_xp_per_goal.sql` — migration dropping `xp_per_goal` from `goal_categories`
- `backend/internal/db/generated/goal_categories.sql.go` — all category queries had `xp_per_goal` removed
- `backend/internal/db/generated/models.go` — `GoalCategory` struct removed `XpPerGoal int32` field
- `backend/internal/goals/service/service.go` — removed `XPPerGoalMax` constant and `xpPerGoal` param from `CreateGoalCategory`
- `backend/internal/users/service/events.go` — `handleGoalUpdatedEvent` previously used `eventData.Xp`; now hardcodes `+1`
- `backend/internal/events/event_types.go` — removed `Xp` field from goal update event payload
- `backend/seeds/data/levels.json` — simplified from varying large XP values (100–1000+) to small uniform values (5–7 tasks per level)
- `frontend/src/features/levels/components/ProgressBar.vue` — updated label from "XP" to "Tasks"
- `frontend/src/utils/schemas.ts` — removed `xp_per_goal` from `GoalCategorySchema`
- `frontend/src/utils/types.ts` — removed `xp_per_goal` from `TGoalCategory`

### Existing Patterns

- XP was previously per-category (each category had its own `xp_per_goal` integer)
- Goal completion fired a `GoalUpdated` event with an `Xp` field carrying the category's XP value
- The user service event handler read that XP and added it to `user.Xp`
- Level-up logic: `if newXp >= level.LevelUpXp { newXp %= level.LevelUpXp; newLevel++ }`

### Constraints & Gotchas

- Had to update all sqlc-generated code after removing the column (re-ran `make generate`)
- Levels seed had to be completely restructured — old seeds had 100 levels with XP thresholds of 100–1000+; new seeds use small values (5–7) so that 1 task = 1 meaningful step toward level-up
- ProgressBar percentage calculation (`xp / level_up_xp * 100`) still works correctly with the new model

## Open Questions

None — implementation is complete.
