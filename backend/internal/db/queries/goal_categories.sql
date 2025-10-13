-- name: CreateGoalCategory :one
INSERT INTO goal_categories (title, xp_per_goal, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetGoalCategoriesByUserId :many
SELECT * FROM goal_categories WHERE user_id = $1 ORDER BY created_at;

-- name: GetGoalCategoryById :one
SELECT * FROM goal_categories WHERE id = $1 LIMIT 1;

-- name: UpdateGoalCategoryById :one
UPDATE goal_categories
SET title = coalesce(sqlc.narg('title'), title),
    xp_per_goal = coalesce(sqlc.narg('xp_per_goal'), xp_per_goal)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteGoalCategoryById :exec
DELETE FROM goal_categories WHERE id = $1;

-- name: GetGoalCategoriesWithGoalsByUserId :many
SELECT
    gc.id, gc.title, gc.xp_per_goal, gc.user_id, gc.created_at, gc.updated_at,
    g.id as goal_id, g.title as goal_title, g.description, g.status,
    g.created_at as goal_created_at, g.updated_at as goal_updated_at
FROM goal_categories gc
LEFT JOIN goals g ON gc.id = g.category_id
WHERE gc.user_id = $1
ORDER BY gc.created_at;

-- name: GetGoalCategoryWithGoalsById :many
SELECT
    gc.id, gc.title, gc.xp_per_goal, gc.user_id, gc.created_at, gc.updated_at,
    g.id as goal_id, g.title as goal_title, g.description, g.status,
    g.created_at as goal_created_at, g.updated_at as goal_updated_at
FROM goal_categories gc
LEFT JOIN goals g ON gc.id = g.category_id
WHERE gc.id = $1;