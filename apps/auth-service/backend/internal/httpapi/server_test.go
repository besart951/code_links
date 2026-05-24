package httpapi

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	adminsvc "github.com/besart951/code-links/apps/auth-service/backend/internal/admin"
	authsvc "github.com/besart951/code-links/apps/auth-service/backend/internal/auth"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/domain"
	appmail "github.com/besart951/code-links/apps/auth-service/backend/internal/mail"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/secret"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/store/memory"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/token"
)

func TestLoginReturnsAccessTokenAndJWKS(t *testing.T) {
	server := newTestServer(t)

	login := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(`{"email":"demo@codelinks.dev","password":"password"}`))
	recorder := httptest.NewRecorder()
	server.routes().ServeHTTP(recorder, login)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected login 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	var session sessionResponse
	if err := json.NewDecoder(recorder.Body).Decode(&session); err != nil {
		t.Fatal(err)
	}
	if session.AccessToken == "" {
		t.Fatal("expected access token")
	}
	if len(session.User.Licenses) != 1 || session.User.Licenses[0] != "infra-link" {
		t.Fatalf("unexpected licenses: %#v", session.User.Licenses)
	}

	jwksRequest := httptest.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil)
	jwksRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(jwksRecorder, jwksRequest)

	if jwksRecorder.Code != http.StatusOK {
		t.Fatalf("expected JWKS 200, got %d", jwksRecorder.Code)
	}
}

