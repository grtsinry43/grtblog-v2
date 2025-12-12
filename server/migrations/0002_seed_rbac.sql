-- +goose Up
-- 基础 RBAC 数据（角色/权限/绑定）

-- 角色
INSERT INTO role (role_name) VALUES ('user') ON CONFLICT (role_name) DO NOTHING;
INSERT INTO role (role_name) VALUES ('admin') ON CONFLICT (role_name) DO NOTHING;

-- 权限
INSERT INTO permission (permission_name) VALUES ('config:read') ON CONFLICT (permission_name) DO NOTHING;
INSERT INTO permission (permission_name) VALUES ('config:write') ON CONFLICT (permission_name) DO NOTHING;

-- 角色-权限：管理员拥有配置读写
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id FROM role r, permission p
WHERE r.role_name = 'admin' AND p.permission_name IN ('config:read', 'config:write')
ON CONFLICT DO NOTHING;

-- 可选：普通用户拥有配置读取（如不需要，可删除以下语句）
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id FROM role r, permission p
WHERE r.role_name = 'user' AND p.permission_name = 'config:read'
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM role_permission WHERE role_id IN (SELECT id FROM role WHERE role_name IN ('admin', 'user'));
DELETE FROM permission WHERE permission_name IN ('config:read', 'config:write');
DELETE FROM role WHERE role_name IN ('admin', 'user');
