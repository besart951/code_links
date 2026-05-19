package auth

import (
	"context"
	"time"
)

type Repository interface {
	FindUserByEmail(ctx context.Context, email string) (User, error)
	FindUserByID(ctx context.Context, userID string) (User, error)
	TouchLastLogin(ctx context.Context, userID string, at time.Time) error
	StoreRefreshToken(ctx context.Context, token RefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (RefreshToken, error)
	RevokeRefreshTokenByHash(ctx context.Context, tokenHash string, at time.Time) error
}
