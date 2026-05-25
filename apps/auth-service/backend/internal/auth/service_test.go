package auth

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/besart951/code-links/apps/auth-service/backend/internal/store/memory"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/token"
)

func TestServiceSignupLoginRefreshAndReset(t *testing.T) {
	service := newTestService(t)
	ctx := context.Background()

	signup, err := service.Signup(ctx, SignupInput{
		Name:          "New User",
		Email:         "new@example.com",
		Password:      "password1234",
		AcceptedTerms: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if signup.DebugVerificationToken == "" {
		t.Fatal("expected debug verification token")
	}

	if _, err := service.Login(ctx, LoginInput{Email: "new@example.com", Password: "password1234"}); !isKind(err, KindForbidden) {
		t.Fatalf("expected unverified login forbidden, got %v", err)
	}
	if _, err := service.VerifyEmail(ctx, signup.DebugVerificationToken); err != nil {
		t.Fatal(err)
	}
	session, err := service.Login(ctx, LoginInput{Email: "new@example.com", Password: "password1234"})
	if err != nil {
		t.Fatal(err)
	}
	if session.AccessToken == "" || session.RefreshToken == "" {
		t.Fatal("expected issued session")
	}

	refreshed, err := service.Refresh(ctx, session.RefreshToken)
	if err != nil {
		t.Fatal(err)
	}
	if refreshed.AccessToken == "" || refreshed.RefreshToken == "" {
		t.Fatal("expected refreshed session")
	}
	if _, err := service.Refresh(ctx, session.RefreshToken); !isKind(err, KindRefreshReuse) {
		t.Fatalf("expected consumed refresh token to be rejected, got %v", err)
	}

	reset, err := service.ForgotPassword(ctx, "demo@codelinks.dev")
	if err != nil {
		t.Fatal(err)
	}
	if reset.DebugResetToken == "" {
		t.Fatal("expected debug reset token")
	}
	if err := service.ResetPassword(ctx, reset.DebugResetToken, "newpassword123"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Login(ctx, LoginInput{Email: "demo@codelinks.dev", Password: "newpassword123"}); err != nil {
		t.Fatal(err)
	}
}

func TestLoginProtectionRateLimitsByEmail(t *testing.T) {
	service := newTestService(t)
	ctx := context.Background()
	input := LoginInput{
		Email:    "demo@codelinks.dev",
		Password: "wrong-password",
		Attempt:  LoginAttemptMetadata{IPAddress: "192.0.2.10", CountryCode: "CH"},
	}

	for attempt := 0; attempt < 5; attempt++ {
		if _, err := service.Login(ctx, input); !isKind(err, KindUnauthorized) {
			t.Fatalf("expected unauthorized attempt %d, got %v", attempt+1, err)
		}
	}
	if _, err := service.Login(ctx, input); !isKind(err, KindRateLimited) {
		t.Fatalf("expected rate limit, got %v", err)
	}
}

func TestLoginProtectionRateLimitsByIP(t *testing.T) {
	service := newTestService(t)
	ctx := context.Background()

	for attempt := 0; attempt < 10; attempt++ {
		_, err := service.Login(ctx, LoginInput{
			Email:    fmt.Sprintf("unknown-%d@example.com", attempt),
			Password: "wrong-password",
			Attempt:  LoginAttemptMetadata{IPAddress: "192.0.2.11", CountryCode: "CH"},
		})
		if !isKind(err, KindUnauthorized) {
			t.Fatalf("expected unauthorized attempt %d, got %v", attempt+1, err)
		}
	}
	_, err := service.Login(ctx, LoginInput{
		Email:    "another-unknown@example.com",
		Password: "wrong-password",
		Attempt:  LoginAttemptMetadata{IPAddress: "192.0.2.11", CountryCode: "CH"},
	})
	if !isKind(err, KindRateLimited) {
		t.Fatalf("expected IP rate limit, got %v", err)
	}
}

func TestExpiredTemporaryLockIsClearedOnSuccessfulLogin(t *testing.T) {
	service, store := newTestServiceWithStore(t)
	ctx := context.Background()
	user, _, err := store.FindUserByEmail(ctx, "demo@codelinks.dev")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SetUserTemporaryLock(ctx, user.ID, time.Now().UTC().Add(-time.Minute)); err != nil {
		t.Fatal(err)
	}

	if _, err := service.Login(ctx, LoginInput{
		Email:    "demo@codelinks.dev",
		Password: "password",
		Attempt:  LoginAttemptMetadata{IPAddress: "192.0.2.12", CountryCode: "CH"},
	}); err != nil {
		t.Fatal(err)
	}

	user, _, err = store.FindUserByEmail(ctx, "demo@codelinks.dev")
	if err != nil {
		t.Fatal(err)
	}
	if user.LockedUntil != nil {
		t.Fatalf("expected temporary lock to be cleared, got %v", user.LockedUntil)
	}
}

func TestRefreshTokenReuseRevokesActiveSessions(t *testing.T) {
	service := newTestService(t)
	ctx := context.Background()
	first, err := service.Login(ctx, LoginInput{
		Email:    "demo@codelinks.dev",
		Password: "password",
		Attempt:  LoginAttemptMetadata{IPAddress: "192.0.2.13", CountryCode: "CH"},
	})
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Refresh(ctx, first.RefreshToken)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Refresh(ctx, first.RefreshToken); !isKind(err, KindRefreshReuse) {
		t.Fatalf("expected refresh reuse error, got %v", err)
	}
	if _, err := service.Refresh(ctx, second.RefreshToken); !isKind(err, KindUnauthorized) {
		t.Fatalf("expected revoked active refresh token to fail, got %v", err)
	}
}

func newTestService(t *testing.T) *Service {
	t.Helper()

	service, _ := newTestServiceWithStore(t)
	return service
}

func newTestServiceWithStore(t *testing.T) (*Service, *memory.Store) {
	t.Helper()

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
	return NewService(Config{
		Environment:          "test",
		PublicFrontendURL:    "http://auth.codelinks.localhost",
		RefreshTokenLifetime: time.Hour,
	}, store, store, store, signer), store
}

func isKind(err error, kind Kind) bool {
	var serviceError *Error
	return errors.As(err, &serviceError) && serviceError.Kind == kind
}
