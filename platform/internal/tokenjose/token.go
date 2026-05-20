package tokenjose

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/besart951/code_links/platform/internal/access"
	jose "github.com/go-jose/go-jose/v4"
)

var ErrKeyNotFound = errors.New("token key not found")
var ErrInvalidToken = errors.New("invalid access token")

type SigningKey struct {
	KeyID      string
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

type AudienceEncryptionKey struct {
	Audience  access.ProductKey
	KeyID     string
	PublicKey *ecdsa.PublicKey
}

type AudienceDecryptionKey struct {
	Audience   access.ProductKey
	KeyID      string
	PrivateKey *ecdsa.PrivateKey
}

type KeyResolver interface {
	SigningKey(ctx context.Context) (SigningKey, error)
	SigningPublicKey(ctx context.Context, keyID string) (ed25519.PublicKey, error)
	AudienceEncryptionKey(ctx context.Context, audience access.ProductKey) (AudienceEncryptionKey, error)
	AudienceDecryptionKey(ctx context.Context, audience access.ProductKey, keyID string) (AudienceDecryptionKey, error)
}

type Issuer struct {
	Keys KeyResolver
}

func (i Issuer) IssueAccessToken(ctx context.Context, snapshot access.AuthorizationSnapshot) (access.IssuedToken, error) {
	signingKey, err := i.Keys.SigningKey(ctx)
	if err != nil {
		return access.IssuedToken{}, err
	}
	claims := claimsFromSnapshot(snapshot)
	payload, err := json.Marshal(claims)
	if err != nil {
		return access.IssuedToken{}, err
	}

	signerOpts := new(jose.SignerOptions).
		WithType("at+jwt").
		WithHeader("kid", signingKey.KeyID)
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.EdDSA,
		Key:       signingKey.PrivateKey,
	}, signerOpts)
	if err != nil {
		return access.IssuedToken{}, err
	}
	jws, err := signer.Sign(payload)
	if err != nil {
		return access.IssuedToken{}, err
	}
	inner, err := jws.CompactSerialize()
	if err != nil {
		return access.IssuedToken{}, err
	}

	audienceKey, err := i.Keys.AudienceEncryptionKey(ctx, snapshot.Audience)
	if err != nil {
		return access.IssuedToken{}, err
	}
	encrypter, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.ECDH_ES,
			Key:       audienceKey.PublicKey,
			KeyID:     audienceKey.KeyID,
		},
		new(jose.EncrypterOptions).WithContentType("JWT"),
	)
	if err != nil {
		return access.IssuedToken{}, err
	}
	jwe, err := encrypter.Encrypt([]byte(inner))
	if err != nil {
		return access.IssuedToken{}, err
	}
	raw, err := jwe.CompactSerialize()
	if err != nil {
		return access.IssuedToken{}, err
	}
	return access.IssuedToken{
		Value:     raw,
		ExpiresAt: snapshot.ExpiresAt,
		JWTID:     snapshot.JWTID,
	}, nil
}

type Validator struct {
	Keys           KeyResolver
	ExpectedIssuer string
	Now            func() time.Time
}

func (v Validator) ValidateAccessToken(ctx context.Context, raw string, audience access.ProductKey) (access.AuthorizationSnapshot, error) {
	jwe, err := jose.ParseEncrypted(raw, []jose.KeyAlgorithm{jose.ECDH_ES}, []jose.ContentEncryption{jose.A256GCM})
	if err != nil {
		return access.AuthorizationSnapshot{}, err
	}
	keyID := jwe.Header.KeyID
	decryptionKey, err := v.Keys.AudienceDecryptionKey(ctx, audience, keyID)
	if err != nil {
		return access.AuthorizationSnapshot{}, err
	}
	inner, err := jwe.Decrypt(decryptionKey.PrivateKey)
	if err != nil {
		return access.AuthorizationSnapshot{}, err
	}
	jws, err := jose.ParseSigned(string(inner), []jose.SignatureAlgorithm{jose.EdDSA})
	if err != nil {
		return access.AuthorizationSnapshot{}, err
	}
	signingKID := jws.Signatures[0].Header.KeyID
	publicKey, err := v.Keys.SigningPublicKey(ctx, signingKID)
	if err != nil {
		return access.AuthorizationSnapshot{}, err
	}
	payload, err := jws.Verify(publicKey)
	if err != nil {
		return access.AuthorizationSnapshot{}, err
	}
	var c tokenClaims
	if err := json.Unmarshal(payload, &c); err != nil {
		return access.AuthorizationSnapshot{}, err
	}
	snapshot := c.snapshot()
	if v.ExpectedIssuer != "" && snapshot.Issuer != v.ExpectedIssuer {
		return access.AuthorizationSnapshot{}, ErrInvalidToken
	}
	if snapshot.Audience != audience {
		return access.AuthorizationSnapshot{}, access.ErrAudienceMismatch
	}
	if snapshot.Subject == "" ||
		snapshot.SessionID == "" ||
		snapshot.JWTID == "" ||
		snapshot.TokenVersion <= 0 ||
		snapshot.EntitlementsVersion <= 0 {
		return access.AuthorizationSnapshot{}, ErrInvalidToken
	}
	now := time.Now().UTC()
	if v.Now != nil {
		now = v.Now().UTC()
	}
	if now.Before(snapshot.NotBefore) || !snapshot.ExpiresAt.After(now) {
		return access.AuthorizationSnapshot{}, ErrInvalidToken
	}
	return snapshot, nil
}

