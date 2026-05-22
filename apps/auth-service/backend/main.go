package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	ctx := context.Background()
	config := loadConfig()

	signer, err := newTokenSigner(config)
	if err != nil {
		log.Fatal(err)
	}

	store, cleanup, err := openStore(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	server := newServer(config, store, signer)
	httpServer := &http.Server{
		Addr:              ":" + config.Port,
		Handler:           server.routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("auth-service listening on :%s", config.Port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

type config struct {
	Port                 string
	DatabaseURL          string
	Environment          string
	Issuer               string
	Audience             string
	CookieDomain         string
	CookieSecure         bool
	PublicAuthBaseURL    string
	PublicFrontendURL    string
	AllowedOrigins       []string
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
	JWTKeyID             string
	JWTPrivateKeyPEM     string
	JWTPrivateKeyFile    string
	SMTPSecretKey        []byte
}

func loadConfig() config {
	environment := env("APP_ENV", "development")
	smtpSecretKey := smtpSecretKey(environment)

	return config{
		Port:                 env("PORT", "8080"),
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		Environment:          environment,
		Issuer:               env("JWT_ISSUER", "http://auth.codelinks.localhost"),
		Audience:             env("JWT_AUDIENCE", "codelinks-products"),
		CookieDomain:         os.Getenv("COOKIE_DOMAIN"),
		CookieSecure:         env("COOKIE_SECURE", environment) == "production" || os.Getenv("COOKIE_SECURE") == "true",
		PublicAuthBaseURL:    env("PUBLIC_AUTH_BASE_URL", "http://auth.codelinks.localhost"),
		PublicFrontendURL:    env("PUBLIC_AUTH_FRONTEND_URL", "http://auth.codelinks.localhost"),
		AllowedOrigins:       splitCSV(env("ALLOWED_ORIGINS", "http://code-links.codelinks.localhost,http://admin-link.codelinks.localhost,http://auth.codelinks.localhost,http://localhost:5173,http://localhost:5174,http://localhost:5175")),
		AccessTokenLifetime:  15 * time.Minute,
		RefreshTokenLifetime: 30 * 24 * time.Hour,
		JWTKeyID:             env("JWT_KEY_ID", "dev-key"),
		JWTPrivateKeyPEM:     os.Getenv("JWT_PRIVATE_KEY_PEM"),
		JWTPrivateKeyFile:    os.Getenv("JWT_PRIVATE_KEY_FILE"),
		SMTPSecretKey:        smtpSecretKey,
	}
}

func smtpSecretKey(environment string) []byte {
	raw := os.Getenv("SMTP_SECRET_KEY")
	if raw == "" {
		if environment == "production" {
			log.Fatal("SMTP_SECRET_KEY must be set in production")
		}
		sum := sha256.Sum256([]byte("codelinks-dev-smtp-secret"))
		return sum[:]
	}

	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err == nil && len(decoded) == 32 {
		return decoded
	}

	sum := sha256.Sum256([]byte(raw))
	return sum[:]
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}

	return items
}
