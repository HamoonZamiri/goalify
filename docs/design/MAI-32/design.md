# Feature: Simplify XP to Per-Task Model

**Status:** review

**Ticket:** MAI-32

**Branch:** hamoondev/mai-32-refactor-simplify-xp-to-per-task-model

**Created:** 2026-03-08 (reconstructed from commit history)

## Goal

Replace per-category variable XP with a flat 1 XP per completed task, so users always know exactly how many tasks bring them to the next level.

## Context

The old model assigned a configurable `xp_per_goal` to each goal category. This made it hard to reason about level progression — the XP required per level was arbitrary and the UX was confusing. The new model: completing any task = +1 XP, levels have small `level_up_xp` thresholds (5–7 tasks), so progress is intuitive.

## Out of Scope

- Changing the level-up mechanics (modulo carry-over logic stays)
- Rewarding different amounts for different task types
- Retroactive XP recalculation for existing users

## Approach

1. Drop `xp_per_goal` column from `goal_categories` via migration
2. Remove the field from all DB queries, generated code, entities, handlers, and schemas
3. Hardcode `+1` XP in `users/service/events.go` when a goal is marked complete
4. Simplify levels seed data to small `level_up_xp` values
5. Update frontend label from "XP" to "Tasks"

## Tasks

- [x] Write and apply migration `drop_xp_per_goal`
- [x] Regenerate sqlc (`make generate`)
- [x] Remove `xpPerGoal` from `CreateGoalCategory` service method and handler
- [x] Remove `Xp` from `GoalUpdated` event payload
- [x] Hardcode `+1` in `users/service/events.go`
- [x] Simplify `backend/seeds/data/levels.json`
- [x] Remove `xp_per_goal` from frontend schemas and types
- [x] Remove `xp_per_goal` from create/edit category forms
- [x] Update ProgressBar label to "Tasks"
- [x] Update all tests
- [ ] Open PR

## Acceptance Criteria

- [x] A user should see that each task increases their progress to the next level by one
- [x] The backend should seamlessly understand the new model with each goal completion always increasing by 1 XP

## Open Questions

None.

## Decisions Log

- **2026-03-08** — Kept XP as the internal unit but locked it to 1 per task. The frontend labels it "Tasks" so the UX reads as task-count-to-level-up without exposing the XP concept.
- **2026-03-08** — Levels seed simplified: levels 1–10 require 5 tasks to level up, levels 11+ require 7. This gives a gentle difficulty ramp that's still intuitive.

## Session Log

- **2026-03-12** — Design doc reconstructed from commit `a3eafda` (merged 2026-03-08). All backend and frontend changes are committed. One uncommitted change: capitalizing "tasks" → "Tasks" in `ProgressBar.vue`. Backend: migration applied, sqlc regenerated, event handler hardcodes +1 XP, levels seed simplified. Frontend: `xp_per_goal` removed from schemas/types/forms, ProgressBar label updated. Tests updated. PR not yet opened.
