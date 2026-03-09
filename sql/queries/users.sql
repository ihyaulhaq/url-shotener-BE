-- name: CreateUser :one
INSERT INTO users (
  username,
  email, 
  hashed_password, 
  is_active
)
VALUES (
  $1, 
  $2, 
  $3, 
  $4
)
RETURNING id, username, email, hashed_password, is_active, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, username, email, hashed_password, is_active, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, username, email, hashed_password, is_active, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT id, username, email, hashed_password, is_active, created_at, updated_at
FROM users
WHERE username = $1;

-- name: UpdateUser :one
UPDATE users
SET username = $2, email = $3, hashed_password = $4, is_active = $5, updated_at = NOW()
WHERE id = $1
RETURNING id, username, email, hashed_password, is_active, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetAllActiveUsers :many
SELECT id, username, email, hashed_password, is_active, created_at, updated_at
FROM users
WHERE is_active = true
ORDER BY created_at DESC;
