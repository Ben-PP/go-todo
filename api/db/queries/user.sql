-- name: GetUserById :one
SELECT id, username, is_admin, created_at
FROM users
WHERE id = $1;

-- name: GetAllUsers :many
SELECT id, username, is_admin, created_at
FROM users;

-- name: CreateUser :one
INSERT INTO users (id, username, password_hash, is_admin)
VALUES ($1, $2, $3, $4)
RETURNING id, username, is_admin, created_at;

-- name: UpdateUser :one
UPDATE users
SET username = $2, password_hash = $3, is_admin = $4
WHERE id = $1
RETURNING id, username, is_admin, created_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetPasswordHashByUsername :one
SELECT password_hash
FROM users
WHERE username = $1;