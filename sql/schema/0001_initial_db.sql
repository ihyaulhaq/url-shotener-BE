-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  username TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  hashed_password TEXT NOT NULL,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

CREATE TABLE refresh_tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  token TEXT UNIQUE NOT NULL,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  expires_at TIMESTAMP NOT NULL,
  revoked_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT now (),
  updated_at TIMESTAMP NOT NULL DEFAULT now ()
);

CREATE TABLE urls (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  url_code TEXT NOT NULL,
  original_url TEXT NOT NULL,
  click_count INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

CREATE TABLE url_users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  url_id UUID NOT NULL REFERENCES urls (id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT now (),
  updated_at TIMESTAMP DEFAULT now (),
  UNIQUE (url_id, user_id)
);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_username ON users (username);

CREATE INDEX idx_users_is_active ON users (is_active);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);

CREATE INDEX idx_refresh_tokens_token ON refresh_tokens (token);

CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens (expires_at);

CREATE INDEX idx_urls_short_url ON urls (url_code);

CREATE INDEX idx_urls_created_at ON urls (created_at);

CREATE INDEX idx_url_users_user_id ON url_users (user_id);

CREATE INDEX idx_url_users_url_id ON url_users (url_id);

CREATE INDEX idx_url_users_user_url ON url_users (user_id, url_id);

-- +goose Down
DROP TABLE IF EXISTS url_users;

DROP TABLE IF EXISTS urls;

DROP TABLE IF EXISTS refresh_tokens;

DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "pgcrypto";
