package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInactiveUser       = errors.New("inactive user")
	ErrTokenRevoked       = errors.New("token revoked")
	ErrSessionExpired     = errors.New("session expired")
)
