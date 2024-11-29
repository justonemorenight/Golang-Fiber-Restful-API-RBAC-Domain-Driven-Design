package repository

import (
	"context"
	"time"
)

type RefreshTokenRepository struct {
	// TODO: add db connection
}

func (r *RefreshTokenRepository) Create(ctx context.Context, userID int32, token string, expiresAt time.Time) error {
	// TODO: implement
	return nil
}

func (r *RefreshTokenRepository) Get(ctx context.Context, token string) (*RefreshToken, error) {
	// TODO: implement
	return nil, nil
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, token string) error {
	// TODO: implement
	return nil
}

func (r *RefreshTokenRepository) GetTokenUsage(ctx context.Context, tokenID int64) (*TokenUsage, error) {
	// TODO: implement
	return nil, nil
}

func (r *RefreshTokenRepository) UpdateTokenUsage(ctx context.Context, tokenID int64, ip, userAgent string) error {
	// TODO: implement
	return nil
}

type RefreshToken struct {
	ID     int64
	UserID int32
	Token  string
}

type TokenUsage struct {
	UsageCount int
	LastUsedAt time.Time
	LastUsedIP string
}