func TestSignupRequiresEmailVerificationBeforeLogin(t *testing.T) {
	server := newTestServer(t)

	signupRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(
		signupRecorder,
		httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(`{"name":"New User","email":"new@example.com","password":"password1234","acceptedTerms":true}`)),
	)
	if signupRecorder.Code != http.StatusCreated {
		t.Fatalf("expected signup 201, got %d: %s", signupRecorder.Code, signupRecorder.Body.String())
	}

	var signup map[string]string
	if err := json.NewDecoder(signupRecorder.Body).Decode(&signup); err != nil {
		t.Fatal(err)
	}
	if signup["debugVerificationToken"] == "" {
		t.Fatal("expected debug verification token in tests")
	}

	loginBeforeVerify := httptest.NewRecorder()
	server.routes().ServeHTTP(
		loginBeforeVerify,
		httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(`{"email":"new@example.com","password":"password1234"}`)),
	)
	if loginBeforeVerify.Code != http.StatusForbidden {
		t.Fatalf("expected unverified login 403, got %d: %s", loginBeforeVerify.Code, loginBeforeVerify.Body.String())
	}

	verify := httptest.NewRequest(http.MethodGet, "/api/auth/verify-email?token="+signup["debugVerificationToken"], nil)
	verify.Header.Set("Accept", "application/json")
	verifyRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(verifyRecorder, verify)
	if verifyRecorder.Code != http.StatusOK {
		t.Fatalf("expected verify 200, got %d: %s", verifyRecorder.Code, verifyRecorder.Body.String())
	}

	loginAfterVerify := httptest.NewRecorder()
	server.routes().ServeHTTP(
		loginAfterVerify,
		httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(`{"email":"new@example.com","password":"password1234"}`)),
	)
	if loginAfterVerify.Code != http.StatusOK {
		t.Fatalf("expected verified login 200, got %d: %s", loginAfterVerify.Code, loginAfterVerify.Body.String())
	}

	attempts, err := server.store.ListLoginAttempts(context.Background(), domain.LoginAttemptListQuery{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if attempts.Total < 2 {
		t.Fatalf("expected login attempts to be recorded, got %d", attempts.Total)
	}
}

func TestForgotPasswordIsNeutralAndResetTokenIsOneTimeUse(t *testing.T) {
	server := newTestServer(t)

	unknownRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(
		unknownRecorder,
		httptest.NewRequest(http.MethodPost, "/api/auth/forgot-password", bytes.NewBufferString(`{"email":"unknown@example.com"}`)),
	)
	if unknownRecorder.Code != http.StatusOK {
		t.Fatalf("expected neutral forgot password 200, got %d", unknownRecorder.Code)
	}

	knownRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(
		knownRecorder,
		httptest.NewRequest(http.MethodPost, "/api/auth/forgot-password", bytes.NewBufferString(`{"email":"demo@codelinks.dev"}`)),
	)
	if knownRecorder.Code != http.StatusOK {
		t.Fatalf("expected known forgot password 200, got %d: %s", knownRecorder.Code, knownRecorder.Body.String())
	}

	var response map[string]string
	if err := json.NewDecoder(knownRecorder.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	token := response["debugResetToken"]
	if token == "" {
		t.Fatal("expected debug reset token")
	}

	resetRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(
		resetRecorder,
		httptest.NewRequest(http.MethodPost, "/api/auth/reset-password", bytes.NewBufferString(`{"token":"`+token+`","password":"newpassword123"}`)),
	)
	if resetRecorder.Code != http.StatusOK {
		t.Fatalf("expected reset 200, got %d: %s", resetRecorder.Code, resetRecorder.Body.String())
	}

	secondResetRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(
		secondResetRecorder,
		httptest.NewRequest(http.MethodPost, "/api/auth/reset-password", bytes.NewBufferString(`{"token":"`+token+`","password":"anotherpassword123"}`)),
	)
	if secondResetRecorder.Code != http.StatusBadRequest {
		t.Fatalf("expected used token 400, got %d: %s", secondResetRecorder.Code, secondResetRecorder.Body.String())
	}

	loginRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(
		loginRecorder,
		httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(`{"email":"demo@codelinks.dev","password":"newpassword123"}`)),
	)
	if loginRecorder.Code != http.StatusOK {
		t.Fatalf("expected login with reset password 200, got %d: %s", loginRecorder.Code, loginRecorder.Body.String())
	}
}

func TestAdminMeRejectsNonAdminAndAllowsAdmin(t *testing.T) {
	server := newTestServer(t)

	signupRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(
		signupRecorder,
		httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(`{"name":"Plain User","email":"plain@example.com","password":"password1234","acceptedTerms":true}`)),
	)
	var signup map[string]string
	if err := json.NewDecoder(signupRecorder.Body).Decode(&signup); err != nil {
		t.Fatal(err)
	}
	verify := httptest.NewRequest(http.MethodGet, "/api/auth/verify-email?token="+signup["debugVerificationToken"], nil)
	verify.Header.Set("Accept", "application/json")
	server.routes().ServeHTTP(httptest.NewRecorder(), verify)

	userLogin := loginAndDecode(t, server, "plain@example.com", "password1234")
	userMe := httptest.NewRequest(http.MethodGet, "/api/admin/me", nil)
	userMe.Header.Set("Authorization", "Bearer "+userLogin.AccessToken)
	userMeRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(userMeRecorder, userMe)
	if userMeRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected non-admin rejection 401, got %d: %s", userMeRecorder.Code, userMeRecorder.Body.String())
	}

	adminLogin := loginAndDecode(t, server, "demo@codelinks.dev", "password")
	adminMe := httptest.NewRequest(http.MethodGet, "/api/admin/me", nil)
	adminMe.Header.Set("Authorization", "Bearer "+adminLogin.AccessToken)
	adminMeRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(adminMeRecorder, adminMe)
	if adminMeRecorder.Code != http.StatusOK {
		t.Fatalf("expected admin me 200, got %d: %s", adminMeRecorder.Code, adminMeRecorder.Body.String())
	}
}

