-- name: CreateURLUser :one
INSERT INTO url_users (url_id, user_id)
VALUES ($1, $2)
RETURNING id, url_id, user_id, created_at, updated_at;

-- name: GetURLsByUserID :many
SELECT u.id, u.url_code, u.original_url, u.click_count, u.created_at, u.updated_at
FROM urls u
JOIN url_users uu ON uu.url_id = u.id
WHERE uu.user_id = $1
ORDER BY u.created_at DESC;
