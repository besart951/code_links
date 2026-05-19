package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInactiveUser       = errors.New("inactive user")
	ErrTokenRevoked       = errors.New("token revoked")
)

type Config struct {
	JWTSecret       string
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type Service struct {
	repo Repository
	cfg  Config
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewService(repo Repository, cfg Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

func (s *Service) Login(ctx context.Context, email, password string, userAgent, ip *string) (TokenPair, error) {
	user, err := s.repo.FindUserByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return TokenPair{}, ErrInvalidCredentials
	}
	if user.Status != "active" {
		return TokenPair{}, ErrInactiveUser
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return TokenPair{}, ErrInvalidCredentials
	}

	now := time.Now().UTC()
	if err := s.repo.TouchLastLogin(ctx, user.ID, now); err != nil {
		return TokenPair{}, err
	}
	user.LastLoginAt = &now

	return s.issueTokens(ctx, user, userAgent, ip)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string, userAgent, ip *string) (TokenPair, error) {
	hash := HashToken(refreshToken)
	record, err := s.repo.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return TokenPair{}, ErrUnauthorized
	}
	now := time.Now().UTC()
	if record.RevokedAt != nil {
		return TokenPair{}, ErrTokenRevoked
	}
	if !record.ExpiresAt.After(now) {
		return TokenPair{}, ErrUnauthorized
	}

	user, err := s.repo.FindUserByID(ctx, record.UserID)
	if err != nil {
		return TokenPair{}, ErrUnauthorized
	}
	if user.Status != "active" {
		return TokenPair{}, ErrInactiveUser
	}

	if err := s.repo.RevokeRefreshTokenByHash(ctx, hash, now); err != nil {
		return TokenPair{}, err
	}
	return s.issueTokens(ctx, user, userAgent, ip)
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return nil
	}
	return s.repo.RevokeRefreshTokenByHash(ctx, HashToken(refreshToken), time.Now().UTC())
}

func (s *Service) ValidateAccessToken(raw string) (Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(raw, &claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrUnauthorized
		}
		return []byte(s.cfg.JWTSecret), nil
	}, jwt.WithIssuer(s.cfg.Issuer))
	if err != nil || !token.Valid {
		return Claims{}, ErrUnauthorized
	}
	if claims.Subject == "" {
		return Claims{}, ErrUnauthorized
	}
	return claims, nil
}

func (s *Service) issueTokens(ctx context.Context, user User, userAgent, ip *string) (TokenPair, error) {
	now := time.Now().UTC()
	accessExpiresAt := now.Add(s.cfg.AccessTokenTTL)
	refreshExpiresAt := now.Add(s.cfg.RefreshTokenTTL)
	refreshToken, err := randomToken(48)
	if err != nil {
		return TokenPair{}, err
	}
	csrfToken, err := randomToken(32)
	if err != nil {
		return TokenPair{}, err
	}

	claims := Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			Issuer:    s.cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return TokenPair{}, err
	}

	recordID, err := randomToken(16)
	if err != nil {
		return TokenPair{}, err
	}
	record := RefreshToken{
		ID:        recordID,
		UserID:    user.ID,
		TokenHash: HashToken(refreshToken),
		UserAgent: userAgent,
		IP:        ip,
		ExpiresAt: refreshExpiresAt,
		CreatedAt: now,
	}
	if err := s.repo.StoreRefreshToken(ctx, record); err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		CSRFToken:        csrfToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		User:             user,
	}, nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func randomToken(bytes int) (string, error) {
	buf := make([]byte, bytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
