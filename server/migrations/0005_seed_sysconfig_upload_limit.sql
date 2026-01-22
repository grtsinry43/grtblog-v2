-- +goose Up
INSERT INTO sys_config (config_key, value) VALUES
    ('upload.maxSizeMB', '50')
ON CONFLICT (config_key) DO NOTHING;

-- +goose Down
DELETE FROM sys_config WHERE config_key IN (
    'upload.maxSizeMB'
);