type tokenClaims struct {
	Issuer              string        `json:"iss"`
	Subject             string        `json:"sub"`
	Audience            string        `json:"aud"`
	IssuedAt            int64         `json:"iat"`
	NotBefore           int64         `json:"nbf"`
	ExpiresAt           int64         `json:"exp"`
	JWTID               string        `json:"jti"`
	TenantID            string        `json:"tenant_id"`
	TenantType          string        `json:"tenant_type"`
	SessionID           string        `json:"session_id"`
	TokenVersion        int           `json:"token_version"`
	EntitlementsVersion int           `json:"entitlements_version"`
	Roles               []string      `json:"roles"`
	Permissions         []string      `json:"permissions"`
	Product             productClaims `json:"product"`
}

type productClaims struct {
	Key      string           `json:"key"`
	Plan     string           `json:"plan,omitempty"`
	Access   bool             `json:"access"`
	Features map[string]bool  `json:"features,omitempty"`
	Limits   map[string]int64 `json:"limits,omitempty"`
}

func claimsFromSnapshot(snapshot access.AuthorizationSnapshot) tokenClaims {
	roles := make([]string, 0, len(snapshot.Roles))
	for _, role := range snapshot.Roles {
		roles = append(roles, string(role))
	}
	permissions := make([]string, 0, len(snapshot.Permissions))
	for _, permission := range snapshot.Permissions {
		permissions = append(permissions, string(permission))
	}
	features := make(map[string]bool, len(snapshot.Product.Entitlements))
	for _, item := range snapshot.Product.Entitlements {
		features[string(item.FeatureKey)] = item.Enabled
	}
	limits := make(map[string]int64, len(snapshot.Product.Limits))
	for _, item := range snapshot.Product.Limits {
		limits[string(item.FeatureKey)+"."+string(item.LimitKey)] = item.Value
	}
	return tokenClaims{
		Issuer:              snapshot.Issuer,
		Subject:             string(snapshot.Subject),
		Audience:            string(snapshot.Audience),
		IssuedAt:            snapshot.IssuedAt.Unix(),
		NotBefore:           snapshot.NotBefore.Unix(),
		ExpiresAt:           snapshot.ExpiresAt.Unix(),
		JWTID:               snapshot.JWTID,
		TenantID:            string(snapshot.TenantID),
		TenantType:          string(snapshot.TenantType),
		SessionID:           string(snapshot.SessionID),
		TokenVersion:        int(snapshot.TokenVersion),
		EntitlementsVersion: int(snapshot.EntitlementsVersion),
		Roles:               roles,
		Permissions:         permissions,
		Product: productClaims{
			Key:      string(snapshot.Product.ProductKey),
			Plan:     string(snapshot.Product.PlanKey),
			Access:   snapshot.Product.Access,
			Features: features,
			Limits:   limits,
		},
	}
}

func (c tokenClaims) snapshot() access.AuthorizationSnapshot {
	roles := make([]access.RoleKey, 0, len(c.Roles))
	for _, role := range c.Roles {
		roles = append(roles, access.RoleKey(role))
	}
	permissions := make([]access.PermissionKey, 0, len(c.Permissions))
	for _, permission := range c.Permissions {
		permissions = append(permissions, access.PermissionKey(permission))
	}
	entitlements := make([]access.Entitlement, 0, len(c.Product.Features))
	for feature, enabled := range c.Product.Features {
		entitlements = append(entitlements, access.Entitlement{
			TenantID:   access.TenantID(c.TenantID),
			ProductKey: access.ProductKey(c.Product.Key),
			FeatureKey: access.FeatureKey(feature),
			Enabled:    enabled,
		})
	}
	limits := make([]access.FeatureLimit, 0, len(c.Product.Limits))
	for combinedKey, value := range c.Product.Limits {
		splitAt := strings.LastIndex(combinedKey, ".")
		if splitAt <= 0 || splitAt == len(combinedKey)-1 {
			continue
		}
		feature := combinedKey[:splitAt]
		limit := combinedKey[splitAt+1:]
		limits = append(limits, access.FeatureLimit{
			TenantID:   access.TenantID(c.TenantID),
			ProductKey: access.ProductKey(c.Product.Key),
			FeatureKey: access.FeatureKey(feature),
			LimitKey:   access.LimitKey(limit),
			Value:      value,
		})
	}
	return access.AuthorizationSnapshot{
		Issuer:              c.Issuer,
		Subject:             access.UserID(c.Subject),
		Audience:            access.ProductKey(c.Audience),
		TenantID:            access.TenantID(c.TenantID),
		TenantType:          access.TenantType(c.TenantType),
		SessionID:           access.SessionID(c.SessionID),
		TokenVersion:        access.TokenVersion(c.TokenVersion),
		EntitlementsVersion: access.EntitlementsVersion(c.EntitlementsVersion),
		Roles:               roles,
		Permissions:         permissions,
		Product: access.ProductSnapshot{
			ProductKey:   access.ProductKey(c.Product.Key),
			PlanKey:      access.PlanKey(c.Product.Plan),
			Access:       c.Product.Access,
			Entitlements: entitlements,
			Limits:       limits,
		},
		IssuedAt:  time.Unix(c.IssuedAt, 0).UTC(),
		NotBefore: time.Unix(c.NotBefore, 0).UTC(),
		ExpiresAt: time.Unix(c.ExpiresAt, 0).UTC(),
		JWTID:     c.JWTID,
	}
}

var _ access.TokenIssuer = (*Issuer)(nil)
var _ access.TokenValidator = (*Validator)(nil)
