package admin

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	appmail "github.com/besart951/code-links/apps/auth-service/backend/internal/mail"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/secret"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/token"
	"github.com/google/uuid"
)

type Store interface {
	GetAdminActor(ctx context.Context, userID uuid.UUID) (AdminActor, error)
	ListAdminUsers(ctx context.Context, query AdminUserListQuery) (AdminUserListResult, error)
	GetManagedUserDetail(ctx context.Context, userID uuid.UUID) (ManagedUserDetail, error)
	SetUserStatus(ctx context.Context, userID uuid.UUID, status UserStatus) (User, error)
	SetUserRole(ctx context.Context, userID uuid.UUID, role AdminRole) (User, error)
	ListLoginAttempts(ctx context.Context, query LoginAttemptListQuery) (LoginAttemptListResult, error)
	ListSecurityEvents(ctx context.Context, limit int) ([]SecurityEvent, error)
	GetDashboardStats(ctx context.Context) (DashboardStats, error)
}

type SettingsStore interface {
	GetSmtpSettings(ctx context.Context) (SmtpSettings, error)
	SaveSmtpSettings(ctx context.Context, settings SmtpSettings) (SmtpSettings, error)
}

type NotificationStore interface {
	ListNotifications(ctx context.Context, limit int) ([]Notification, error)
	CreateNotification(ctx context.Context, notification Notification) error
}

type AuditStore interface {
	RecordAdminAuditEntry(ctx context.Context, entry AdminAuditEntry) error
}

type SessionStore interface {
	FindRefreshSession(ctx context.Context, tokenHash string) (uuid.UUID, error)
}

type TokenParser interface {
	Parse(rawToken string) (token.Claims, error)
}

type Config struct {
	SMTPSecretKey []byte
}

type Service struct {
	config        Config
	store         Store
	settings      SettingsStore
	notifications NotificationStore
	audit         AuditStore
	sessions      SessionStore
	tokens        TokenParser
	emailSender   appmail.Sender
	projection    AdminProjectionPolicy
}

func NewService(config Config, store Store, settings SettingsStore, notifications NotificationStore, audit AuditStore, sessions SessionStore, tokens TokenParser, emailSender appmail.Sender) *Service {
	return &Service{
		config:        config,
		store:         store,
		settings:      settings,
		notifications: notifications,
		audit:         audit,
		sessions:      sessions,
		tokens:        tokens,
		emailSender:   emailSender,
		projection:    AdminProjectionPolicy{},
	}
}

type Kind string

const (
	KindBadRequest   Kind = "bad_request"
	KindUnauthorized Kind = "unauthorized"
	KindForbidden    Kind = "forbidden"
	KindNotFound     Kind = "not_found"
	KindInternal     Kind = "internal"
	KindBadGateway   Kind = "bad_gateway"
)

type Error struct {
	Kind    Kind
	Message string
}

func (e *Error) Error() string { return e.Message }

func serviceError(kind Kind, message string) *Error {
	return &Error{Kind: kind, Message: message}
}

type Authn struct {
	BearerToken  string
	RefreshToken string
}

type RequestMeta struct {
	IPAddress string
}

type SMTPSettingsInput struct {
	Host         string
	Port         int
	Username     string
	Password     string
	Encryption   SmtpEncryption
	FromEmail    string
	FromName     string
	ReplyToEmail string
	Active       bool
}

func (s *Service) ResolveActor(ctx context.Context, authn Authn, permission AdminPermission) (AdminActor, error) {
	actor, err := s.actorFromAuthn(ctx, authn)
	if err != nil {
		return AdminActor{}, serviceError(KindUnauthorized, "admin authentication required")
	}
	if !HasPermission(actor, permission) {
		return AdminActor{}, serviceError(KindForbidden, "missing permission")
	}
	return actor, nil
}

func (s *Service) actorFromAuthn(ctx context.Context, authn Authn) (AdminActor, error) {
	if authn.BearerToken != "" {
		claims, err := s.tokens.Parse(authn.BearerToken)
		if err != nil {
			return AdminActor{}, err
		}
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return AdminActor{}, err
		}
		return s.store.GetAdminActor(ctx, userID)
	}

	if authn.RefreshToken == "" {
		return AdminActor{}, serviceError(KindUnauthorized, "admin authentication required")
	}
	userID, err := s.sessions.FindRefreshSession(ctx, secret.HashRefreshToken(authn.RefreshToken))
	if err != nil {
		return AdminActor{}, err
	}
	return s.store.GetAdminActor(ctx, userID)
}

