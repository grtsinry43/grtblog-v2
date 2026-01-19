-- +goose Up
ALTER TABLE article
    ADD COLUMN IF NOT EXISTS content_hash VARCHAR(32) NOT NULL DEFAULT '';
ALTER TABLE moment
    ADD COLUMN IF NOT EXISTS content_hash VARCHAR(32) NOT NULL DEFAULT '';
ALTER TABLE page
    ADD COLUMN IF NOT EXISTS content_hash VARCHAR(32) NOT NULL DEFAULT '';

UPDATE article
SET content_hash = md5(coalesce(title, '') || coalesce(lead_in, '') || coalesce(content, ''));
UPDATE moment
SET content_hash = md5(coalesce(title, '') || coalesce(summary, '') || coalesce(content, ''));
UPDATE page
SET content_hash = md5(coalesce(title, '') || coalesce(description, '') || coalesce(content, ''));

-- +goose Down
ALTER TABLE page
    DROP COLUMN IF EXISTS content_hash;
ALTER TABLE moment
    DROP COLUMN IF EXISTS content_hash;
ALTER TABLE article
    DROP COLUMN IF EXISTS content_hash;
