-- +goose Up
-- Tạo roles mặc định
INSERT INTO roles (name, description) VALUES
('admin', 'Quản trị viên với toàn quyền'),
('member', 'Thành viên thông thường');

-- Tạo permissions
INSERT INTO permissions (name, description, resource, action) VALUES
('users.create', 'Tạo user mới', 'users', 'create'),
('users.read', 'Xem thông tin user', 'users', 'read'),
('users.update', 'Cập nhật user', 'users', 'update'),
('users.delete', 'Xóa user', 'users', 'delete'),
('users.list', 'Xem danh sách users', 'users', 'list'),
('users.read_self', 'Xem thông tin cá nhân', 'users', 'read_self');

-- Gán permissions cho admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id 
FROM roles r, permissions p
WHERE r.name = 'admin';

-- Gán permission read_self cho member role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'member' AND p.name = 'users.read_self';

-- +goose Down
DELETE FROM role_permissions;
DELETE FROM permissions;
DELETE FROM roles; 