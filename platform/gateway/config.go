package gateway

import (
	"os"
	"strconv"
	"time"

	"github.com/besart951/code_links/platform/auth"
)

type Config struct {
	HTTPAddr        string
	DatabaseURL     string
	JWTSecret       string
	JWTIssuer       string
	InternalToken   string
	CookieDomain    string
	CookieSecure    bool
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func LoadConfig() Config {
	return Config{
		HTTPAddr:        env("HTTP_ADDR", ":8080"),
		DatabaseURL:     env("DATABASE_URL", ""),
		JWTSecret:       env("JWT_SECRET", "dev-only-change-me"),
		JWTIssuer:       env("JWT_ISSUER", "codelinks-platform"),
		InternalToken:   env("PLATFORM_INTERNAL_TOKEN", ""),
		CookieDomain:    env("COOKIE_DOMAIN", ""),
		CookieSecure:    envBool("COOKIE_SECURE", false),
		AccessTokenTTL:  envDuration("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL: envDuration("REFRESH_TOKEN_TTL", 30*24*time.Hour),
	}
}

func (c Config) AuthConfig() auth.Config {
	return auth.Config{
		JWTSecret:       c.JWTSecret,
		Issuer:          c.JWTIssuer,
		AccessTokenTTL:  c.AccessTokenTTL,
		RefreshTokenTTL: c.RefreshTokenTTL,
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
