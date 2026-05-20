package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/besart951/code_links/platform/internal/access"
	"github.com/besart951/code_links/platform/internal/admin"
	"github.com/besart951/code_links/platform/internal/auth"
	"github.com/besart951/code_links/platform/internal/productclients"
)

func TestLoginSetsHttpOnlyRefreshAndReadableCSRFCookie(t *testing.T) {
	now := time.Date(2026, 5, 20, 10, 0, 0, 0, time.UTC)
	handler := NewServer(Server{
		Login: fakeLogin{result: auth.LoginResult{
			User:             auth.User{ID: "user_1"},
			Session:          auth.Session{ID: "session_1", ExpiresAt: now.Add(time.Hour)},
			SessionToken:     "session_raw",
			RefreshToken:     "refresh_raw",
			CSRFToken:        "csrf_raw",
			RefreshExpiresAt: now.Add(time.Hour),
		}},
		Config: Config{CookieDomain: ".codelinks.localhost", CookieSecure: true},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"email":"owner@example.com","password":"secret"}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	headers := rec.Header().Values("Set-Cookie")
	joined := strings.Join(headers, "\n")
	if !strings.Contains(joined, "platform_session=session_raw") || !strings.Contains(joined, "refresh_token=refresh_raw") || !strings.Contains(joined, "HttpOnly") || !strings.Contains(joined, "Secure") {
		t.Fatalf("session/refresh cookie missing secure attributes:\n%s", joined)
	}
	if !strings.Contains(joined, "SameSite=Lax") {
		t.Fatalf("expected SameSite=Lax:\n%s", joined)
	}
	for _, header := range headers {
		if strings.HasPrefix(header, "csrf_token=") && strings.Contains(header, "HttpOnly") {
			t.Fatalf("csrf token cookie must be readable for double-submit strategy: %s", header)
		}
	}
}

func TestTokenExchangeDerivesIdentityFromSessionAndProductClient(t *testing.T) {
	now := time.Date(2026, 5, 20, 10, 0, 0, 0, time.UTC)
	issuer := &fakeIssueToken{}
	handler := NewServer(Server{
		Authenticate: fakeAuthenticate{result: auth.AuthenticatedSession{
			User:    auth.User{ID: "user_1"},
			Session: auth.Session{ID: "session_1", ExpiresAt: now.Add(time.Hour)},
		}},
		ProductClients: fakeProductClients{client: productclients.ProductClient{
			ID:       "infra_link",
			Audience: "infra_link",
		}},
		IssueToken: issuer,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", strings.NewReader(`{"user_id":"attacker","session_id":"evil","audience":"planer_link"}`))
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session_raw"})
	req.AddCookie(&http.Cookie{Name: selectedTenantCookieName, Value: "tenant_1"})
	req.Header.Set("X-Product-Key", "infra_link")
	req.Header.Set("Authorization", "Bearer product-secret")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	if issuer.input.SessionID != "session_1" || issuer.input.TenantID != "tenant_1" || issuer.input.Audience != "infra_link" {
		t.Fatalf("token exchange trusted request body instead of authenticated context: %#v", issuer.input)
	}
}

func TestAdminMeRequiresSession(t *testing.T) {
	handler := NewServer(Server{Admin: &fakeAdmin{allow: true}})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/me", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestAdminMeRequiresSuperadminPermission(t *testing.T) {
	handler := NewServer(Server{
		Authenticate: fakeAuthenticate{result: auth.AuthenticatedSession{
			User:    auth.User{ID: "user_1"},
			Session: auth.Session{ID: "session_1", ExpiresAt: time.Now().Add(time.Hour)},
		}},
		Admin: &fakeAdmin{allow: false},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/me", nil)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session_raw"})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestAdminDashboardUsesAuthenticatedSuperadmin(t *testing.T) {
	adminUseCase := &fakeAdmin{allow: true}
	handler := NewServer(Server{
		Authenticate: fakeAuthenticate{result: auth.AuthenticatedSession{
			User:    auth.User{ID: "user_1", Email: "owner@example.com", DisplayName: "Owner"},
			Session: auth.Session{ID: "session_1", ExpiresAt: time.Now().Add(time.Hour)},
		}},
		Admin: adminUseCase,
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/dashboard?user_id=attacker", nil)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session_raw"})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	if adminUseCase.lastUser.ID != "user_1" {
		t.Fatalf("admin auth trusted request input: %#v", adminUseCase.lastUser)
	}
}

type fakeLogin struct {
	result auth.LoginResult
}

func (f fakeLogin) Execute(ctx context.Context, input auth.LoginInput) (auth.LoginResult, error) {
	return f.result, nil
}

type fakeAuthenticate struct {
	result auth.AuthenticatedSession
}

func (f fakeAuthenticate) Execute(ctx context.Context, rawSessionToken string) (auth.AuthenticatedSession, error) {
	return f.result, nil
}

type fakeProductClients struct {
	client productclients.ProductClient
}

func (f fakeProductClients) Authenticate(ctx context.Context, productKey, bearerToken string) (productclients.ProductClient, error) {
	return f.client, nil
}

type fakeIssueToken struct {
	input access.IssueAccessTokenInput
}

func (f *fakeIssueToken) Execute(ctx context.Context, input access.IssueAccessTokenInput) (access.IssuedToken, error) {
	f.input = input
	return access.IssuedToken{Value: "nested", ExpiresAt: time.Now().Add(time.Minute), JWTID: "token_1"}, nil
}

type fakeAdmin struct {
	allow    bool
	lastUser auth.User
}

func (f *fakeAdmin) Me(ctx context.Context, user auth.User) (admin.Me, error) {
	f.lastUser = user
	if !f.allow {
		return admin.Me{}, admin.ErrSuperadminRequired
	}
	return admin.Me{User: admin.UserToSummary(user), Permissions: []string{string(admin.PermissionRead)}, Superadmin: true}, nil
}

func (f *fakeAdmin) Dashboard(ctx context.Context) (admin.DashboardSummary, error) {
	return admin.DashboardSummary{GeneratedAt: time.Now()}, nil
}

func (f *fakeAdmin) Search(ctx context.Context, query string, limit, offset int) (admin.SearchResponse, error) {
	return admin.SearchResponse{}, nil
}

func (f *fakeAdmin) ListTenants(ctx context.Context, limit, offset int) (admin.Page[admin.TenantSummary], error) {
	return admin.Page[admin.TenantSummary]{}, nil
}

func (f *fakeAdmin) GetTenant(ctx context.Context, tenantID string) (admin.TenantSummary, error) {
	return admin.TenantSummary{}, nil
}

func (f *fakeAdmin) ListUsers(ctx context.Context, limit, offset int) (admin.Page[admin.UserSummary], error) {
	return admin.Page[admin.UserSummary]{}, nil
}

func (f *fakeAdmin) GetUser(ctx context.Context, userID string) (admin.UserSummary, error) {
	return admin.UserSummary{}, nil
}

func (f *fakeAdmin) ListProducts(ctx context.Context) ([]admin.ProductSummary, error) {
	return nil, nil
}

func (f *fakeAdmin) ListSubscriptions(ctx context.Context, limit, offset int) (admin.Page[admin.SubscriptionSummary], error) {
	return admin.Page[admin.SubscriptionSummary]{}, nil
}

func (f *fakeAdmin) ListAudit(ctx context.Context, limit, offset int) (admin.Page[admin.AuditLogEntry], error) {
	return admin.Page[admin.AuditLogEntry]{}, nil
}

func (f *fakeAdmin) ListNotifications(ctx context.Context) (admin.NotificationsSummary, error) {
	return admin.NotificationsSummary{}, nil
}

func (f *fakeAdmin) ListSecurity(ctx context.Context, limit, offset int) (admin.Page[admin.SecurityEventSummary], error) {
	return admin.Page[admin.SecurityEventSummary]{}, nil
}

func (f *fakeAdmin) ListSettings(ctx context.Context, limit, offset int) (admin.Page[admin.Setting], error) {
	return admin.Page[admin.Setting]{}, nil
}
