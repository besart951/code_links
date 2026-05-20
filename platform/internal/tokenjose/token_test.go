package tokenjose

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"testing"
	"time"

	"github.com/besart951/code_links/platform/internal/access"
)

func TestNestedJWTRoundTripAndWrongKeyFailure(t *testing.T) {
	signPub, signPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	audiencePriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	wrongPriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 5, 20, 10, 0, 0, 0, time.UTC)
	registry := fakeKeys{
		signing: SigningKey{KeyID: "sig-1", PrivateKey: signPriv, PublicKey: signPub},
		Encryption: map[access.ProductKey]AudienceEncryptionKey{
			"infra_link": {Audience: "infra_link", KeyID: "infra-enc-1", PublicKey: &audiencePriv.PublicKey},
		},
		Decryption: map[access.ProductKey]AudienceDecryptionKey{
			"infra_link": {Audience: "infra_link", KeyID: "infra-enc-1", PrivateKey: audiencePriv},
		},
	}
	snapshot := access.AuthorizationSnapshot{
		Issuer:              "https://auth.codelinks.ch",
		Subject:             "user_1",
		Audience:            "infra_link",
		TenantID:            "tenant_1",
		TenantType:          access.TenantTypeCompany,
		SessionID:           "session_1",
		TokenVersion:        3,
		EntitlementsVersion: 12,
		Roles:               []access.RoleKey{"owner", "admin"},
		Permissions:         []access.PermissionKey{"infra_link.project.read"},
		Product: access.ProductSnapshot{
			ProductKey: "infra_link",
			PlanKey:    "business",
			Access:     true,
			Entitlements: []access.Entitlement{{
				TenantID:   "tenant_1",
				ProductKey: "infra_link",
				FeatureKey: "infra_link.project.read",
				Enabled:    true,
			}},
			Limits: []access.FeatureLimit{{
				TenantID:   "tenant_1",
				ProductKey: "infra_link",
				FeatureKey: "infra_link.project.read",
				LimitKey:   "max_projects",
				Value:      100,
			}},
		},
		IssuedAt:  now,
		NotBefore: now,
		ExpiresAt: now.Add(10 * time.Minute),
		JWTID:     "token_1",
	}

	issued, err := Issuer{Keys: registry}.IssueAccessToken(context.Background(), snapshot)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Validator{Keys: registry, Now: func() time.Time { return now.Add(time.Minute) }}.
		ValidateAccessToken(context.Background(), issued.Value, "infra_link")
	if err != nil {
		t.Fatal(err)
	}
	if got.Subject != "user_1" || got.Audience != "infra_link" || got.TokenVersion != 3 || got.EntitlementsVersion != 12 {
		t.Fatalf("unexpected snapshot %#v", got)
	}
	if len(got.Product.Entitlements) != 1 || got.Product.Entitlements[0].FeatureKey != "infra_link.project.read" {
		t.Fatalf("expected feature snapshot to round trip, got %#v", got.Product.Entitlements)
	}
	if len(got.Product.Limits) != 1 || got.Product.Limits[0].Value != 100 {
		t.Fatalf("expected limit snapshot to round trip, got %#v", got.Product.Limits)
	}

	wrongRegistry := registry
	wrongRegistry.Decryption = map[access.ProductKey]AudienceDecryptionKey{
		"infra_link": {Audience: "infra_link", KeyID: "infra-enc-1", PrivateKey: wrongPriv},
	}
	if _, err := (Validator{Keys: wrongRegistry, Now: func() time.Time { return now.Add(time.Minute) }}).
		ValidateAccessToken(context.Background(), issued.Value, "infra_link"); err == nil {
		t.Fatal("expected wrong decryption key to fail")
	}
}

