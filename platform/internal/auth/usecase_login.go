package auth

import (
	"context"
	"strings"
	"time"
)

type LoginUser struct {
	Users           UserRepository
	Sessions        SessionRepository
	Passwords       PasswordVerifier
	Tokens          TokenHasher
	Clock           Clock
	IDs             IDGenerator
	SessionTTL      time.Duration
	RefreshTokenTTL time.Duration
}

type LoginInput struct {
	Email     string
	Password  string
	UserAgent string
	IP        string
}

type LoginResult struct {
	User             User
	Session          Session
	SessionToken     string
	RefreshToken     string
	CSRFToken        string
	RefreshExpiresAt time.Time
}

func (uc LoginUser) Execute(ctx context.Context, input LoginInput) (LoginResult, error) {
	user, err := uc.Users.FindByEmail(ctx, normalizeEmail(input.Email))
	if err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}
	if !user.IsActive() {
		return LoginResult{}, ErrInactiveUser
	}
	if !uc.Passwords.VerifyPassword(input.Password, user.PasswordHash) {
		return LoginResult{}, ErrInvalidCredentials
	}

	now := uc.Clock.Now().UTC()
	sessionID, err := uc.IDs.NewID("session")
	if err != nil {
		return LoginResult{}, err
	}
	refreshID, err := uc.IDs.NewID("refresh")
	if err != nil {
		return LoginResult{}, err
	}
	rawRefresh, err := uc.IDs.NewSecret(48)
	if err != nil {
		return LoginResult{}, err
	}
	rawSession, err := uc.IDs.NewSecret(48)
	if err != nil {
		return LoginResult{}, err
	}
	csrf, err := uc.IDs.NewSecret(32)
	if err != nil {
		return LoginResult{}, err
	}

	session := Session{
		ID:           SessionID(sessionID),
		UserID:       user.ID,
		TokenHash:    uc.Tokens.HashToken(rawSession),
		TokenVersion: user.TokenVersion,
		UserAgent:    input.UserAgent,
		IP:           input.IP,
		CreatedAt:    now,
		ExpiresAt:    now.Add(uc.SessionTTL),
	}
	refresh := RefreshToken{
		ID:        RefreshTokenID(refreshID),
		SessionID: session.ID,
		UserID:    user.ID,
		TokenHash: uc.Tokens.HashToken(rawRefresh),
		UserAgent: input.UserAgent,
		IP:        input.IP,
		ExpiresAt: now.Add(uc.RefreshTokenTTL),
		CreatedAt: now,
	}

	if err := uc.Sessions.CreateSession(ctx, session, refresh); err != nil {
		return LoginResult{}, err
	}
	if err := uc.Users.TouchLastLogin(ctx, user.ID, now); err != nil {
		return LoginResult{}, err
	}
	user.LastLoginAt = &now

	return LoginResult{
		User:             user,
		Session:          session,
		SessionToken:     rawSession,
		RefreshToken:     rawRefresh,
		CSRFToken:        csrf,
		RefreshExpiresAt: refresh.ExpiresAt,
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
