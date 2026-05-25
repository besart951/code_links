package config

import (
	"crypto/sha256"
	"strings"
	"testing"
)

func TestValidateProductionConfigRejectsUnsafeSettings(t *testing.T) {
	err := ValidateProductionConfig(Config{
		Environment:        "production",
		CookieSecure:       false,
		CookieDomain:       "codelinks.localhost",
		PublicAuthBaseURL:  "http://auth.codelinks.localhost",
		PublicFrontendURL:  "http://localhost:5173",
		AllowedOrigins:     []string{"http://admin-link.codelinks.localhost"},
		EnableMockPurchase: true,
		TrustedProxyCIDRs:  []string{"0.0.0.0/0"},
	})
	if err == nil {
		t.Fatal("expected unsafe production config to fail")
	}
	message := err.Error()
	for _, expected := range []string{
		"JWT_PRIVATE_KEY_PEM or JWT_PRIVATE_KEY_FILE",
		"SMTP_SECRET_KEY",
		"COOKIE_SECURE",
		"COOKIE_DOMAIN",
		"PUBLIC_AUTH_BASE_URL",
		"PUBLIC_AUTH_FRONTEND_URL",
		"ALLOWED_ORIGINS",
		"ENABLE_MOCK_PURCHASE",
		"TRUSTED_PROXY_CIDRS",
	} {
		if !strings.Contains(message, expected) {
			t.Fatalf("expected %q in %q", expected, message)
		}
	}
}

func TestValidateProductionConfigAcceptsLockedDownSettings(t *testing.T) {
	secret := sha256.Sum256([]byte("smtp-secret"))
	if err := ValidateProductionConfig(Config{
		Environment:        "production",
		CookieSecure:       true,
		CookieDomain:       ".codelinks.dev",
		PublicAuthBaseURL:  "https://auth.codelinks.dev",
		PublicFrontendURL:  "https://auth.codelinks.dev",
		AllowedOrigins:     []string{"https://admin.codelinks.dev"},
		JWTPrivateKeyPEM:   "pem",
		SMTPSecretKey:      secret[:],
		TrustedProxyCIDRs:  []string{"10.0.0.0/8"},
		EnableMockPurchase: false,
	}); err != nil {
		t.Fatal(err)
	}
}
