-- name: CreateURL :one
INSERT INTO urls (url_code, original_url)
VALUES ($1, $2)
RETURNING id, url_code, original_url, click_count, created_at, updated_at;

-- name: GetURLByID :one
SELECT id, url_code, original_url, click_count, created_at, updated_at
FROM urls
WHERE id = $1;

-- name: GetURLByURLCode :one
SELECT id, url_code, original_url, click_count, created_at, updated_at
FROM urls
WHERE url_code = $1;

-- name: UpdateURL :one
UPDATE urls
SET url_code = $2, original_url = $3, click_count = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, url_code, original_url, click_count, created_at, updated_at;

-- name: IncrementURLCount :one
UPDATE urls
SET click_count = click_count + 1, updated_at = NOW()
WHERE id = $1
RETURNING id, url_code, original_url, click_count, created_at, updated_at;

-- name: DeleteURL :exec
DELETE FROM urls
WHERE id = $1;

-- name: GetAllURLs :many
SELECT id, url_code, original_url, click_count, created_at, updated_at
FROM urls
ORDER BY created_at DESC;

