-- +goose Up
ALTER TABLE thinking ADD COLUMN author_id BIGINT;

UPDATE thinking t
SET author_id = u.id
FROM app_user u
WHERE t.author_id IS NULL
  AND (u.username = t.author OR u.nickname = t.author);

UPDATE thinking
SET author_id = (
    SELECT id
    FROM app_user
    WHERE is_admin = TRUE
    ORDER BY id
    LIMIT 1
)
WHERE author_id IS NULL;

ALTER TABLE thinking
    ADD CONSTRAINT fk_thinking_author FOREIGN KEY (author_id) REFERENCES app_user (id) ON DELETE SET NULL;

ALTER TABLE thinking DROP COLUMN author;

-- +goose Down
ALTER TABLE thinking ADD COLUMN author VARCHAR(45) NOT NULL DEFAULT '原创';

UPDATE thinking t
SET author = COALESCE(u.nickname, u.username, '原创')
FROM app_user u
WHERE t.author_id = u.id;

ALTER TABLE thinking DROP CONSTRAINT IF EXISTS fk_thinking_author;
ALTER TABLE thinking DROP COLUMN author_id;
