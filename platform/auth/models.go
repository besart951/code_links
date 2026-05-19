package auth

import "time"

type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	DisplayName  string     `json:"display_name"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	UserAgent *string
	IP        *string
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time
}

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	CSRFToken        string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	User             User
}
