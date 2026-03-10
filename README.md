# URL Shotener BE

Go + Postgres backend for a simple URL shortener.

## Features (current)

- Create a short URL code for an `original_url`
- Redirect from `/{shortUrl}` to the original URL
- Basic health check endpoint
- Postgres persistence via `sqlc`-generated queries
- Click count increment on redirect

## Tech stack

- **Go** (see `go.mod` for version)
- **net/http** ServeMux routing (Go 1.22+ style patterns)
- **PostgreSQL**
- **sqlc** for type-safe query generation (output in `internal/database`)
- **gotenv** for loading `.env`

## Project structure

- `cmd/api-server/`: app entrypoint (`main.go`)
- `internal/config/`: environment-based configuration
- `internal/handler/`: HTTP handlers and route wiring
- `internal/service/`: business logic (short-code generation, validation)
- `internal/store/`: DB store wrapper around sqlc `Queries`
- `internal/database/`: sqlc generated models + queries
- `sql/schema/`: SQL migrations/schema (currently written for `goose`)
- `sql/queries/`: sqlc query definitions
- `pkg/utils/`: shared helpers (JSON responses)

## Requirements

- **PostgreSQL** running and reachable
- **Go** installed

## Configuration (environment variables)

The server loads `.env` if present (see `internal/config/config.go`). Defaults are provided for most values.

Create a `.env` file at the repo root:

```env
# App
APP_ENV=development
BASE_URL=http://localhost:8080

# HTTP server
PORT=8080
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
SERVER_IDLE_TIMEOUT=60s

# Database (Postgres)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=url_shortener
DB_SSLMODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
```

## Database setup

The schema lives in `sql/schema/0001_initial_db.sql` and creates:

- `urls`: `url_code`, `original_url`, `click_count`, timestamps
- `users`, `refresh_tokens`, `url_users`: present in schema/queries but **not wired to HTTP routes yet**

### Apply schema

How you apply migrations depends on your workflow. The migration file contains `goose` directives (`-- +goose Up/Down`).

- **Option A (recommended for now)**: run the SQL manually in your DB client against your target database
- **Option B**: use `goose` if you have it installed/configured

### sqlc generation

`sqlc` config is in `sqlc.yaml` and generates Go code into `internal/database`.

If you edit SQL under `sql/queries` or schema under `sql/schema`, re-generate:

```bash
sqlc generate
```

## Run locally

From repo root:

```bash
go run ./cmd/api-server
```

The server starts on `:${PORT}` (default `:8080`).

## HTTP API

### Health

`GET /health`

Response:

```json
{"data":{"status":"ok"}}
```

### Shorten URL

`POST /api/urls/shorten`

Body:

```json
{"original_url":"https://example.com/some/path"}
```

Response (`201`):

```json
{
  "data": {
    "short_url": "http://localhost:8080/2a1b3c4d5e6f",
    "original_url": "https://example.com/some/path"
  }
}
```

Curl:

```bash
curl -s -X POST "http://localhost:8080/api/urls/shorten" ^
  -H "Content-Type: application/json" ^
  -d "{\"original_url\":\"https://example.com\"}"
```

### Redirect

`GET /{shortUrl}`

Example:

- `GET /2a1b3c4d5e6f` → `302 Found` redirect to the stored `original_url`
- On success, the server increments `click_count` for that URL

## Notes / behavior

- **Short code generation**: SHA-256 of the URL plus a retry salt, first 6 bytes hex-encoded (12 chars). Collisions are checked against the DB and retried (up to 5 attempts).
- **Validation**: request body uses `validator` tag `required,url`, and service additionally enforces `http`/`https` scheme.
- **Errors**: JSON envelope format `{ "data": ..., "error": "..." }` (see `pkg/utils/response.go`).

## Troubleshooting

- **DB connection fails**: verify `.env` values and that Postgres is reachable; the server pings the DB on startup.
- **Short URL host is wrong**: set `BASE_URL` (used to build `short_url` in responses).
