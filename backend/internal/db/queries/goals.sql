-- name: CreateGoal :one
INSERT INTO goals (title, description, user_id, category_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateGoalStatus :one
UPDATE goals
SET status = $1
WHERE id = $2 AND user_id = $3
RETURNING *;

-- name: GetGoalsByUserId :many
SELECT * FROM goals WHERE user_id = $1;

-- name: GetGoalById :one
SELECT * FROM goals WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: UpdateGoalById :one
UPDATE goals
SET title = coalesce(sqlc.narg('title'), title),
    description = coalesce(sqlc.narg('description'), description),
    status = coalesce(sqlc.narg('status'), status),
    category_id = coalesce(sqlc.narg('category_id'), category_id)
WHERE id = sqlc.arg('id') AND user_id = sqlc.arg('user_id')
RETURNING *;

-- name: DeleteGoalById :execrows
DELETE FROM goals WHERE id = $1 AND user_id = $2;

-- name: ResetGoalsByCategory :exec
UPDATE goals
SET status = 'not_complete'
WHERE category_id = $1 AND user_id = $2;
