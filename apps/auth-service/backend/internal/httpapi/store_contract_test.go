package httpapi

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"testing"
	"time"

	adminsvc "github.com/besart951/code-links/apps/auth-service/backend/internal/admin"
	authsvc "github.com/besart951/code-links/apps/auth-service/backend/internal/auth"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/domain"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/secret"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/store/memory"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/store/postgres"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type contractStore interface {
	authsvc.Store
	authsvc.RoleStore
	adminsvc.Store
	adminsvc.SettingsStore
}

func TestStoreContractMemory(t *testing.T) {
	store, err := memory.New()
	if err != nil {
		t.Fatal(err)
	}

	runStoreContract(t, store)
}

func TestStoreContractPostgres(t *testing.T) {
	databaseURL := os.Getenv("CODELINKS_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set CODELINKS_TEST_DATABASE_URL to run Postgres store contract tests")
	}

	store, cleanup, err := postgres.Open(context.Background(), databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	runStoreContract(t, store)
}

func runStoreContract(t *testing.T, store contractStore) {
	t.Helper()
	ctx := context.Background()
	now := time.Now().UTC()
	email := fmt.Sprintf("contract-%s@example.com", uuid.NewString())
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password1234"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	user, err := store.CreateUser(ctx, "Contract User", email, string(passwordHash))
	if err != nil {
		t.Fatal(err)
	}
	if user.ID == uuid.Nil || user.Email != email {
		t.Fatalf("unexpected created user: %#v", user)
	}

	if _, _, err := store.FindUserByEmail(ctx, email); err != nil {
		t.Fatal(err)
	}
	if _, err := store.GrantLicense(ctx, user.ID, "infra-link"); err != nil {
		t.Fatal(err)
	}
	_, licenses, err := store.FindUserByID(ctx, user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(licenses, "infra-link") {
		t.Fatalf("expected infra-link license, got %#v", licenses)
	}

	verificationToken := secret.HashOpaqueToken(uuid.NewString())
	if err := store.CreateEmailVerificationToken(ctx, user.ID, verificationToken, now.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	verifiedUser, err := store.VerifyEmailToken(ctx, verificationToken, now)
	if err != nil {
		t.Fatal(err)
	}
	if verifiedUser.EmailVerifiedAt == nil {
		t.Fatal("expected email verification timestamp")
	}

	resetToken := secret.HashOpaqueToken(uuid.NewString())
	if err := store.CreatePasswordResetToken(ctx, user.ID, resetToken, now.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if _, err := store.ResetPasswordByToken(ctx, resetToken, "new-hash", now); err != nil {
		t.Fatal(err)
	}
	if _, err := store.ResetPasswordByToken(ctx, resetToken, "newer-hash", now); err == nil {
		t.Fatal("expected used reset token to fail")
	}

	refreshTokenHash := secret.HashRefreshToken(uuid.NewString())
	if err := store.CreateRefreshSession(ctx, refreshTokenHash, user.ID, now.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	foundUserID, err := store.FindRefreshSession(ctx, refreshTokenHash)
	if err != nil {
		t.Fatal(err)
	}
	if foundUserID != user.ID {
		t.Fatalf("expected refresh session for %s, got %s", user.ID, foundUserID)
	}
	if err := store.DeleteRefreshSession(ctx, refreshTokenHash); err != nil {
		t.Fatal(err)
	}

	ipHash := sha256.Sum256([]byte("192.0.2.10"))
	if err := store.RecordLoginAttempt(ctx, domain.LoginAttempt{
		ID:             uuid.New(),
		UserID:         &user.ID,
		EmailAttempted: email,
		OccurredAt:     now,
		IPAddress:      "192.0.2.10",
		IPHash:         base64.RawURLEncoding.EncodeToString(ipHash[:]),
		CountryCode:    "CH",
		Success:        true,
		AuthMethod:     "password",
		RiskScore:      8,
	}); err != nil {
		t.Fatal(err)
	}
	result, err := store.ListLoginAttempts(ctx, domain.LoginAttemptListQuery{UserID: &user.ID, Page: 1, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total == 0 {
		t.Fatal("expected recorded login attempt")
	}

	if _, err := store.SetUserRole(ctx, user.ID, domain.AdminRoleSupport); err != nil {
		t.Fatal(err)
	}
	users, err := store.ListAdminUsers(ctx, domain.AdminUserListQuery{
		Query:     "contract",
		Role:      string(domain.AdminRoleSupport),
		Status:    domain.UserStatusActive,
		Page:      1,
		PageSize:  1,
		Sort:      "email",
		Direction: "asc",
	})
	if err != nil {
		t.Fatal(err)
	}
	if users.Total != 1 || len(users.Items) != 1 {
		t.Fatalf("expected one filtered admin user, got total=%d items=%d", users.Total, len(users.Items))
	}
	if users.Items[0].Email != email || users.Items[0].PrimaryRole != domain.AdminRoleSupport {
		t.Fatalf("unexpected admin user projection: %#v", users.Items[0])
	}
	if users.Items[0].SuccessfulLoginCount != 1 || users.Items[0].FailedLoginCount != 0 {
		t.Fatalf("unexpected login counts: %#v", users.Items[0])
	}
}
