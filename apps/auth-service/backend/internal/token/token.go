package token

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

	"github.com/besart951/code-links/apps/auth-service/backend/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Config struct {
	KeyID          string
	Issuer         string
	Audience       string
	Lifetime       time.Duration
	PrivateKeyPEM  string
	PrivateKeyFile string
}

type Signer struct {
	privateKey *rsa.PrivateKey
	keyID      string
	issuer     string
	audience   string
	lifetime   time.Duration
}

type Claims struct {
	Email         string                   `json:"email"`
	Name          string                   `json:"name"`
	Status        domain.UserStatus        `json:"status"`
	EmailVerified bool                     `json:"emailVerified"`
	Licenses      []string                 `json:"licenses"`
	Roles         []domain.AdminRole       `json:"roles"`
	Permissions   []domain.AdminPermission `json:"permissions"`
	jwt.RegisteredClaims
}

func NewSigner(config Config) (*Signer, error) {
	privateKey, err := loadPrivateKey(config)
	if err != nil {
		return nil, err
	}

	return &Signer{
		privateKey: privateKey,
		keyID:      config.KeyID,
		issuer:     config.Issuer,
		audience:   config.Audience,
		lifetime:   config.Lifetime,
	}, nil
}

func loadPrivateKey(config Config) (*rsa.PrivateKey, error) {
	if config.PrivateKeyPEM != "" {
		return parsePrivateKey([]byte(config.PrivateKeyPEM))
	}
	if config.PrivateKeyFile != "" {
		content, err := os.ReadFile(config.PrivateKeyFile)
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

func (s *Signer) Issue(user domain.User, licenses []string, roles []domain.AdminRole, permissions []domain.AdminPermission) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(s.lifetime)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims{
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

func (s *Signer) Parse(rawToken string) (Claims, error) {
	claims := Claims{}
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
		return Claims{}, err
	}
	if !token.Valid {
		return Claims{}, errors.New("invalid token")
	}
	if _, err := uuid.Parse(claims.Subject); err != nil {
		return Claims{}, errors.New("invalid subject")
	}

	return claims, nil
}

type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func (s *Signer) JWKS() JWKSResponse {
	publicKey := &s.privateKey.PublicKey

	return JWKSResponse{Keys: []JWK{{
		Kty: "RSA",
		Use: "sig",
		Kid: s.keyID,
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
	}}}
}
