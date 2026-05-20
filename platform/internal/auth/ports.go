package auth

import (
	"context"
	"time"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, userID UserID) (User, error)
	TouchLastLogin(ctx context.Context, userID UserID, at time.Time) error
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session Session, refresh RefreshToken) error
	FindBySessionTokenHash(ctx context.Context, tokenHash string) (AuthenticatedSession, error)
	FindByRefreshTokenHash(ctx context.Context, tokenHash string) (Session, RefreshToken, error)
	RotateRefreshToken(ctx context.Context, oldTokenHash string, next RefreshToken, revokedAt time.Time) error
	RevokeSession(ctx context.Context, sessionID SessionID, revokedAt time.Time) error
}

type PasswordVerifier interface {
	VerifyPassword(password, passwordHash string) bool
}

type TokenHasher interface {
	HashToken(raw string) string
}

type Clock interface {
	Now() time.Time
}

type IDGenerator interface {
	NewID(prefix string) (string, error)
	NewSecret(bytes int) (string, error)
}
