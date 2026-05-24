package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/mail"
	"time"

	"github.com/besart951/code-links/apps/auth-service/backend/internal/domain"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/secret"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Store interface {
	FindUserByEmail(ctx context.Context, email string) (domain.User, []string, error)
	FindUserByID(ctx context.Context, userID uuid.UUID) (domain.User, []string, error)
	CreateUser(ctx context.Context, name string, email string, passwordHash string) (domain.User, error)
	GrantLicense(ctx context.Context, userID uuid.UUID, productID string) ([]string, error)
	CreateRefreshSession(ctx context.Context, tokenHash string, userID uuid.UUID, expiresAt time.Time) error
	FindRefreshSession(ctx context.Context, tokenHash string) (uuid.UUID, error)
	DeleteRefreshSession(ctx context.Context, tokenHash string) error
	CreateEmailVerificationToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	VerifyEmailToken(ctx context.Context, tokenHash string, now time.Time) (domain.User, error)
	CreatePasswordResetToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	ResetPasswordByToken(ctx context.Context, tokenHash string, passwordHash string, now time.Time) (domain.User, error)
	RecordLoginAttempt(ctx context.Context, attempt domain.LoginAttempt) error
}

type RoleStore interface {
	ListUserRoles(ctx context.Context, userID uuid.UUID) ([]domain.AdminRole, []domain.AdminPermission, error)
}

type NotificationStore interface {
	CreateNotification(ctx context.Context, notification domain.Notification) error
}

type TokenIssuer interface {
	Issue(user domain.User, licenses []string, roles []domain.AdminRole, permissions []domain.AdminPermission) (string, time.Time, error)
}

type Config struct {
	Environment          string
	PublicFrontendURL    string
	RefreshTokenLifetime time.Duration
}

type Service struct {
	config        Config
	store         Store
	roleStore     RoleStore
	notifications NotificationStore
	tokens        TokenIssuer
}

func NewService(config Config, store Store, roleStore RoleStore, notifications NotificationStore, tokens TokenIssuer) *Service {
	return &Service{config: config, store: store, roleStore: roleStore, notifications: notifications, tokens: tokens}
}

type Kind string

const (
	KindBadRequest    Kind = "bad_request"
	KindUnauthorized  Kind = "unauthorized"
	KindForbidden     Kind = "forbidden"
	KindConflict      Kind = "conflict"
	KindInternal      Kind = "internal"
	KindBadGateway    Kind = "bad_gateway"
	KindInvalidToken  Kind = "invalid_token"
	KindMissingToken  Kind = "missing_token"
	KindEmailConflict Kind = "email_conflict"
)

type Error struct {
	Kind    Kind
	Message string
}

func (e *Error) Error() string { return e.Message }

func serviceError(kind Kind, message string) *Error {
	return &Error{Kind: kind, Message: message}
}

type SignupInput struct {
	Name                    string
	Email                   string
	Password                string
	AcceptedTerms           bool
	AcceptedTermsAndPrivacy bool
}

type SignupResult struct {
	Message                string
	VerificationURL        string
	DebugVerificationToken string
}

func (s *Service) Signup(ctx context.Context, input SignupInput) (SignupResult, error) {
	if _, err := mail.ParseAddress(input.Email); err != nil {
		return SignupResult{}, serviceError(KindBadRequest, "invalid email")
	}
	if len(input.Password) < 10 {
		return SignupResult{}, serviceError(KindBadRequest, "password must be at least 10 characters")
	}
	if !input.AcceptedTerms && !input.AcceptedTermsAndPrivacy {
		return SignupResult{}, serviceError(KindBadRequest, "terms must be accepted")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return SignupResult{}, serviceError(KindInternal, "could not hash password")
	}

	user, err := s.store.CreateUser(ctx, input.Name, input.Email, string(passwordHash))
	if errors.Is(err, domain.ErrConflict) {
		return SignupResult{}, serviceError(KindConflict, "email already exists")
	}
	if err != nil {
		return SignupResult{}, serviceError(KindInternal, "could not create user")
	}

	token, err := secret.NewOpaqueToken()
	if err != nil {
		return SignupResult{}, serviceError(KindInternal, "could not create verification token")
	}
	if err := s.store.CreateEmailVerificationToken(ctx, user.ID, secret.HashOpaqueToken(token), time.Now().UTC().Add(24*time.Hour)); err != nil {
		return SignupResult{}, serviceError(KindInternal, "could not create verification token")
	}

	_ = s.notifications.CreateNotification(ctx, domain.Notification{
		ID:        uuid.New(),
		Type:      "email_verification",
		Channel:   "email",
		Recipient: user.Email,
		Subject:   "CodeLinks E-Mail bestätigen",
		Status:    "queued",
		CreatedAt: time.Now().UTC(),
	})

	result := SignupResult{
		Message:         "Account created. Please verify your email before logging in.",
		VerificationURL: s.config.PublicFrontendURL + "/verify-email?token=" + token,
	}
	if s.config.Environment != "production" {
		result.DebugVerificationToken = token
	}
	return result, nil
}

