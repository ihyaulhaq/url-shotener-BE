-- name: CreateURL :one
INSERT INTO urls (short_url, original_url, count_code)
VALUES ($1, $2, $3)
RETURNING id, short_url, original_url, count_code, created_at, updated_at;

-- name: GetURLByID :one
SELECT id, short_url, original_url, count_code, created_at, updated_at
FROM urls
WHERE id = $1;

-- name: GetURLByShortURL :one
SELECT id, short_url, original_url, count_code, created_at, updated_at
FROM urls
WHERE short_url = $1;

-- name: UpdateURL :one
UPDATE urls
SET short_url = $2, original_url = $3, count_code = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, short_url, original_url, count_code, created_at, updated_at;

-- name: IncrementURLCount :one
UPDATE urls
SET count_code = count_code + 1, updated_at = NOW()
WHERE id = $1
RETURNING id, short_url, original_url, count_code, created_at, updated_at;

-- name: DeleteURL :exec
DELETE FROM urls
WHERE id = $1;

-- name: GetAllURLs :many
SELECT id, short_url, original_url, count_code, created_at, updated_at
FROM urls
ORDER BY created_at DESC;

