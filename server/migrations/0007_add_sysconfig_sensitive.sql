-- +goose Up
ALTER TABLE sys_config
    ADD COLUMN IF NOT EXISTS is_sensitive BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS group_path TEXT,
    ADD COLUMN IF NOT EXISTS label VARCHAR(100),
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS value_type VARCHAR(20) NOT NULL DEFAULT 'string',
    ADD COLUMN IF NOT EXISTS enum_options JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS default_value TEXT,
    ADD COLUMN IF NOT EXISTS visible_when JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS sort INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS meta JSONB NOT NULL DEFAULT '{}'::jsonb;

UPDATE sys_config
SET is_sensitive = TRUE
WHERE config_key IN ('turnstile.secret');

UPDATE sys_config
SET group_path = 'security/turnstile',
    label = '启用',
    value_type = 'bool',
    sort = 10,
    meta = '{"inputType":"switch"}'::jsonb
WHERE config_key = 'turnstile.enabled';

UPDATE sys_config
SET group_path = 'security/turnstile',
    label = '密钥',
    value_type = 'string',
    visible_when = '[{"key":"turnstile.enabled","op":"eq","value":true}]'::jsonb,
    sort = 20,
    meta = '{"inputType":"password"}'::jsonb
WHERE config_key = 'turnstile.secret';

UPDATE sys_config
SET group_path = 'security/turnstile',
    label = '站点 Key',
    value_type = 'string',
    visible_when = '[{"key":"turnstile.enabled","op":"eq","value":true}]'::jsonb,
    sort = 30
WHERE config_key = 'turnstile.siteKey';

UPDATE sys_config
SET group_path = 'security/turnstile',
    label = '校验地址',
    value_type = 'string',
    visible_when = '[{"key":"turnstile.enabled","op":"eq","value":true}]'::jsonb,
    sort = 40
WHERE config_key = 'turnstile.verifyURL';

UPDATE sys_config
SET group_path = 'security/turnstile',
    label = '超时(秒)',
    value_type = 'number',
    visible_when = '[{"key":"turnstile.enabled","op":"eq","value":true}]'::jsonb,
    sort = 50,
    meta = '{"unit":"s"}'::jsonb
WHERE config_key = 'turnstile.timeoutSeconds';

UPDATE sys_config
SET group_path = 'storage/upload',
    label = '最大上传(MB)',
    value_type = 'number',
    sort = 10,
    meta = '{"min":1,"max":50}'::jsonb
WHERE config_key = 'upload.maxSizeMB';

UPDATE sys_config
SET group_path = 'webhook',
    label = '超时(秒)',
    value_type = 'number',
    sort = 10,
    meta = '{"unit":"s"}'::jsonb
WHERE config_key = 'webhook.timeoutSeconds';

UPDATE sys_config
SET group_path = 'webhook',
    label = '并发数',
    value_type = 'number',
    sort = 20
WHERE config_key = 'webhook.workers';

UPDATE sys_config
SET group_path = 'webhook',
    label = '队列长度',
    value_type = 'number',
    sort = 30
WHERE config_key = 'webhook.queueSize';

-- +goose Down
ALTER TABLE sys_config
    DROP COLUMN IF EXISTS meta,
    DROP COLUMN IF EXISTS sort,
    DROP COLUMN IF EXISTS visible_when,
    DROP COLUMN IF EXISTS default_value,
    DROP COLUMN IF EXISTS enum_options,
    DROP COLUMN IF EXISTS value_type,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS label,
    DROP COLUMN IF EXISTS group_path,
    DROP COLUMN IF EXISTS is_sensitive;
