-- +goose Up
ALTER TABLE refresh_tokens
ADD COLUMN last_used_at TIMESTAMPTZ,
ADD COLUMN usage_count INTEGER DEFAULT 0,
ADD COLUMN last_used_ip TEXT,
ADD COLUMN last_user_agent TEXT;

-- +goose Down
ALTER TABLE refresh_tokens
DROP COLUMN last_used_at,
DROP COLUMN usage_count,
DROP COLUMN last_used_ip,
DROP COLUMN last_user_agent;