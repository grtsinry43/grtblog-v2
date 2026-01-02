-- +goose Up
INSERT INTO permission (permission_name) VALUES ('log:read') ON CONFLICT (permission_name) DO NOTHING;
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id FROM role r, permission p
WHERE r.role_name = 'admin' AND p.permission_name = 'log:read'
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM role_permission WHERE permission_id IN (SELECT id FROM permission WHERE permission_name = 'log:read');
DELETE FROM permission WHERE permission_name = 'log:read';
