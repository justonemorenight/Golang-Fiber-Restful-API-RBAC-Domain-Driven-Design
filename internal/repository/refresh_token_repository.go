package repository

import (
	"backend-fiber/internal/db"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type RefreshTokenRepository struct {
	queries *db.Queries
}

func NewRefreshTokenRepository(queries *db.Queries) *RefreshTokenRepository {
	return &RefreshTokenRepository{queries: queries}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, userID int32, token string, expiresAt time.Time) error {
	params := db.CreateRefreshTokenParams{
		UserID: userID,
		Token:  token,
		ExpiresAt: pgtype.Timestamp{
			Time:  expiresAt,
			Valid: true,
		},
	}
	_, err := r.queries.CreateRefreshToken(ctx, params)
	return err
}

func (r *RefreshTokenRepository) Get(ctx context.Context, token string) (*db.RefreshToken, error) {
	refreshToken, err := r.queries.GetRefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, token string) error {
	return r.queries.DeleteRefreshToken(ctx, token)
}

func (r *RefreshTokenRepository) DeleteAllForUser(ctx context.Context, userID int32) error {
	return r.queries.DeleteUserRefreshTokens(ctx, userID)
}
