-- name: CreateList :one
INSERT INTO lists (id, user_id, title, description)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, title, description, created_at, updated_at;
