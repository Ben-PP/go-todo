-- name: CreateTodo :one
INSERT INTO todos (id, list_id, user_id, parent_id, title, description, complete_before)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;