package productauth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email         string   `json:"email"`
	Name          string   `json:"name"`
	Status        string   `json:"status"`
	EmailVerified bool     `json:"emailVerified"`
	Licenses      []string `json:"licenses"`
	Roles         []string `json:"roles"`
	Permissions   []string `json:"permissions"`
	jwt.RegisteredClaims
}

func (c Claims) HasLicense(productID string) bool {
	for _, license := range c.Licenses {
		if license == productID {
			return true
		}
	}

	return false
}

func (c Claims) CurrentUser() CurrentUser {
	return CurrentUser{
		ID:            c.Subject,
		Email:         c.Email,
		Name:          c.Name,
		Status:        c.Status,
		EmailVerified: c.EmailVerified,
		Licenses:      append([]string(nil), c.Licenses...),
		Roles:         append([]string(nil), c.Roles...),
		Permissions:   append([]string(nil), c.Permissions...),
		ExpiresAt:     claimsExpiresAt(c),
	}
}

type CurrentUser struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	Status        string    `json:"status"`
	EmailVerified bool      `json:"emailVerified"`
	Licenses      []string  `json:"licenses"`
	Roles         []string  `json:"roles"`
	Permissions   []string  `json:"permissions"`
	ExpiresAt     time.Time `json:"-"`
}

func (u CurrentUser) HasLicense(productID string) bool {
	for _, license := range u.Licenses {
		if license == productID {
			return true
		}
	}

	return false
}

func (u CurrentUser) HasRole(role string) bool {
	for _, candidate := range u.Roles {
		if candidate == role {
			return true
		}
	}

	return false
}

func (u CurrentUser) HasPermission(permission string) bool {
	for _, candidate := range u.Permissions {
		if candidate == permission {
			return true
		}
	}

	return false
}

func claimsExpiresAt(claims Claims) time.Time {
	if claims.ExpiresAt == nil {
		return time.Time{}
	}
	return claims.ExpiresAt.Time
}
