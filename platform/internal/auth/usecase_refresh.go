package auth

import (
	"context"
	"time"
)

type RefreshSession struct {
	Users           UserRepository
	Sessions        SessionRepository
	Tokens          TokenHasher
	Clock           Clock
	IDs             IDGenerator
	RefreshTokenTTL time.Duration
}

type RefreshInput struct {
	RefreshToken string
	UserAgent    string
	IP           string
}

type RefreshResult struct {
	User             User
	Session          Session
	RefreshToken     string
	CSRFToken        string
	RefreshExpiresAt time.Time
}

func (uc RefreshSession) Execute(ctx context.Context, input RefreshInput) (RefreshResult, error) {
	now := uc.Clock.Now().UTC()
	hash := uc.Tokens.HashToken(input.RefreshToken)
	session, current, err := uc.Sessions.FindByRefreshTokenHash(ctx, hash)
	if err != nil {
		return RefreshResult{}, ErrUnauthorized
	}
	if !session.IsActive(now) {
		return RefreshResult{}, ErrSessionExpired
	}
	if !current.IsUsable(now) {
		return RefreshResult{}, ErrTokenRevoked
	}

	user, err := uc.Users.FindByID(ctx, session.UserID)
	if err != nil {
		return RefreshResult{}, ErrUnauthorized
	}
	if !user.IsActive() {
		return RefreshResult{}, ErrInactiveUser
	}
	if session.TokenVersion != user.TokenVersion {
		return RefreshResult{}, ErrTokenRevoked
	}

	refreshID, err := uc.IDs.NewID("refresh")
	if err != nil {
		return RefreshResult{}, err
	}
	nextRaw, err := uc.IDs.NewSecret(48)
	if err != nil {
		return RefreshResult{}, err
	}
	csrf, err := uc.IDs.NewSecret(32)
	if err != nil {
		return RefreshResult{}, err
	}
	next := RefreshToken{
		ID:        RefreshTokenID(refreshID),
		SessionID: session.ID,
		UserID:    user.ID,
		TokenHash: uc.Tokens.HashToken(nextRaw),
		UserAgent: input.UserAgent,
		IP:        input.IP,
		ExpiresAt: now.Add(uc.RefreshTokenTTL),
		CreatedAt: now,
	}
	if err := uc.Sessions.RotateRefreshToken(ctx, hash, next, now); err != nil {
		return RefreshResult{}, err
	}

	return RefreshResult{
		User:             user,
		Session:          session,
		RefreshToken:     nextRaw,
		CSRFToken:        csrf,
		RefreshExpiresAt: next.ExpiresAt,
	}, nil
}

type AuthenticateSession struct {
	Sessions SessionRepository
	Tokens   TokenHasher
	Clock    Clock
}

func (uc AuthenticateSession) Execute(ctx context.Context, rawSessionToken string) (AuthenticatedSession, error) {
	if rawSessionToken == "" {
		return AuthenticatedSession{}, ErrUnauthorized
	}
	authenticated, err := uc.Sessions.FindBySessionTokenHash(ctx, uc.Tokens.HashToken(rawSessionToken))
	if err != nil {
		return AuthenticatedSession{}, ErrUnauthorized
	}
	now := uc.Clock.Now().UTC()
	if !authenticated.Session.IsActive(now) {
		return AuthenticatedSession{}, ErrSessionExpired
	}
	if !authenticated.User.IsActive() {
		return AuthenticatedSession{}, ErrInactiveUser
	}
	if authenticated.Session.TokenVersion != authenticated.User.TokenVersion {
		return AuthenticatedSession{}, ErrTokenRevoked
	}
	return authenticated, nil
}

type LogoutSession struct {
	Sessions SessionRepository
	Clock    Clock
}

func (uc LogoutSession) Execute(ctx context.Context, sessionID SessionID) error {
	if sessionID == "" {
		return nil
	}
	return uc.Sessions.RevokeSession(ctx, sessionID, uc.Clock.Now().UTC())
}
