-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at)
VALUES ($1, $2, $3)
RETURNING id, token, user_id, expires_at, revoked_at, created_at, updated_at;

-- name: GetRefreshTokenByToken :one
SELECT id, token, user_id, expires_at, revoked_at, created_at, updated_at
FROM refresh_tokens
WHERE token = $1;

-- name: GetRefreshTokenByID :one
SELECT id, token, user_id, expires_at, revoked_at, created_at, updated_at
FROM refresh_tokens
WHERE id = $1;

-- name: GetRefreshTokensByUserID :many
SELECT id, token, user_id, expires_at, revoked_at, created_at, updated_at
FROM refresh_tokens
WHERE user_id = $1 AND revoked_at IS NULL
ORDER BY created_at DESC;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING id, token, user_id, expires_at, revoked_at, created_at, updated_at;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW();

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE id = $1;

-- name: ValidateRefreshToken :one
SELECT id, token, user_id, expires_at, revoked_at, created_at, updated_at
FROM refresh_tokens
WHERE token = $1 AND expires_at > NOW() AND revoked_at IS NULL;

