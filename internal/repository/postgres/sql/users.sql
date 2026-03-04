-- name: CreateUser :one
INSERT INTO users (name, email)
VALUES ($1, $2)
RETURNING id, name, email, created_at, is_active, role;

-- name: GetUserByID :one
SELECT id, name, email, created_at, is_active, role
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, name, email, created_at, is_active, role
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, name, email, created_at, is_active, role
FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: SetUserActive :exec
UPDATE users
SET is_active = $2
WHERE id = $1;