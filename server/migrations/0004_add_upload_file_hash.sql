-- +goose Up
ALTER TABLE upload_file
    ADD COLUMN IF NOT EXISTS hash VARCHAR(64);
CREATE UNIQUE INDEX IF NOT EXISTS uq_upload_file_hash ON upload_file (hash);

-- +goose Down
DROP INDEX IF EXISTS uq_upload_file_hash;
ALTER TABLE upload_file
    DROP COLUMN IF EXISTS hash;
