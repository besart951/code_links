package productauth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestRemoteValidatorValidatesTokenAndLicense(t *testing.T) {
	privateKey := newTestKey(t)
	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(jwksDocument{Keys: []jwk{testJWK("test-key", &privateKey.PublicKey)}})
	}))
	defer jwksServer.Close()

	validator, err := NewRemoteValidator(context.Background(), RemoteValidatorConfig{
		JWKSURL:   jwksServer.URL,
		Issuer:    "http://auth.codelinks.localhost",
		Audience:  "codelinks-products",
		ProductID: "infra-link",
	})
	if err != nil {
		t.Fatal(err)
	}

	token := signToken(t, privateKey, "test-key", []string{"infra-link"})
	claims, err := validator.ValidateToken(token)
	if err != nil {
		t.Fatal(err)
	}

	if claims.Subject != "user-1" {
		t.Fatalf("expected subject user-1, got %q", claims.Subject)
	}
	if !claims.HasLicense("infra-link") {
		t.Fatal("expected infra-link license")
	}
}

func TestRemoteValidatorRejectsMissingLicense(t *testing.T) {
	privateKey := newTestKey(t)
	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(jwksDocument{Keys: []jwk{testJWK("test-key", &privateKey.PublicKey)}})
	}))
	defer jwksServer.Close()

	validator, err := NewRemoteValidator(context.Background(), RemoteValidatorConfig{
		JWKSURL:   jwksServer.URL,
		Issuer:    "http://auth.codelinks.localhost",
		Audience:  "codelinks-products",
		ProductID: "infra-link",
	})
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	request.Header.Set("Authorization", "Bearer "+signToken(t, privateKey, "test-key", []string{"planer-link"}))
	recorder := httptest.NewRecorder()

	validator.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", recorder.Code)
	}
}

func TestMiddlewareAddsCurrentUserToContext(t *testing.T) {
	privateKey := newTestKey(t)
	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(jwksDocument{Keys: []jwk{testJWK("test-key", &privateKey.PublicKey)}})
	}))
	defer jwksServer.Close()

	validator, err := NewRemoteValidator(context.Background(), RemoteValidatorConfig{
		JWKSURL:   jwksServer.URL,
		Issuer:    "http://auth.codelinks.localhost",
		Audience:  "codelinks-products",
		ProductID: "infra-link",
	})
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	request.Header.Set("Authorization", "Bearer "+signToken(t, privateKey, "test-key", []string{"infra-link"}))
	recorder := httptest.NewRecorder()

	validator.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := CurrentUserFromContext(r.Context())
		if !ok {
			t.Fatal("expected current user in context")
		}
		if user.ID != "user-1" || user.Email != "demo@codelinks.dev" || !user.EmailVerified || !user.HasRole("user") {
			t.Fatalf("unexpected current user: %#v", user)
		}
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", recorder.Code)
	}
}

func newTestKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	return key
}

func testJWK(kid string, key *rsa.PublicKey) jwk {
	return jwk{
		Kty: "RSA",
		Use: "sig",
		Kid: kid,
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
	}
}

func signToken(t *testing.T, key *rsa.PrivateKey, kid string, licenses []string) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims{
		Email:         "demo@codelinks.dev",
		Name:          "Demo User",
		Status:        "active",
		EmailVerified: true,
		Licenses:      licenses,
		Roles:         []string{"user"},
		Permissions:   []string{"product.read"},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-1",
			Issuer:    "http://auth.codelinks.localhost",
			Audience:  []string{"codelinks-products"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	token.Header["kid"] = kid

	raw, err := token.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}

	return raw
}