func (s *Service) DashboardStats(ctx context.Context, actor AdminActor) (DashboardStats, error) {
	stats, err := s.store.GetDashboardStats(ctx)
	if err != nil {
		return DashboardStats{}, serviceError(KindInternal, "could not load dashboard stats")
	}
	return s.projection.ProjectDashboardStats(actor, stats), nil
}

func (s *Service) ListUsers(ctx context.Context, actor AdminActor, query AdminUserListQuery) (AdminUserListResult, error) {
	result, err := s.store.ListAdminUsers(ctx, query)
	if err != nil {
		return AdminUserListResult{}, serviceError(KindInternal, "could not load users")
	}
	return s.projection.ProjectUserList(actor, result), nil
}

func (s *Service) UserDetail(ctx context.Context, actor AdminActor, userID uuid.UUID) (ManagedUserDetail, error) {
	user, err := s.store.GetManagedUserDetail(ctx, userID)
	if errors.Is(err, errNotFound) {
		return ManagedUserDetail{}, serviceError(KindNotFound, "user not found")
	}
	if err != nil {
		return ManagedUserDetail{}, serviceError(KindInternal, "could not load user")
	}
	return s.projection.ProjectManagedUserDetail(actor, user), nil
}

func (s *Service) SetUserStatus(ctx context.Context, actor AdminActor, userID uuid.UUID, status UserStatus, meta RequestMeta) (User, error) {
	user, err := s.store.SetUserStatus(ctx, userID, status)
	if errors.Is(err, errNotFound) {
		return User{}, serviceError(KindNotFound, "user not found")
	}
	if err != nil {
		return User{}, serviceError(KindInternal, "could not update user")
	}
	s.recordAudit(ctx, actor, "admin.users.update_status", "user", userID.String(), string(status), meta)
	return user, nil
}

func (s *Service) LockUser(ctx context.Context, actor AdminActor, userID uuid.UUID, meta RequestMeta) (User, error) {
	user, err := s.store.SetUserStatus(ctx, userID, UserStatusLocked)
	if errors.Is(err, errNotFound) {
		return User{}, serviceError(KindNotFound, "user not found")
	}
	if err != nil {
		return User{}, serviceError(KindInternal, "could not lock user")
	}
	s.recordAudit(ctx, actor, "admin.users.lock", "user", userID.String(), "", meta)
	return user, nil
}

func (s *Service) UnlockUser(ctx context.Context, actor AdminActor, userID uuid.UUID, meta RequestMeta) (User, error) {
	user, err := s.store.SetUserStatus(ctx, userID, UserStatusActive)
	if errors.Is(err, errNotFound) {
		return User{}, serviceError(KindNotFound, "user not found")
	}
	if err != nil {
		return User{}, serviceError(KindInternal, "could not unlock user")
	}
	s.recordAudit(ctx, actor, "admin.users.unlock", "user", userID.String(), "", meta)
	return user, nil
}

func (s *Service) SetUserRole(ctx context.Context, actor AdminActor, userID uuid.UUID, role AdminRole, meta RequestMeta) (User, error) {
	user, err := s.store.SetUserRole(ctx, userID, role)
	if errors.Is(err, errNotFound) {
		return User{}, serviceError(KindNotFound, "user not found")
	}
	if err != nil {
		return User{}, serviceError(KindInternal, "could not change role")
	}
	s.recordAudit(ctx, actor, "admin.users.change_role", "user", userID.String(), string(role), meta)
	return user, nil
}

func (s *Service) ListLoginAttempts(ctx context.Context, actor AdminActor, query LoginAttemptListQuery) (LoginAttemptListResult, error) {
	result, err := s.store.ListLoginAttempts(ctx, query)
	if err != nil {
		return LoginAttemptListResult{}, serviceError(KindInternal, "could not load login attempts")
	}
	return s.projection.ProjectLoginAttempts(actor, result), nil
}

