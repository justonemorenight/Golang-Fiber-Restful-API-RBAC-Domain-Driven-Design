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