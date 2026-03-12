-- +goose Up
ALTER TABLE goal_categories DROP COLUMN xp_per_goal;
UPDATE users SET xp = 0;

-- +goose Down
ALTER TABLE goal_categories ADD COLUMN xp_per_goal INTEGER NOT NULL DEFAULT 1;