func (s *Service) VerifyEmail(ctx context.Context, token string) (domain.User, error) {
	if token == "" {
		return domain.User{}, serviceError(KindBadRequest, "missing token")
	}
	user, err := s.store.VerifyEmailToken(ctx, secret.HashOpaqueToken(token), time.Now().UTC())
	if err != nil {
		return domain.User{}, serviceError(KindBadRequest, "invalid or expired token")
	}
	return user, nil
}

type ForgotPasswordResult struct {
	Message         string
	DebugResetToken string
	DebugResetURL   string
}

func (s *Service) ForgotPassword(ctx context.Context, email string) (ForgotPasswordResult, error) {
	result := ForgotPasswordResult{Message: "Falls ein Konto mit dieser E-Mail existiert, wurde ein Reset-Link versendet."}
	user, _, err := s.store.FindUserByEmail(ctx, email)
	if err != nil {
		return result, nil
	}

	token, err := secret.NewOpaqueToken()
	if err != nil {
		return ForgotPasswordResult{}, serviceError(KindInternal, "could not create reset token")
	}
	if err := s.store.CreatePasswordResetToken(ctx, user.ID, secret.HashOpaqueToken(token), time.Now().UTC().Add(time.Hour)); err != nil {
		return ForgotPasswordResult{}, serviceError(KindInternal, "could not create reset token")
	}
	_ = s.notifications.CreateNotification(ctx, domain.Notification{
		ID:        uuid.New(),
		Type:      "password_reset",
		Channel:   "email",
		Recipient: user.Email,
		Subject:   "CodeLinks Passwort zurücksetzen",
		Status:    "queued",
		CreatedAt: time.Now().UTC(),
	})
	if s.config.Environment != "production" {
		result.DebugResetToken = token
		result.DebugResetURL = s.config.PublicFrontendURL + "/reset-password?token=" + token
	}
	return result, nil
}

func (s *Service) ResetPassword(ctx context.Context, token string, password string) error {
	if token == "" {
		return serviceError(KindBadRequest, "missing token")
	}
	if len(password) < 10 {
		return serviceError(KindBadRequest, "password must be at least 10 characters")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return serviceError(KindInternal, "could not hash password")
	}
	if _, err := s.store.ResetPasswordByToken(ctx, secret.HashOpaqueToken(token), string(passwordHash), time.Now().UTC()); err != nil {
		return serviceError(KindBadRequest, "invalid or expired token")
	}
	return nil
}

type LoginInput struct {
	Email    string
	Password string
	Attempt  LoginAttemptMetadata
}

type LoginAttemptMetadata struct {
	IPAddress       string
	CountryCode     string
	City            string
	UserAgent       string
	Browser         string
	OperatingSystem string
	CorrelationID   string
}

type Session struct {
	AccessToken      string
	AccessExpiresAt  time.Time
	RefreshToken     string
	RefreshExpiresAt time.Time
	User             domain.User
	Licenses         []string
	Roles            []domain.AdminRole
	Permissions      []domain.AdminPermission
}

