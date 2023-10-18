-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, email, apikey)
VALUES ($1, $2, $3, $4, $5, encode(sha256(random()::text::bytea), 'hex'))
RETURNING *;

-- name: GetUserByApiKey :one
SELECT * FROM users
WHERE apikey = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;
