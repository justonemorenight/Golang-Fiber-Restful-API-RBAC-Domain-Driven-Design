// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: base.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id,
    token,
    expires_at
) VALUES ($1, $2, $3)
RETURNING id, user_id, token, expires_at, created_at
`

type CreateRefreshTokenParams struct {
	UserID    int32            `json:"user_id"`
	Token     string           `json:"token"`
	ExpiresAt pgtype.Timestamp `json:"expires_at"`
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRow(ctx, createRefreshToken, arg.UserID, arg.Token, arg.ExpiresAt)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteRefreshToken = `-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token = $1
`

func (q *Queries) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := q.db.Exec(ctx, deleteRefreshToken, token)
	return err
}

const deleteUserRefreshTokens = `-- name: DeleteUserRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE user_id = $1
`

func (q *Queries) DeleteUserRefreshTokens(ctx context.Context, userID int32) error {
	_, err := q.db.Exec(ctx, deleteUserRefreshTokens, userID)
	return err
}

const getRefreshToken = `-- name: GetRefreshToken :one
SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens
WHERE token = $1 AND expires_at > NOW()
LIMIT 1
`

func (q *Queries) GetRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	row := q.db.QueryRow(ctx, getRefreshToken, token)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}