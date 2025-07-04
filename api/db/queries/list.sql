-- name: GetList :one
SELECT * FROM lists
WHERE id = $1;

-- name: GetListIdsAccessible :many
SELECT id FROM lists l
WHERE l.user_id = $1 OR id IN (
    SELECT list_id FROM list_shares ls WHERE ls.user_id = $1
);

-- name: CreateList :one
INSERT INTO lists (id, user_id, title, description)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, title, description, created_at, updated_at;

-- name: UpdateList :one
UPDATE lists
SET title = $1, description = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $3
RETURNING *;

-- name: DeleteList :execrows
DELETE FROM lists
WHERE id = $1;