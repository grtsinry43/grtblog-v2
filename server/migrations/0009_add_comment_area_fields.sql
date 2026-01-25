-- +goose Up
ALTER TABLE comment_area
    ADD COLUMN IF NOT EXISTS area_type VARCHAR(20) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS content_id BIGINT,
    ADD COLUMN IF NOT EXISTS is_closed BOOLEAN NOT NULL DEFAULT FALSE,
    ALTER COLUMN area_name TYPE VARCHAR(255);

-- backfill area_type/content_id from area_name = "<type>:<content_id>"
UPDATE comment_area
SET area_type = split_part(area_name, ':', 1),
    content_id = NULLIF(split_part(area_name, ':', 2), '')::BIGINT
WHERE area_name LIKE '%:%';

CREATE UNIQUE INDEX IF NOT EXISTS uq_comment_area_type_content ON comment_area (area_type, content_id);
CREATE INDEX IF NOT EXISTS idx_comment_area_is_closed ON comment_area (is_closed);

ALTER TABLE comment
    ALTER COLUMN location TYPE VARCHAR(255);

-- +goose Down
DROP INDEX IF EXISTS idx_comment_area_is_closed;
DROP INDEX IF EXISTS uq_comment_area_type_content;

ALTER TABLE comment_area
    DROP COLUMN IF EXISTS is_closed,
    DROP COLUMN IF EXISTS content_id,
    DROP COLUMN IF EXISTS area_type,
    ALTER COLUMN area_name TYPE VARCHAR(45);

ALTER TABLE comment
    ALTER COLUMN location TYPE VARCHAR(45);
