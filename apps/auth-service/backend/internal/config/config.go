package config

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/netip"
	"net/url"
	"os"
	"strings"
	"time"
)

type Config struct {
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
	TrustedProxyCIDRs    []string
	EnableMockPurchase   bool
	RuntimeLogFile       string
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
	JWTKeyID             string
	JWTPrivateKeyPEM     string
	JWTPrivateKeyFile    string
	SMTPSecretKey        []byte
}

func Load() Config {
	environment := env("APP_ENV", "development")
	smtpSecretKey := smtpSecretKey(environment)
	jwtPrivateKeyPEM := os.Getenv("JWT_PRIVATE_KEY_PEM")
	jwtPrivateKeyFile := os.Getenv("JWT_PRIVATE_KEY_FILE")
	if environment == "production" && jwtPrivateKeyPEM == "" && jwtPrivateKeyFile == "" {
		log.Fatal("fatal: JWT_PRIVATE_KEY_PEM or JWT_PRIVATE_KEY_FILE must be set in production")
	}

	cfg := Config{
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
		TrustedProxyCIDRs:    splitCSV(os.Getenv("TRUSTED_PROXY_CIDRS")),
		EnableMockPurchase:   os.Getenv("ENABLE_MOCK_PURCHASE") == "true",
		RuntimeLogFile:       env("AUTH_SERVICE_LOG_FILE", "logs/auth-service.log"),
		AccessTokenLifetime:  15 * time.Minute,
		RefreshTokenLifetime: 30 * 24 * time.Hour,
		JWTKeyID:             env("JWT_KEY_ID", "dev-key"),
		JWTPrivateKeyPEM:     jwtPrivateKeyPEM,
		JWTPrivateKeyFile:    jwtPrivateKeyFile,
		SMTPSecretKey:        smtpSecretKey,
	}
	if err := ValidateProductionConfig(cfg); err != nil {
		log.Fatal("fatal: " + err.Error())
	}

	return cfg
}

func smtpSecretKey(environment string) []byte {
	raw := os.Getenv("SMTP_SECRET_KEY")
	if raw == "" {
		if environment == "production" {
			log.Fatal("fatal: SMTP_SECRET_KEY must be set in production")
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

func ValidateProductionConfig(cfg Config) error {
	if cfg.Environment != "production" {
		return nil
	}

	problems := []string{}
	if cfg.JWTPrivateKeyPEM == "" && cfg.JWTPrivateKeyFile == "" {
		problems = append(problems, "JWT_PRIVATE_KEY_PEM or JWT_PRIVATE_KEY_FILE must be set")
	}
	if len(cfg.SMTPSecretKey) == 0 {
		problems = append(problems, "SMTP_SECRET_KEY must be set")
	}
	if !cfg.CookieSecure {
		problems = append(problems, "COOKIE_SECURE must resolve true")
	}
	if containsLocalhost(cfg.CookieDomain) {
		problems = append(problems, "COOKIE_DOMAIN must not contain localhost")
	}
	for name, value := range map[string]string{
		"PUBLIC_AUTH_FRONTEND_URL": cfg.PublicFrontendURL,
		"PUBLIC_AUTH_BASE_URL":     cfg.PublicAuthBaseURL,
	} {
		if containsLocalhost(value) {
			problems = append(problems, name+" must not contain localhost")
		}
		if value != "" && !isHTTPSURL(value) {
			problems = append(problems, name+" must be HTTPS")
		}
	}
	for _, origin := range cfg.AllowedOrigins {
		if containsLocalhost(origin) {
			problems = append(problems, "ALLOWED_ORIGINS must not contain localhost")
		}
		if !isHTTPSURL(origin) {
			problems = append(problems, "ALLOWED_ORIGINS must be HTTPS")
		}
	}
	if cfg.EnableMockPurchase {
		problems = append(problems, "ENABLE_MOCK_PURCHASE must be false")
	}
	for _, cidr := range cfg.TrustedProxyCIDRs {
		if trustedProxyCIDRIsWildcard(cidr) {
			problems = append(problems, "TRUSTED_PROXY_CIDRS must not contain wildcard CIDR")
		}
	}

	if len(problems) > 0 {
		return fmt.Errorf("invalid production config: %s", strings.Join(problems, "; "))
	}
	return nil
}

func containsLocalhost(value string) bool {
	lower := strings.ToLower(value)
	return strings.Contains(lower, "localhost") || strings.Contains(lower, "codelinks.localhost")
}

func isHTTPSURL(value string) bool {
	parsed, err := url.Parse(value)
	return err == nil && parsed.Scheme == "https" && parsed.Host != ""
}

func trustedProxyCIDRIsWildcard(value string) bool {
	prefix, err := netip.ParsePrefix(value)
	if err != nil {
		return false
	}
	return prefix.Bits() == 0
}