func (s *Service) ListSecurityEvents(ctx context.Context, actor AdminActor, limit int) ([]SecurityEvent, error) {
	events, err := s.store.ListSecurityEvents(ctx, limit)
	if err != nil {
		return nil, serviceError(KindInternal, "could not load security events")
	}
	return s.projection.ProjectSecurityEvents(actor, events), nil
}

func (s *Service) ListNotifications(ctx context.Context, limit int) ([]Notification, error) {
	notifications, err := s.notifications.ListNotifications(ctx, limit)
	if err != nil {
		return nil, serviceError(KindInternal, "could not load notifications")
	}
	return notifications, nil
}

func (s *Service) SMTPSettings(ctx context.Context) (SmtpSettings, error) {
	settings, err := s.settings.GetSmtpSettings(ctx)
	if err != nil {
		return SmtpSettings{}, serviceError(KindInternal, "could not load smtp settings")
	}
	return settings, nil
}

func (s *Service) UpdateSMTPSettings(ctx context.Context, actor AdminActor, input SMTPSettingsInput, meta RequestMeta) (SmtpSettings, error) {
	encryptedPassword := ""
	if input.Password != "" {
		var err error
		encryptedPassword, err = secret.Encrypt(input.Password, s.config.SMTPSecretKey)
		if err != nil {
			return SmtpSettings{}, serviceError(KindInternal, "could not encrypt smtp password")
		}
	}

	settings, err := s.settings.SaveSmtpSettings(ctx, SmtpSettings{
		Host:              strings.TrimSpace(input.Host),
		Port:              input.Port,
		Username:          strings.TrimSpace(input.Username),
		PasswordEncrypted: encryptedPassword,
		Encryption:        input.Encryption,
		FromEmail:         strings.TrimSpace(input.FromEmail),
		FromName:          strings.TrimSpace(input.FromName),
		ReplyToEmail:      strings.TrimSpace(input.ReplyToEmail),
		Active:            input.Active,
	})
	if err != nil {
		return SmtpSettings{}, serviceError(KindInternal, "could not save smtp settings")
	}
	s.recordAudit(ctx, actor, "admin.smtp_settings.update", "smtp_settings", "default", "", meta)
	return settings, nil
}

func (s *Service) SendTestEmail(ctx context.Context, actor AdminActor, recipient string, meta RequestMeta) error {
	if _, err := mail.ParseAddress(recipient); err != nil {
		return serviceError(KindBadRequest, "invalid recipient")
	}
	settings, err := s.settings.GetSmtpSettings(ctx)
	if err != nil {
		return serviceError(KindInternal, "could not load smtp settings")
	}
	password := ""
	if settings.PasswordEncrypted != "" {
		password, err = secret.Decrypt(settings.PasswordEncrypted, s.config.SMTPSecretKey)
		if err != nil {
			return serviceError(KindInternal, "could not decrypt smtp password")
		}
	}
	if err := s.emailSender.Send(ctx, settings, appmail.Message{
		To:      recipient,
		Subject: "CodeLinks SMTP Test",
		Body:    "Diese Test-E-Mail wurde von CodeLinks gesendet.",
	}, password); err != nil {
		return serviceError(KindBadGateway, err.Error())
	}

	now := time.Now().UTC()
	_ = s.notifications.CreateNotification(ctx, Notification{
		ID:        uuid.New(),
		Type:      "smtp_test",
		Channel:   "email",
		Recipient: recipient,
		Subject:   "CodeLinks SMTP Test",
		Status:    "sent",
		CreatedAt: now,
		SentAt:    &now,
	})
	s.recordAudit(ctx, actor, "admin.smtp_settings.test_email", "smtp_settings", "default", recipient, meta)
	return nil
}

func (s *Service) recordAudit(ctx context.Context, actor AdminActor, action string, targetType string, targetID string, reason string, meta RequestMeta) {
	actorID, err := uuid.Parse(actor.ID)
	if err != nil {
		return
	}
	_ = s.audit.RecordAdminAuditEntry(ctx, AdminAuditEntry{
		ID:          uuid.New(),
		ActorUserID: actorID,
		Action:      action,
		TargetType:  targetType,
		TargetID:    targetID,
		Reason:      reason,
		IPAddress:   meta.IPAddress,
		CreatedAt:   time.Now().UTC(),
	})
}
