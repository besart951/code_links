package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type tokenSigner struct {
	privateKey *rsa.PrivateKey
	keyID      string
	issuer     string
	audience   string
	lifetime   time.Duration
}

type accessClaims struct {
	Email         string            `json:"email"`
	Name          string            `json:"name"`
	Status        UserStatus        `json:"status"`
	EmailVerified bool              `json:"emailVerified"`
	Licenses      []string          `json:"licenses"`
	Roles         []AdminRole       `json:"roles"`
	Permissions   []AdminPermission `json:"permissions"`
	jwt.RegisteredClaims
}

func newTokenSigner(config config) (*tokenSigner, error) {
	privateKey, err := loadPrivateKey(config)
	if err != nil {
		return nil, err
	}

	return &tokenSigner{
		privateKey: privateKey,
		keyID:      config.JWTKeyID,
		issuer:     config.Issuer,
		audience:   config.Audience,
		lifetime:   config.AccessTokenLifetime,
	}, nil
}

func loadPrivateKey(config config) (*rsa.PrivateKey, error) {
	if config.JWTPrivateKeyPEM != "" {
		return parsePrivateKey([]byte(config.JWTPrivateKeyPEM))
	}
	if config.JWTPrivateKeyFile != "" {
		content, err := os.ReadFile(config.JWTPrivateKeyFile)
		if err != nil {
			return nil, err
		}
		return parsePrivateKey(content)
	}

	return rsa.GenerateKey(rand.Reader, 2048)
}

func parsePrivateKey(content []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(content)
	if block == nil {
		return nil, errors.New("JWT private key must be PEM encoded")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("JWT private key must be RSA")
	}

	return key, nil
}

func (s *tokenSigner) issue(user User, licenses []string, roles []AdminRole, permissions []AdminPermission) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(s.lifetime)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims{
		Email:         user.Email,
		Name:          user.Name,
		Status:        user.Status,
		EmailVerified: user.EmailVerifiedAt != nil,
		Licenses:      licenses,
		Roles:         roles,
		Permissions:   permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			Issuer:    s.issuer,
			Audience:  []string{s.audience},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})
	token.Header["kid"] = s.keyID

	rawToken, err := token.SignedString(s.privateKey)
	return rawToken, expiresAt, err
}

func (s *tokenSigner) parse(rawToken string) (accessClaims, error) {
	claims := accessClaims{}
	parser := jwt.NewParser(
		jwt.WithIssuer(s.issuer),
		jwt.WithAudience(s.audience),
		jwt.WithExpirationRequired(),
	)

	token, err := parser.ParseWithClaims(rawToken, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return &s.privateKey.PublicKey, nil
	})
	if err != nil {
		return accessClaims{}, err
	}
	if !token.Valid {
		return accessClaims{}, errors.New("invalid token")
	}
	if _, err := uuid.Parse(claims.Subject); err != nil {
		return accessClaims{}, errors.New("invalid subject")
	}

	return claims, nil
}

type jwksResponse struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func (s *tokenSigner) jwks() jwksResponse {
	publicKey := &s.privateKey.PublicKey

	return jwksResponse{Keys: []jwk{{
		Kty: "RSA",
		Use: "sig",
		Kid: s.keyID,
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
	}}}
}
