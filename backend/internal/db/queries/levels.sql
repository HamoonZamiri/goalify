-- name: UpsertLevel :exec
INSERT INTO levels (id, level_up_xp, cash_reward, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
	level_up_xp = EXCLUDED.level_up_xp,
	cash_reward = EXCLUDED.cash_reward,
	updated_at = NOW();

-- name: GetAllLevels :many
SELECT * FROM levels ORDER BY id;
