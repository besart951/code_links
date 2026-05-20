package access

import (
	"context"
	"time"
)

type Repository interface {
	SessionAccess(ctx context.Context, sessionID SessionID) (SessionAccess, error)
	TenantAccess(ctx context.Context, tenantID TenantID, product ProductKey) (TenantAccess, error)
	MemberAccess(ctx context.Context, userID UserID, tenantID TenantID, product ProductKey) (MemberAccess, error)
}

type TokenIssuer interface {
	IssueAccessToken(ctx context.Context, snapshot AuthorizationSnapshot) (IssuedToken, error)
}

type TokenValidator interface {
	ValidateAccessToken(ctx context.Context, raw string, audience ProductKey) (AuthorizationSnapshot, error)
}

type IssuedToken struct {
	Value     string
	ExpiresAt time.Time
	JWTID     string
}

type Clock interface {
	Now() time.Time
}

type IDGenerator interface {
	NewID(prefix string) (string, error)
}
