package auth

import "time"

type UserID string
type SessionID string
type RefreshTokenID string
type TokenVersion int

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

type User struct {
	ID           UserID
	Email        string
	DisplayName  string
	PasswordHash string
	Status       UserStatus
	TokenVersion TokenVersion
	CreatedAt    time.Time
	LastLoginAt  *time.Time
}

func (u User) IsActive() bool {
	return u.Status == UserStatusActive
}

type Session struct {
	ID           SessionID
	UserID       UserID
	TokenHash    string
	TokenVersion TokenVersion
	UserAgent    string
	IP           string
	CreatedAt    time.Time
	LastSeenAt   *time.Time
	RevokedAt    *time.Time
	ExpiresAt    time.Time
}

func (s Session) IsActive(now time.Time) bool {
	return s.RevokedAt == nil && s.ExpiresAt.After(now)
}

func (s *Session) Revoke(at time.Time) {
	if s.RevokedAt == nil {
		s.RevokedAt = &at
	}
}

type AuthenticatedSession struct {
	User    User
	Session Session
}

type RefreshToken struct {
	ID        RefreshTokenID
	SessionID SessionID
	UserID    UserID
	TokenHash string
	UserAgent string
	IP        string
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time
}

func (t RefreshToken) IsUsable(now time.Time) bool {
	return t.RevokedAt == nil && t.ExpiresAt.After(now)
}
