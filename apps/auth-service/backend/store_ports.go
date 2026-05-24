package main

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthStore interface {
	FindUserByEmail(ctx context.Context, email string) (User, []string, error)
	FindUserByID(ctx context.Context, userID uuid.UUID) (User, []string, error)
	CreateUser(ctx context.Context, name string, email string, passwordHash string) (User, error)
	GrantLicense(ctx context.Context, userID uuid.UUID, productID string) ([]string, error)
	CreateRefreshSession(ctx context.Context, tokenHash string, userID uuid.UUID, expiresAt time.Time) error
	FindRefreshSession(ctx context.Context, tokenHash string) (uuid.UUID, error)
	DeleteRefreshSession(ctx context.Context, tokenHash string) error
	CreateEmailVerificationToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	VerifyEmailToken(ctx context.Context, tokenHash string, now time.Time) (User, error)
	CreatePasswordResetToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	ResetPasswordByToken(ctx context.Context, tokenHash string, passwordHash string, now time.Time) (User, error)
	RecordLoginAttempt(ctx context.Context, attempt LoginAttempt) error
}

type AdminStore interface {
	GetAdminActor(ctx context.Context, userID uuid.UUID) (AdminActor, error)
	ListUserRoles(ctx context.Context, userID uuid.UUID) ([]AdminRole, []AdminPermission, error)
	ListAdminUsers(ctx context.Context, query AdminUserListQuery) (AdminUserListResult, error)
	GetManagedUserDetail(ctx context.Context, userID uuid.UUID) (ManagedUserDetail, error)
	SetUserStatus(ctx context.Context, userID uuid.UUID, status UserStatus) (User, error)
	SetUserRole(ctx context.Context, userID uuid.UUID, role AdminRole) (User, error)
	ListLoginAttempts(ctx context.Context, query LoginAttemptListQuery) (LoginAttemptListResult, error)
	ListSecurityEvents(ctx context.Context, limit int) ([]SecurityEvent, error)
	GetDashboardStats(ctx context.Context) (DashboardStats, error)
}

type SmtpSettingsStore interface {
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

type Store interface {
	AuthStore
	AdminStore
	SmtpSettingsStore
	NotificationStore
	AuditStore
}
