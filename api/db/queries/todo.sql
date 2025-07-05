-- name: CreateTodo :one
INSERT INTO todos (id, list_id, user_id, parent_id, title, description, complete_before)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetTodoByIdWithListId :one
SELECT * FROM todos
WHERE id = $1 AND list_id = $2;

-- name: UpdateTodo :one
UPDATE todos
SET title = $1, description = $2, completed = $3, complete_before = $4, updated_at = CURRENT_TIMESTAMP, completed_at = CASE WHEN $3 THEN CURRENT_TIMESTAMP ELSE NULL END
WHERE id = $5
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;

-- name: DeleteTodoByIdWithListId :exec
DELETE FROM todos
WHERE id = $1 AND list_id = $2;