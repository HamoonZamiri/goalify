-- name: CreateUser :one
INSERT INTO users (email, password, refresh_token_expiry, level_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: UpdateRefreshToken :one
UPDATE users
    SET refresh_token = $1,
    refresh_token_expiry = $2
    WHERE id = $3 
    RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: DeleteUserById :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUserById :one
UPDATE users
    SET email = coalesce(sqlc.narg('email'), email),
    password = coalesce(sqlc.narg('password'), password),
    refresh_token = coalesce(sqlc.narg('refresh_token'), refresh_token),
    refresh_token_expiry = coalesce(sqlc.narg('refresh_token_expiry'), refresh_token_expiry),
    level_id = coalesce(sqlc.narg('level_id'), level_id),
    xp = coalesce(sqlc.narg('xp'), xp),
    cash_available = coalesce(sqlc.narg('cash_available'), cash_available)
    WHERE id = sqlc.arg('id')
    RETURNING *;

-- name: GetLevelById :one
SELECT * FROM levels WHERE id = $1 LIMIT 1;