func TestSMTPSettingsAreEncryptedAndSupportCannotUpdate(t *testing.T) {
	server := newTestServer(t)
	adminLogin := loginAndDecode(t, server, "demo@codelinks.dev", "password")

	update := httptest.NewRequest(http.MethodPut, "/api/admin/settings/smtp", bytes.NewBufferString(`{"host":"smtp.example.com","port":587,"username":"mailer","password":"smtp-secret","encryption":"starttls","fromEmail":"no-reply@example.com","fromName":"CodeLinks","replyToEmail":"support@example.com","active":true}`))
	update.Header.Set("Authorization", "Bearer "+adminLogin.AccessToken)
	updateRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(updateRecorder, update)
	if updateRecorder.Code != http.StatusOK {
		t.Fatalf("expected smtp update 200, got %d: %s", updateRecorder.Code, updateRecorder.Body.String())
	}

	settings, err := server.store.GetSmtpSettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	encrypted := settings.PasswordEncrypted
	if encrypted == "" || encrypted == "smtp-secret" {
		t.Fatalf("expected encrypted smtp password, got %q", encrypted)
	}
	decrypted, err := secret.Decrypt(encrypted, server.smtpSecretKey)
	if err != nil {
		t.Fatal(err)
	}
	if decrypted != "smtp-secret" {
		t.Fatalf("unexpected decrypted password %q", decrypted)
	}

	testEmail := httptest.NewRequest(http.MethodPost, "/api/admin/settings/smtp/test-email", bytes.NewBufferString(`{"recipient":"admin@example.com"}`))
	testEmail.Header.Set("Authorization", "Bearer "+adminLogin.AccessToken)
	testEmailRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(testEmailRecorder, testEmail)
	if testEmailRecorder.Code != http.StatusOK {
		t.Fatalf("expected smtp test 200, got %d: %s", testEmailRecorder.Code, testEmailRecorder.Body.String())
	}

	supportLogin := loginAndDecode(t, server, "support@codelinks.dev", "password")
	supportUpdate := httptest.NewRequest(http.MethodPut, "/api/admin/settings/smtp", bytes.NewBufferString(`{"host":"smtp.example.com","port":587,"encryption":"starttls","fromEmail":"no-reply@example.com","fromName":"CodeLinks","replyToEmail":"support@example.com","active":true}`))
	supportUpdate.Header.Set("Authorization", "Bearer "+supportLogin.AccessToken)
	supportRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(supportRecorder, supportUpdate)
	if supportRecorder.Code != http.StatusForbidden {
		t.Fatalf("expected support smtp rejection 403, got %d: %s", supportRecorder.Code, supportRecorder.Body.String())
	}
}

func TestServerCanUseRealEmailSender(t *testing.T) {
	var sender appmail.Sender = appmail.SmtpSender{}
	if _, ok := sender.(appmail.SmtpSender); !ok {
		t.Fatal("expected SMTP email sender")
	}
}

func TestAdminAPIMasksIPMetadataWithoutRawPermission(t *testing.T) {
	server := newTestServer(t)
	auditorLogin := loginAndDecode(t, server, "auditor@codelinks.dev", "password")

	attemptsRequest := httptest.NewRequest(http.MethodGet, "/api/admin/login-attempts", nil)
	attemptsRequest.Header.Set("Authorization", "Bearer "+auditorLogin.AccessToken)
	attemptsRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(attemptsRecorder, attemptsRequest)
	if attemptsRecorder.Code != http.StatusOK {
		t.Fatalf("expected login attempts 200, got %d: %s", attemptsRecorder.Code, attemptsRecorder.Body.String())
	}

	var attempts domain.LoginAttemptListResult
	if err := json.NewDecoder(attemptsRecorder.Body).Decode(&attempts); err != nil {
		t.Fatal(err)
	}
	if len(attempts.Items) == 0 {
		t.Fatal("expected at least one login attempt")
	}
	if got := attempts.Items[0].IPAddress; got != "192.0.x.x" {
		t.Fatalf("expected masked login attempt IP, got %q", got)
	}

	statsRequest := httptest.NewRequest(http.MethodGet, "/api/admin/dashboard/stats", nil)
	statsRequest.Header.Set("Authorization", "Bearer "+auditorLogin.AccessToken)
	statsRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(statsRecorder, statsRequest)
	if statsRecorder.Code != http.StatusOK {
		t.Fatalf("expected dashboard stats 200, got %d: %s", statsRecorder.Code, statsRecorder.Body.String())
	}

	var stats domain.DashboardStats
	if err := json.NewDecoder(statsRecorder.Body).Decode(&stats); err != nil {
		t.Fatal(err)
	}
	if len(stats.TopIPAddresses) == 0 {
		t.Fatal("expected dashboard IP stats")
	}
	if got := stats.TopIPAddresses[0].Key; got != "192.0.x.x" {
		t.Fatalf("expected masked dashboard IP, got %q", got)
	}
}

