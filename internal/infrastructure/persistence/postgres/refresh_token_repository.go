package repository

import (
	"backend-fiber/internal/db"
	"backend-fiber/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	queries *db.Queries
	db      *pgxpool.Pool
}

func NewRefreshTokenRepository(queries *db.Queries, db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		queries: queries,
		db:      db,
	}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, userID int32, token string, expiresAt time.Time) error {
	params := db.CreateRefreshTokenParams{
		UserID:  userID,
		Token:   token,
		Column3: pgtype.Timestamptz{Time: expiresAt, Valid: true},
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

func (r *RefreshTokenRepository) UpdateTokenUsage(ctx context.Context, tokenID int32, ip string, userAgent string) error {
	query := `
		UPDATE refresh_tokens 
		SET last_used_at = NOW(),
			usage_count = usage_count + 1,
			last_used_ip = $2,
			last_user_agent = $3
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, tokenID, ip, userAgent)
	return err
}

func (r *RefreshTokenRepository) GetTokenUsage(ctx context.Context, tokenID int32) (*models.TokenUsage, error) {
	query := `
		SELECT id, user_id, last_used_at, usage_count, last_used_ip, last_user_agent
		FROM refresh_tokens
		WHERE id = $1
	`
	var usage models.TokenUsage
	err := r.db.QueryRow(ctx, query, tokenID).Scan(
		&usage.TokenID,
		&usage.UserID,
		&usage.LastUsedAt,
		&usage.UsageCount,
		&usage.LastUsedIP,
		&usage.LastUserAgent,
	)
	if err != nil {
		return nil, err
	}
	return &usage, nil
}
