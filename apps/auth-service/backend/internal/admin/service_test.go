package admin

import (
	"context"
	"errors"
	"testing"
	"time"

	authsvc "github.com/besart951/code-links/apps/auth-service/backend/internal/auth"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/domain"
	appmail "github.com/besart951/code-links/apps/auth-service/backend/internal/mail"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/store/memory"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/token"
)

func TestServicePermissionAndSMTPFlow(t *testing.T) {
	ctx := context.Background()
	store, err := memory.New()
	if err != nil {
		t.Fatal(err)
	}
	signer, err := token.NewSigner(token.Config{
		KeyID:    "test-key",
		Issuer:   "http://auth.codelinks.localhost",
		Audience: "codelinks-products",
		Lifetime: 15 * time.Minute,
	})
	if err != nil {
		t.Fatal(err)
	}
	authService := authsvc.NewService(authsvc.Config{
		Environment:          "test",
		PublicFrontendURL:    "http://auth.codelinks.localhost",
		RefreshTokenLifetime: time.Hour,
	}, store, store, store, signer)
	email := &spySender{}
	adminService := NewService(Config{
		SMTPSecretKey: []byte("01234567890123456789012345678901"),
	}, store, store, store, store, staticRuntimeLogs{}, store, signer, email)

	support, err := authService.Login(ctx, authsvc.LoginInput{Email: "support@codelinks.dev", Password: "password"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := adminService.ResolveActor(ctx, Authn{BearerToken: support.AccessToken}, domain.PermissionSMTPSettingsUpdate); !adminKind(err, KindForbidden) {
		t.Fatalf("expected support forbidden, got %v", err)
	}

	adminSession, err := authService.Login(ctx, authsvc.LoginInput{Email: "demo@codelinks.dev", Password: "password"})
	if err != nil {
		t.Fatal(err)
	}
	actor, err := adminService.ResolveActor(ctx, Authn{BearerToken: adminSession.AccessToken}, domain.PermissionSMTPSettingsUpdate)
	if err != nil {
		t.Fatal(err)
	}
	settings, err := adminService.UpdateSMTPSettings(ctx, actor, SMTPSettingsInput{
		Host:         "smtp.example.com",
		Port:         587,
		Username:     "mailer",
		Password:     "smtp-secret",
		Encryption:   domain.SmtpEncryptionStartTLS,
		FromEmail:    "no-reply@example.com",
		FromName:     "CodeLinks",
		ReplyToEmail: "support@example.com",
		Active:       true,
	}, RequestMeta{IPAddress: "192.0.2.10"})
	if err != nil {
		t.Fatal(err)
	}
	if !settings.HasPassword || settings.PasswordEncrypted == "smtp-secret" {
		t.Fatalf("expected encrypted password metadata, got %#v", settings)
	}

	if err := adminService.SendTestEmail(ctx, actor, "admin@example.com", RequestMeta{IPAddress: "192.0.2.10"}); err != nil {
		t.Fatal(err)
	}
	if email.password != "smtp-secret" || email.message.To != "admin@example.com" {
		t.Fatalf("unexpected test email call: %#v password=%q", email.message, email.password)
	}
}

type staticRuntimeLogs struct{}

func (staticRuntimeLogs) ListRuntimeLogs(context.Context, int) ([]domain.RuntimeLogEntry, error) {
	return []domain.RuntimeLogEntry{}, nil
}

type spySender struct {
	message  appmail.Message
	password string
}

func (s *spySender) Send(_ context.Context, _ domain.SmtpSettings, message appmail.Message, password string) error {
	s.message = message
	s.password = password
	return nil
}

func adminKind(err error, kind Kind) bool {
	var serviceError *Error
	return errors.As(err, &serviceError) && serviceError.Kind == kind
}
