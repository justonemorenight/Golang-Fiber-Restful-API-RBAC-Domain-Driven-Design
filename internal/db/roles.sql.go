// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: roles.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createRole = `-- name: CreateRole :one
INSERT INTO roles (name, description)
VALUES ($1, $2)
RETURNING id, name, description
`

type CreateRoleParams struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
}

func (q *Queries) CreateRole(ctx context.Context, arg CreateRoleParams) (Role, error) {
	row := q.db.QueryRow(ctx, createRole, arg.Name, arg.Description)
	var i Role
	err := row.Scan(&i.ID, &i.Name, &i.Description)
	return i, err
}

const deleteRole = `-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1
`

func (q *Queries) DeleteRole(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteRole, id)
	return err
}

const getRole = `-- name: GetRole :one
SELECT id, name, description FROM roles WHERE id = $1
`

func (q *Queries) GetRole(ctx context.Context, id int32) (Role, error) {
	row := q.db.QueryRow(ctx, getRole, id)
	var i Role
	err := row.Scan(&i.ID, &i.Name, &i.Description)
	return i, err
}

const getRoleByName = `-- name: GetRoleByName :one
SELECT id, name, description FROM roles WHERE name = $1
`

func (q *Queries) GetRoleByName(ctx context.Context, name string) (Role, error) {
	row := q.db.QueryRow(ctx, getRoleByName, name)
	var i Role
	err := row.Scan(&i.ID, &i.Name, &i.Description)
	return i, err
}

const getRolePermissions = `-- name: GetRolePermissions :many
SELECT p.id, p.name, p.description, p.resource, p.action
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = $1
`

func (q *Queries) GetRolePermissions(ctx context.Context, roleID int32) ([]Permission, error) {
	rows, err := q.db.Query(ctx, getRolePermissions, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Permission{}
	for rows.Next() {
		var i Permission
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Resource,
			&i.Action,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoleUsers = `-- name: GetRoleUsers :many
SELECT u.id, u.name, u.email, u.password, u.created_at, u.updated_at
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
WHERE ur.role_id = $1
`

func (q *Queries) GetRoleUsers(ctx context.Context, roleID int32) ([]User, error) {
	rows, err := q.db.Query(ctx, getRoleUsers, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Password,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRoles = `-- name: ListRoles :many
SELECT id, name, description FROM roles
`

func (q *Queries) ListRoles(ctx context.Context) ([]Role, error) {
	rows, err := q.db.Query(ctx, listRoles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Role{}
	for rows.Next() {
		var i Role
		if err := rows.Scan(&i.ID, &i.Name, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removePermissionFromRole = `-- name: RemovePermissionFromRole :exec
DELETE FROM role_permissions
WHERE role_id = $1 AND permission_id = $2
`

type RemovePermissionFromRoleParams struct {
	RoleID       int32 `json:"role_id"`
	PermissionID int32 `json:"permission_id"`
}

func (q *Queries) RemovePermissionFromRole(ctx context.Context, arg RemovePermissionFromRoleParams) error {
	_, err := q.db.Exec(ctx, removePermissionFromRole, arg.RoleID, arg.PermissionID)
	return err
}

const removeUserFromRole = `-- name: RemoveUserFromRole :exec
DELETE FROM user_roles
WHERE role_id = $1 AND user_id = $2
`

type RemoveUserFromRoleParams struct {
	RoleID int32 `json:"role_id"`
	UserID int32 `json:"user_id"`
}

func (q *Queries) RemoveUserFromRole(ctx context.Context, arg RemoveUserFromRoleParams) error {
	_, err := q.db.Exec(ctx, removeUserFromRole, arg.RoleID, arg.UserID)
	return err
}

const updateRole = `-- name: UpdateRole :exec
UPDATE roles SET name = $2, description = $3 WHERE id = $1
`

type UpdateRoleParams struct {
	ID          int32       `json:"id"`
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
}

func (q *Queries) UpdateRole(ctx context.Context, arg UpdateRoleParams) error {
	_, err := q.db.Exec(ctx, updateRole, arg.ID, arg.Name, arg.Description)
	return err
}