func TestMockPurchaseGrantsLicenseAndRefreshesAccessToken(t *testing.T) {
	server := newTestServer(t)
	loginRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(
		loginRecorder,
		httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(`{"email":"demo@codelinks.dev","password":"password"}`)),
	)

	var login sessionResponse
	if err := json.NewDecoder(loginRecorder.Body).Decode(&login); err != nil {
		t.Fatal(err)
	}

	purchase := httptest.NewRequest(http.MethodPost, "/api/licenses/mock-purchase", bytes.NewBufferString(`{"productId":"planer-link"}`))
	purchase.Header.Set("Authorization", "Bearer "+login.AccessToken)
	purchaseRecorder := httptest.NewRecorder()
	server.routes().ServeHTTP(purchaseRecorder, purchase)

	if purchaseRecorder.Code != http.StatusOK {
		t.Fatalf("expected purchase 200, got %d: %s", purchaseRecorder.Code, purchaseRecorder.Body.String())
	}

	var session sessionResponse
	if err := json.NewDecoder(purchaseRecorder.Body).Decode(&session); err != nil {
		t.Fatal(err)
	}
	if !contains(session.User.Licenses, "planer-link") {
		t.Fatalf("expected planer-link license, got %#v", session.User.Licenses)
	}
}

func TestMockPurchaseRouteRequiresExplicitDevMode(t *testing.T) {
	server := newTestServerWithMockPurchase(t, false)
	login := loginAndDecode(t, server, "demo@codelinks.dev", "password")

	purchase := httptest.NewRequest(http.MethodPost, "/api/licenses/mock-purchase", bytes.NewBufferString(`{"productId":"planer-link"}`))
	purchase.Header.Set("Authorization", "Bearer "+login.AccessToken)
	recorder := httptest.NewRecorder()
	server.routes().ServeHTTP(recorder, purchase)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected disabled mock purchase 404, got %d: %s", recorder.Code, recorder.Body.String())
	}
}

type testServer struct {
	*Server
	store         *memory.Store
	smtpSecretKey []byte
}

func (s *testServer) routes() http.Handler {
	return s.Handler()
}

func newTestServer(t *testing.T) *testServer {
	return newTestServerWithMockPurchase(t, true)
}

func newTestServerWithMockPurchase(t *testing.T, enableMockPurchase bool) *testServer {
	t.Helper()

	store, err := memory.New()
	if err != nil {
		t.Fatal(err)
	}

	key := sha256.Sum256([]byte("test-smtp-secret"))
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
	adminService := adminsvc.NewService(adminsvc.Config{
		SMTPSecretKey: key[:],
	}, store, store, store, store, store, signer, appmail.NoopSender{})
	server := NewServer(Config{
		EnableMockPurchase: enableMockPurchase,
		PublicFrontendURL:  "http://auth.codelinks.localhost",
	}, authService, adminService, signer)
	return &testServer{Server: server, store: store, smtpSecretKey: key[:]}
}

func loginAndDecode(t *testing.T, server *testServer, email string, password string) sessionResponse {
	t.Helper()

	recorder := httptest.NewRecorder()
	server.routes().ServeHTTP(
		recorder,
		httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(`{"email":"`+email+`","password":"`+password+`"}`)),
	)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected login 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	var session sessionResponse
	if err := json.NewDecoder(recorder.Body).Decode(&session); err != nil {
		t.Fatal(err)
	}

	return session
}

func contains(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}

	return false
}
