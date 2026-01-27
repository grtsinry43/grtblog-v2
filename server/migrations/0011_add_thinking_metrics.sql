-- +goose Up
ALTER TABLE thinking ADD COLUMN updated_at TIMESTAMPTZ DEFAULT now();

ALTER TABLE thinking ADD COLUMN comment_id BIGINT NOT NULL;

CREATE TABLE IF NOT EXISTS thinking_metrics
(
    thinking_id   BIGINT PRIMARY KEY
        REFERENCES thinking (id) ON DELETE CASCADE,
    views      BIGINT  NOT NULL DEFAULT 0,
    likes      INTEGER NOT NULL DEFAULT 0,
    comments   INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ      DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS thinking_metrics;
ALTER TABLE thinking DROP COLUMN IF EXISTS comment_id;
ALTER TABLE thinking DROP COLUMN IF EXISTS updated_at;
