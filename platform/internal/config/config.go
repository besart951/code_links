package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr               string
	DatabaseURL            string
	Issuer                 string
	CookieDomain           string
	CookieSecure           bool
	AccessTokenTTL         time.Duration
	SessionTTL             time.Duration
	RefreshTokenTTL        time.Duration
	SigningPrivateKeyPath  string
	SigningKeyID           string
	AudiencePublicKeyPaths map[string]string
	ProductClientTokens    map[string]string
}

func Load() Config {
	return Config{
		HTTPAddr:               env("HTTP_ADDR", ":8080"),
		DatabaseURL:            env("DATABASE_URL", ""),
		Issuer:                 env("TOKEN_ISSUER", env("JWT_ISSUER", "https://auth.codelinks.ch")),
		CookieDomain:           env("COOKIE_DOMAIN", ""),
		CookieSecure:           envBool("COOKIE_SECURE", true),
		AccessTokenTTL:         envDuration("ACCESS_TOKEN_TTL", 10*time.Minute),
		SessionTTL:             envDuration("SESSION_TTL", 30*24*time.Hour),
		RefreshTokenTTL:        envDuration("REFRESH_TOKEN_TTL", 30*24*time.Hour),
		SigningPrivateKeyPath:  env("TOKEN_SIGNING_PRIVATE_KEY_PATH", ""),
		SigningKeyID:           env("TOKEN_SIGNING_KID", "platform-ed25519-1"),
		AudiencePublicKeyPaths: envMap("AUDIENCE_PUBLIC_KEY_PATHS"),
		ProductClientTokens:    envMap("PRODUCT_CLIENT_TOKENS"),
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

func envMap(key string) map[string]string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return map[string]string{}
	}
	result := map[string]string{}
	for _, item := range strings.Split(raw, ",") {
		name, value, ok := strings.Cut(item, "=")
		if !ok {
			continue
		}
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		if name != "" && value != "" {
			result[name] = value
		}
	}
	return result
}
