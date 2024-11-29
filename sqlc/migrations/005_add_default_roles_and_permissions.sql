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
('users.read_self', 'View personal information', 'users', 'read_self'),
('roles.create', 'Create new role', 'roles', 'create'),
('roles.read', 'View role information', 'roles', 'read'),
('roles.update', 'Update role information', 'roles', 'update'),
('roles.delete', 'Delete role', 'roles', 'delete'),
('roles.list', 'View all roles', 'roles', 'list'),
('roles.assign_permission', 'Assign permissions to role', 'roles', 'assign_permission'),
('roles.remove_permission', 'Remove permissions from role', 'roles', 'remove_permission'),
('roles.view_permissions', 'View role permissions', 'roles', 'view_permissions'),
('roles.assign_user', 'Assign role to user', 'roles', 'assign_user'),
('roles.remove_user', 'Remove role from user', 'roles', 'remove_user'),
('roles.view_users', 'View users with specific role', 'roles', 'view_users');

-- Assign permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id 
FROM roles r, permissions p
WHERE r.name = 'admin' 
AND p.resource = 'roles';

-- Assign read_self permission to member role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'member' AND p.name = 'users.read_self';

-- +goose Down
DELETE FROM role_permissions;
DELETE FROM permissions;
DELETE FROM roles; 