func TestValidateAccessTokenRejectsInvalidRequiredClaims(t *testing.T) {
	now := time.Date(2026, 5, 20, 10, 0, 0, 0, time.UTC)
	registry := newTestKeys(t)
	base := access.AuthorizationSnapshot{
		Issuer:              "https://auth.codelinks.ch",
		Subject:             "user_1",
		Audience:            "infra_link",
		TenantID:            "tenant_1",
		TenantType:          access.TenantTypeCompany,
		SessionID:           "session_1",
		TokenVersion:        3,
		EntitlementsVersion: 12,
		Product:             access.ProductSnapshot{ProductKey: "infra_link", Access: true},
		IssuedAt:            now,
		NotBefore:           now,
		ExpiresAt:           now.Add(10 * time.Minute),
		JWTID:               "token_1",
	}

	tests := []struct {
		name      string
		snapshot  access.AuthorizationSnapshot
		validator Validator
	}{
		{
			name:     "wrong issuer",
			snapshot: base,
			validator: Validator{
				Keys:           registry,
				ExpectedIssuer: "https://other.example",
				Now:            func() time.Time { return now.Add(time.Minute) },
			},
		},
		{
			name: "missing subject",
			snapshot: func() access.AuthorizationSnapshot {
				s := base
				s.Subject = ""
				return s
			}(),
			validator: Validator{Keys: registry, ExpectedIssuer: base.Issuer, Now: func() time.Time { return now.Add(time.Minute) }},
		},
		{
			name: "missing jti",
			snapshot: func() access.AuthorizationSnapshot {
				s := base
				s.JWTID = ""
				return s
			}(),
			validator: Validator{Keys: registry, ExpectedIssuer: base.Issuer, Now: func() time.Time { return now.Add(time.Minute) }},
		},
		{
			name: "expired",
			snapshot: func() access.AuthorizationSnapshot {
				s := base
				s.ExpiresAt = now.Add(-time.Minute)
				return s
			}(),
			validator: Validator{Keys: registry, ExpectedIssuer: base.Issuer, Now: func() time.Time { return now.Add(time.Minute) }},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issued, err := Issuer{Keys: registry}.IssueAccessToken(context.Background(), tt.snapshot)
			if err != nil {
				t.Fatal(err)
			}
			_, err = tt.validator.ValidateAccessToken(context.Background(), issued.Value, "infra_link")
			if !errors.Is(err, ErrInvalidToken) {
				t.Fatalf("expected ErrInvalidToken, got %v", err)
			}
		})
	}
}

func newTestKeys(t *testing.T) fakeKeys {
	t.Helper()
	signPub, signPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	audiencePriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return fakeKeys{
		signing: SigningKey{KeyID: "sig-1", PrivateKey: signPriv, PublicKey: signPub},
		Encryption: map[access.ProductKey]AudienceEncryptionKey{
			"infra_link": {Audience: "infra_link", KeyID: "infra-enc-1", PublicKey: &audiencePriv.PublicKey},
		},
		Decryption: map[access.ProductKey]AudienceDecryptionKey{
			"infra_link": {Audience: "infra_link", KeyID: "infra-enc-1", PrivateKey: audiencePriv},
		},
	}
}

type fakeKeys struct {
	signing    SigningKey
	Encryption map[access.ProductKey]AudienceEncryptionKey
	Decryption map[access.ProductKey]AudienceDecryptionKey
}

func (k fakeKeys) SigningKey(ctx context.Context) (SigningKey, error) {
	return k.signing, nil
}

func (k fakeKeys) SigningPublicKey(ctx context.Context, keyID string) (ed25519.PublicKey, error) {
	if keyID != k.signing.KeyID {
		return nil, ErrKeyNotFound
	}
	return k.signing.PublicKey, nil
}

func (k fakeKeys) AudienceEncryptionKey(ctx context.Context, audience access.ProductKey) (AudienceEncryptionKey, error) {
	key, ok := k.Encryption[audience]
	if !ok {
		return AudienceEncryptionKey{}, ErrKeyNotFound
	}
	return key, nil
}

func (k fakeKeys) AudienceDecryptionKey(ctx context.Context, audience access.ProductKey, keyID string) (AudienceDecryptionKey, error) {
	key, ok := k.Decryption[audience]
	if !ok || key.KeyID != keyID {
		return AudienceDecryptionKey{}, ErrKeyNotFound
	}
	return key, nil
}
