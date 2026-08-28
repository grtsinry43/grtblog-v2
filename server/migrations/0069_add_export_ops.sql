-- +goose Up
CREATE SCHEMA IF NOT EXISTS export_ops;

CREATE TABLE IF NOT EXISTS export_ops.export_record
(
    id               VARCHAR(36) PRIMARY KEY,
    filename         VARCHAR(255) NOT NULL,
    status           VARCHAR(24) NOT NULL,
    stage            VARCHAR(48) NOT NULL DEFAULT 'queued',
    trigger_type     VARCHAR(24) NOT NULL DEFAULT 'manual',
    mode             VARCHAR(16) NOT NULL DEFAULT 'both',
    size_bytes       BIGINT NOT NULL DEFAULT 0,
    sha256           VARCHAR(64),
    app_version      VARCHAR(64),
    site_name        TEXT,
    site_url         TEXT,
    article_count    BIGINT NOT NULL DEFAULT 0,
    moments_count    BIGINT NOT NULL DEFAULT 0,
    pages_count      BIGINT NOT NULL DEFAULT 0,
    thinkings_count  BIGINT NOT NULL DEFAULT 0,
    image_count      BIGINT NOT NULL DEFAULT 0,
    failed_image_count BIGINT NOT NULL DEFAULT 0,
    error_message    TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at       TIMESTAMPTZ,
    completed_at     TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_export_record_created_at
    ON export_ops.export_record (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_export_record_status
    ON export_ops.export_record (status);

CREATE TABLE IF NOT EXISTS export_ops.download_ticket
(
    token_hash VARCHAR(64) PRIMARY KEY,
    export_id  VARCHAR(36) NOT NULL REFERENCES export_ops.export_record (id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_export_download_ticket_expires_at
    ON export_ops.download_ticket (expires_at);

-- +goose Down
DROP SCHEMA IF EXISTS export_ops CASCADE;
