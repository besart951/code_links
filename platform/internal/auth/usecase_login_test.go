package auth

import (
	"context"
	"testing"
	"time"
)

func TestLoginUserCreatesSessionAndRefreshToken(t *testing.T) {
	now := time.Date(2026, 5, 20, 10, 0, 0, 0, time.UTC)
	users := &fakeUsers{
		byEmail: map[string]User{
			"owner@example.com": {
				ID:           "user_1",
				Email:        "owner@example.com",
				PasswordHash: "hash",
				Status:       UserStatusActive,
				TokenVersion: 3,
			},
		},
	}
	sessions := &fakeSessions{}
	uc := LoginUser{
		Users:           users,
		Sessions:        sessions,
		Passwords:       acceptingPasswords{},
		Tokens:          staticHasher{},
		Clock:           authClock{now: now},
		IDs:             &sequenceIDs{ids: []string{"session_1", "refresh_1"}, secrets: []string{"refresh_raw", "session_raw", "csrf_raw"}},
		SessionTTL:      time.Hour,
		RefreshTokenTTL: 24 * time.Hour,
	}

	result, err := uc.Execute(context.Background(), LoginInput{
		Email:     " OWNER@example.com ",
		Password:  "secret",
		UserAgent: "test-agent",
		IP:        "127.0.0.1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Session.ID != "session_1" || result.Session.TokenVersion != 3 {
		t.Fatalf("unexpected session %#v", result.Session)
	}
	if result.RefreshToken != "refresh_raw" || result.SessionToken != "session_raw" || result.CSRFToken != "csrf_raw" {
		t.Fatalf("unexpected returned secrets %#v", result)
	}
	if sessions.refresh.TokenHash != "hashed:refresh_raw" {
		t.Fatalf("refresh token was not hashed: %#v", sessions.refresh)
	}
	if sessions.session.TokenHash != "hashed:session_raw" {
		t.Fatalf("session token was not hashed: %#v", sessions.session)
	}
	if !users.touched.Equal(now) {
		t.Fatalf("last login not touched at fixed time: %s", users.touched)
	}
}

type fakeUsers struct {
	byEmail map[string]User
	touched time.Time
}

func (r *fakeUsers) FindByEmail(ctx context.Context, email string) (User, error) {
	user, ok := r.byEmail[email]
	if !ok {
		return User{}, ErrInvalidCredentials
	}
	return user, nil
}

func (r *fakeUsers) FindByID(ctx context.Context, userID UserID) (User, error) {
	for _, user := range r.byEmail {
		if user.ID == userID {
			return user, nil
		}
	}
	return User{}, ErrUnauthorized
}

func (r *fakeUsers) TouchLastLogin(ctx context.Context, userID UserID, at time.Time) error {
	r.touched = at
	return nil
}

type fakeSessions struct {
	session Session
	refresh RefreshToken
}

func (r *fakeSessions) CreateSession(ctx context.Context, session Session, refresh RefreshToken) error {
	r.session = session
	r.refresh = refresh
	return nil
}

func (r *fakeSessions) FindBySessionTokenHash(ctx context.Context, tokenHash string) (AuthenticatedSession, error) {
	return AuthenticatedSession{}, ErrUnauthorized
}

func (r *fakeSessions) FindByRefreshTokenHash(ctx context.Context, tokenHash string) (Session, RefreshToken, error) {
	return r.session, r.refresh, nil
}

func (r *fakeSessions) RotateRefreshToken(ctx context.Context, oldTokenHash string, next RefreshToken, revokedAt time.Time) error {
	r.refresh = next
	return nil
}

func (r *fakeSessions) RevokeSession(ctx context.Context, sessionID SessionID, revokedAt time.Time) error {
	r.session.RevokedAt = &revokedAt
	return nil
}

type acceptingPasswords struct{}

func (acceptingPasswords) VerifyPassword(password, passwordHash string) bool {
	return true
}

type staticHasher struct{}

func (staticHasher) HashToken(raw string) string {
	return "hashed:" + raw
}

type authClock struct {
	now time.Time
}

func (c authClock) Now() time.Time {
	return c.now
}

type sequenceIDs struct {
	ids         []string
	secrets     []string
	idIndex     int
	secretIndex int
}

func (g *sequenceIDs) NewID(prefix string) (string, error) {
	value := g.ids[g.idIndex]
	g.idIndex++
	return value, nil
}

func (g *sequenceIDs) NewSecret(bytes int) (string, error) {
	value := g.secrets[g.secretIndex]
	g.secretIndex++
	return value, nil
}