func (s *Service) Login(ctx context.Context, input LoginInput) (Session, error) {
	email := domain.NormalizeEmail(input.Email)
	user, licenses, err := s.store.FindUserByEmail(ctx, email)
	if errors.Is(err, domain.ErrNotFound) {
		s.recordLoginAttempt(ctx, nil, email, false, domain.LoginFailureUnknownEmail, input.Attempt)
		return Session{}, serviceError(KindUnauthorized, "invalid email or password")
	}
	if err != nil {
		return Session{}, serviceError(KindInternal, "could not load user")
	}

	if user.Status == domain.UserStatusLocked || user.Status == domain.UserStatusDisabled {
		s.recordLoginAttempt(ctx, &user.ID, email, false, domain.LoginFailureAccountLocked, input.Attempt)
		return Session{}, serviceError(KindForbidden, "account is locked")
	}
	if user.EmailVerifiedAt == nil {
		s.recordLoginAttempt(ctx, &user.ID, email, false, domain.LoginFailureEmailNotConfirmed, input.Attempt)
		return Session{}, serviceError(KindForbidden, "email_not_confirmed")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
		s.recordLoginAttempt(ctx, &user.ID, email, false, domain.LoginFailureWrongPassword, input.Attempt)
		return Session{}, serviceError(KindUnauthorized, "invalid email or password")
	}

	s.recordLoginAttempt(ctx, &user.ID, email, true, "", input.Attempt)
	return s.issueSession(ctx, user, licenses, true)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (Session, error) {
	if refreshToken == "" {
		return Session{}, serviceError(KindUnauthorized, "missing refresh token")
	}
	tokenHash := secret.HashRefreshToken(refreshToken)
	userID, err := s.store.FindRefreshSession(ctx, tokenHash)
	if err != nil {
		return Session{}, serviceError(KindUnauthorized, "invalid refresh token")
	}
	_ = s.store.DeleteRefreshSession(ctx, tokenHash)

	user, licenses, err := s.store.FindUserByID(ctx, userID)
	if err != nil {
		return Session{}, serviceError(KindUnauthorized, "invalid refresh token")
	}
	return s.issueSession(ctx, user, licenses, true)
}

func (s *Service) Logout(ctx context.Context, refreshToken string) {
	if refreshToken != "" {
		_ = s.store.DeleteRefreshSession(ctx, secret.HashRefreshToken(refreshToken))
	}
}

func (s *Service) MockPurchase(ctx context.Context, userID uuid.UUID, productID string) (Session, error) {
	licenses, err := s.store.GrantLicense(ctx, userID, productID)
	if err != nil {
		return Session{}, serviceError(KindInternal, "could not grant license")
	}
	user, _, err := s.store.FindUserByID(ctx, userID)
	if err != nil {
		return Session{}, serviceError(KindInternal, "could not load user")
	}
	return s.issueSession(ctx, user, licenses, false)
}

func (s *Service) issueSession(ctx context.Context, user domain.User, licenses []string, rotateRefresh bool) (Session, error) {
	session := Session{User: user, Licenses: licenses}
	if rotateRefresh {
		refreshToken, expiresAt, err := s.createRefreshSession(ctx, user.ID)
		if err != nil {
			return Session{}, serviceError(KindInternal, "could not create refresh session")
		}
		session.RefreshToken = refreshToken
		session.RefreshExpiresAt = expiresAt
	}

	roles, permissions, err := s.roleStore.ListUserRoles(ctx, user.ID)
	if err != nil {
		return Session{}, serviceError(KindInternal, "could not load roles")
	}
	accessToken, accessExpiresAt, err := s.tokens.Issue(user, licenses, roles, permissions)
	if err != nil {
		return Session{}, serviceError(KindInternal, "could not issue token")
	}
	session.Roles = roles
	session.Permissions = permissions
	session.AccessToken = accessToken
	session.AccessExpiresAt = accessExpiresAt
	return session, nil
}

func (s *Service) createRefreshSession(ctx context.Context, userID uuid.UUID) (string, time.Time, error) {
	token, err := secret.NewOpaqueToken()
	if err != nil {
		return "", time.Time{}, err
	}
	expiresAt := time.Now().UTC().Add(s.config.RefreshTokenLifetime)
	if err := s.store.CreateRefreshSession(ctx, secret.HashRefreshToken(token), userID, expiresAt); err != nil {
		return "", time.Time{}, err
	}
	return token, expiresAt, nil
}

func (s *Service) recordLoginAttempt(ctx context.Context, userID *uuid.UUID, email string, success bool, failureReason domain.LoginFailureReason, metadata LoginAttemptMetadata) {
	var reason *domain.LoginFailureReason
	if failureReason != "" {
		reason = &failureReason
	}
	ipHash := sha256.Sum256([]byte(metadata.IPAddress))
	_ = s.store.RecordLoginAttempt(ctx, domain.LoginAttempt{
		ID:              uuid.New(),
		UserID:          userID,
		EmailAttempted:  email,
		OccurredAt:      time.Now().UTC(),
		IPAddress:       metadata.IPAddress,
		IPHash:          base64.RawURLEncoding.EncodeToString(ipHash[:]),
		CountryCode:     metadata.CountryCode,
		City:            metadata.City,
		UserAgent:       metadata.UserAgent,
		Browser:         metadata.Browser,
		OperatingSystem: metadata.OperatingSystem,
		Success:         success,
		FailureReason:   reason,
		AuthMethod:      "password",
		RiskScore:       riskScore(success, failureReason),
		CorrelationID:   metadata.CorrelationID,
	})
}

func riskScore(success bool, failureReason domain.LoginFailureReason) int {
	if success {
		return 8
	}
	switch failureReason {
	case domain.LoginFailureAccountLocked:
		return 85
	case domain.LoginFailureUnknownEmail:
		return 55
	case domain.LoginFailureEmailNotConfirmed:
		return 30
	default:
		return 65
	}
}
