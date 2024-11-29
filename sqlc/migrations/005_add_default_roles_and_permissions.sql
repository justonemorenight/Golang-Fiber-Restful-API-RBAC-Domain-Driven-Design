-- +goose Up
-- Create default roles
INSERT INTO roles (name, description) VALUES
('admin', 'Administrator with full access'),
('member', 'Regular member');

-- Create permissions
INSERT INTO permissions (name, description, resource, action) VALUES
('users.create', 'Create new user', 'users', 'create'),
('users.read', 'View user information', 'users', 'read'),
('users.update', 'Update user', 'users', 'update'),
('users.delete', 'Delete user', 'users', 'delete'),
('users.list', 'View user list', 'users', 'list'),
('users.read_self', 'View personal information', 'users', 'read_self');

-- Assign permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id 
FROM roles r, permissions p
WHERE r.name = 'admin';

-- Assign read_self permission to member role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'member' AND p.name = 'users.read_self';

-- +goose Down
DELETE FROM role_permissions;
DELETE FROM permissions;
DELETE FROM roles; 