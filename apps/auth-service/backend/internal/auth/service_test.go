package auth

import (
	"context"
	"errors"
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
	if _, err := service.Refresh(ctx, session.RefreshToken); !isKind(err, KindUnauthorized) {
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

func newTestService(t *testing.T) *Service {
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
	}, store, store, store, signer)
}

func isKind(err error, kind Kind) bool {
	var serviceError *Error
	return errors.As(err, &serviceError) && serviceError.Kind == kind
}
