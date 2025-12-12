-- +goose Up
-- 基础 Turnstile 配置占位，便于运行时动态开关

INSERT INTO sys_config (config_key, value) VALUES
    ('turnstile.enabled', 'false'),
    ('turnstile.secret', ''),
    ('turnstile.siteKey', ''),
    ('turnstile.verifyURL', ''),
    ('turnstile.timeoutSeconds', '')
ON CONFLICT (config_key) DO NOTHING;

-- +goose Down
DELETE FROM sys_config WHERE config_key IN (
    'turnstile.enabled',
    'turnstile.secret',
    'turnstile.siteKey',
    'turnstile.verifyURL',
    'turnstile.timeoutSeconds'
);
