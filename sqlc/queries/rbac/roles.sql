-- name: CreateRole :one
INSERT INTO roles (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetRole :one
SELECT * FROM roles WHERE id = $1;

-- name: ListRoles :many
SELECT * FROM roles;

-- name: UpdateRole :exec
UPDATE roles SET name = $2, description = $3 WHERE id = $1;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1;

-- name: GetRoleByName :one
SELECT * FROM roles WHERE name = $1;

-- name: GetRolePermissions :many
SELECT p.*
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = $1;

-- name: RemovePermissionFromRole :exec
DELETE FROM role_permissions
WHERE role_id = $1 AND permission_id = $2;

-- name: RemoveUserFromRole :exec
DELETE FROM user_roles
WHERE role_id = $1 AND user_id = $2;

-- name: GetRoleUsers :many
SELECT u.*
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
WHERE ur.role_id = $1;