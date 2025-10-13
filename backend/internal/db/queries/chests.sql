-- Chest Operations
-- name: CreateChest :one
INSERT INTO chests (type, description, price)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetChestById :one
SELECT * FROM chests WHERE id = $1;

-- name: GetAllChests :many
SELECT * FROM chests ORDER BY price ASC;

-- name: UpdateChestById :one
UPDATE chests
SET type = coalesce(sqlc.narg('type'), type),
    description = coalesce(sqlc.narg('description'), description),
    price = coalesce(sqlc.narg('price'), price)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteChestById :exec
DELETE FROM chests WHERE id = $1;

-- Chest Item Operations
-- name: CreateChestItem :one
INSERT INTO chest_items (image_url, title, rarity, price)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetChestItemById :one
SELECT * FROM chest_items WHERE id = $1;

-- name: GetAllChestItems :many
SELECT * FROM chest_items ORDER BY created_at DESC;

-- name: UpdateChestItemById :one
UPDATE chest_items
SET image_url = coalesce(sqlc.narg('image_url'), image_url),
    title = coalesce(sqlc.narg('title'), title),
    rarity = coalesce(sqlc.narg('rarity'), rarity),
    price = coalesce(sqlc.narg('price'), price)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteChestItemById :exec
DELETE FROM chest_items WHERE id = $1;

-- Chest Item Drop Rate Operations
-- name: CreateChestItemDropRate :one
INSERT INTO chest_item_drop_rates (item_id, chest_id, drop_rate)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetDropRatesByChestId :many
SELECT * FROM chest_item_drop_rates WHERE chest_id = $1;

-- name: GetDropRatesByItemId :many
SELECT * FROM chest_item_drop_rates WHERE item_id = $1;

-- name: UpdateDropRate :one
UPDATE chest_item_drop_rates
SET drop_rate = $2
WHERE id = $1
RETURNING *;

-- name: DeleteDropRate :exec
DELETE FROM chest_item_drop_rates WHERE id = $1;

-- User Chest Operations
-- name: CreateUserChest :one
INSERT INTO user_chests (user_id, chest_id, quantity_owned)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserChestsByUserId :many
SELECT * FROM user_chests WHERE user_id = $1;

-- name: GetUserChestByUserIdAndChestId :one
SELECT * FROM user_chests WHERE user_id = $1 AND chest_id = $2;

-- name: UpdateUserChestQuantity :one
UPDATE user_chests
SET quantity_owned = $2
WHERE user_id = $1 AND chest_id = $3
RETURNING *;

-- name: DeleteUserChest :exec
DELETE FROM user_chests WHERE user_id = $1 AND chest_id = $2;

-- User Item Operations
-- name: CreateUserItem :one
INSERT INTO user_items (user_id, item_id, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserItemsByUserId :many
SELECT * FROM user_items WHERE user_id = $1;

-- name: GetUserItemById :one
SELECT * FROM user_items WHERE id = $1;

-- name: UpdateUserItemStatus :one
UPDATE user_items
SET status = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUserItem :exec
DELETE FROM user_items WHERE id = $1;