-- +goose Up
ALTER TABLE nav_menu
    ADD COLUMN IF NOT EXISTS icon VARCHAR(64);

-- +goose Down
ALTER TABLE nav_menu
    DROP COLUMN IF EXISTS icon;
