package httpapi

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/besart951/code_links/platform/internal/access"
	"github.com/besart951/code_links/platform/internal/admin"
	"github.com/besart951/code_links/platform/internal/auth"
	"github.com/besart951/code_links/platform/internal/productclients"
	"github.com/besart951/code_links/platform/internal/tenants"
)

const (
	sessionCookieName        = "platform_session"
	refreshCookieName        = "refresh_token"
	csrfCookieName           = "csrf_token"
	selectedTenantCookieName = "codelinks_tenant_id"
)

type LoginUseCase interface {
	Execute(ctx context.Context, input auth.LoginInput) (auth.LoginResult, error)
}

type RefreshUseCase interface {
	Execute(ctx context.Context, input auth.RefreshInput) (auth.RefreshResult, error)
}

type LogoutUseCase interface {
	Execute(ctx context.Context, sessionID auth.SessionID) error
}

type AuthenticateSessionUseCase interface {
	Execute(ctx context.Context, rawSessionToken string) (auth.AuthenticatedSession, error)
}

type SelectTenantUseCase interface {
	Execute(ctx context.Context, input tenants.SelectTenantInput) (tenants.SelectTenantResult, error)
}

type IssueTokenUseCase interface {
	Execute(ctx context.Context, input access.IssueAccessTokenInput) (access.IssuedToken, error)
}

type CheckProductUseCase interface {
	Execute(ctx context.Context, input access.CheckProductAccessInput) (access.AccessDecision, error)
}

type CheckFeatureUseCase interface {
	Execute(ctx context.Context, input access.CheckFeatureAccessInput) (access.AccessDecision, error)
}

type ProductClientAuthenticator interface {
	Authenticate(ctx context.Context, productKey, bearerToken string) (productclients.ProductClient, error)
}

type TenantLister interface {
	ListTenantsForUser(ctx context.Context, userID tenants.UserID) ([]tenants.Tenant, error)
}

type ProductLister interface {
	ListProducts(ctx context.Context) ([]access.Product, error)
}

type EntitlementLister interface {
	ListEntitlements(ctx context.Context, tenantID access.TenantID) ([]access.Entitlement, error)
	ListFeatureLimits(ctx context.Context, tenantID access.TenantID) ([]access.FeatureLimit, error)
}

type AdminReadUseCase interface {
	Me(ctx context.Context, user auth.User) (admin.Me, error)
	Dashboard(ctx context.Context) (admin.DashboardSummary, error)
	Search(ctx context.Context, query string, limit, offset int) (admin.SearchResponse, error)
	ListTenants(ctx context.Context, limit, offset int) (admin.Page[admin.TenantSummary], error)
	GetTenant(ctx context.Context, tenantID string) (admin.TenantSummary, error)
	ListUsers(ctx context.Context, limit, offset int) (admin.Page[admin.UserSummary], error)
	GetUser(ctx context.Context, userID string) (admin.UserSummary, error)
	ListProducts(ctx context.Context) ([]admin.ProductSummary, error)
	ListSubscriptions(ctx context.Context, limit, offset int) (admin.Page[admin.SubscriptionSummary], error)
	ListAudit(ctx context.Context, limit, offset int) (admin.Page[admin.AuditLogEntry], error)
	ListNotifications(ctx context.Context) (admin.NotificationsSummary, error)
	ListSecurity(ctx context.Context, limit, offset int) (admin.Page[admin.SecurityEventSummary], error)
	ListSettings(ctx context.Context, limit, offset int) (admin.Page[admin.Setting], error)
}

type Config struct {
	CookieDomain string
	CookieSecure bool
}

type Server struct {
	Login          LoginUseCase
	Refresh        RefreshUseCase
	Logout         LogoutUseCase
	Authenticate   AuthenticateSessionUseCase
	SelectTenant   SelectTenantUseCase
	IssueToken     IssueTokenUseCase
	CheckProduct   CheckProductUseCase
	CheckFeature   CheckFeatureUseCase
	ProductClients ProductClientAuthenticator
	Tenants        TenantLister
	Products       ProductLister
	Entitlements   EntitlementLister
	Admin          AdminReadUseCase
	Config         Config
}

func NewServer(server Server) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", server.handleHealth)
	mux.HandleFunc("POST /api/v1/auth/login", server.handleLogin)
	mux.HandleFunc("POST /api/v1/auth/refresh", server.handleRefresh)
	mux.HandleFunc("POST /api/v1/auth/logout", server.handleLogout)
	mux.HandleFunc("GET /api/v1/auth/me", server.handleMe)
	mux.HandleFunc("POST /api/v1/auth/select-tenant", server.handleSelectTenant)
	mux.HandleFunc("GET /api/v1/auth/products", server.handleProducts)
	mux.HandleFunc("GET /api/v1/auth/entitlements", server.handleEntitlements)
	mux.HandleFunc("POST /api/v1/auth/token", server.handleIssueToken)
	mux.HandleFunc("POST /api/v1/auth/check-product-access", server.handleCheckProductAccess)
	mux.HandleFunc("POST /api/v1/auth/check-feature-access", server.handleCheckFeatureAccess)
	mux.HandleFunc("GET /api/v1/admin/me", server.handleAdminMe)
	mux.HandleFunc("GET /api/v1/admin/dashboard", server.handleAdminDashboard)
	mux.HandleFunc("GET /api/v1/admin/search", server.handleAdminSearch)
	mux.HandleFunc("GET /api/v1/admin/tenants", server.handleAdminTenants)
	mux.HandleFunc("GET /api/v1/admin/tenants/{tenantId}", server.handleAdminTenant)
	mux.HandleFunc("GET /api/v1/admin/users", server.handleAdminUsers)
	mux.HandleFunc("GET /api/v1/admin/users/{userId}", server.handleAdminUser)
	mux.HandleFunc("GET /api/v1/admin/products", server.handleAdminProducts)
	mux.HandleFunc("GET /api/v1/admin/subscriptions", server.handleAdminSubscriptions)
	mux.HandleFunc("GET /api/v1/admin/audit", server.handleAdminAudit)
	mux.HandleFunc("GET /api/v1/admin/notifications", server.handleAdminNotifications)
	mux.HandleFunc("GET /api/v1/admin/security", server.handleAdminSecurity)
	mux.HandleFunc("GET /api/v1/admin/settings", server.handleAdminSettings)
	return securityHeaders(mux)
}

func (s Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	UserID           string    `json:"user_id"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	CSRFToken        string    `json:"csrf_token"`
}

func (s Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if s.Login == nil {
		writeError(w, http.StatusNotImplemented, "login_not_configured")
		return
	}
	var req loginRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	result, err := s.Login.Execute(r.Context(), auth.LoginInput{
		Email:     req.Email,
		Password:  req.Password,
		UserAgent: r.UserAgent(),
		IP:        clientIP(r),
	})
	if err != nil {
		writeAuthError(w, err)
		return
	}
	s.setCookie(w, sessionCookieName, result.SessionToken, result.Session.ExpiresAt, true)
	s.setCookie(w, refreshCookieName, result.RefreshToken, result.RefreshExpiresAt, true)
	s.setCookie(w, csrfCookieName, result.CSRFToken, time.Now().Add(24*time.Hour), false)
	writeJSON(w, http.StatusOK, loginResponse{
		UserID:           string(result.User.ID),
		RefreshExpiresAt: result.RefreshExpiresAt,
		CSRFToken:        result.CSRFToken,
	})
}

func (s Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	if s.Refresh == nil {
		writeError(w, http.StatusNotImplemented, "refresh_not_configured")
		return
	}
	if !requireCSRF(w, r) {
		return
	}
	cookie, err := r.Cookie(refreshCookieName)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	result, err := s.Refresh.Execute(r.Context(), auth.RefreshInput{
		RefreshToken: cookie.Value,
		UserAgent:    r.UserAgent(),
		IP:           clientIP(r),
	})
	if err != nil {
		writeAuthError(w, err)
		return
	}
	s.setCookie(w, refreshCookieName, result.RefreshToken, result.RefreshExpiresAt, true)
	s.setCookie(w, csrfCookieName, result.CSRFToken, time.Now().Add(24*time.Hour), false)
	writeJSON(w, http.StatusOK, loginResponse{
		UserID:           string(result.User.ID),
		RefreshExpiresAt: result.RefreshExpiresAt,
		CSRFToken:        result.CSRFToken,
	})
}

func (s Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if !requireCSRF(w, r) {
		return
	}
	authenticated, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	if s.Logout != nil {
		if err := s.Logout.Execute(r.Context(), authenticated.Session.ID); err != nil {
			writeError(w, http.StatusInternalServerError, "logout_failed")
			return
		}
	}
	s.clearCookie(w, sessionCookieName, true)
	s.clearCookie(w, refreshCookieName, true)
	s.clearCookie(w, csrfCookieName, false)
	s.clearCookie(w, selectedTenantCookieName, true)
	w.WriteHeader(http.StatusNoContent)
}

func (s Server) handleMe(w http.ResponseWriter, r *http.Request) {
	authenticated, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	if s.Tenants == nil {
		writeError(w, http.StatusNotImplemented, "tenants_not_configured")
		return
	}
	items, err := s.Tenants.ListTenantsForUser(r.Context(), tenants.UserID(authenticated.User.ID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "tenant_lookup_failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user": map[string]any{
			"id":           authenticated.User.ID,
			"email":        authenticated.User.Email,
			"display_name": authenticated.User.DisplayName,
		},
		"tenants": items,
	})
}

type selectTenantRequest struct {
	TenantID string `json:"tenant_id"`
}

func (s Server) handleSelectTenant(w http.ResponseWriter, r *http.Request) {
	if !requireCSRF(w, r) {
		return
	}
	authenticated, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	if s.SelectTenant == nil {
		writeError(w, http.StatusNotImplemented, "select_tenant_not_configured")
		return
	}
	var req selectTenantRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	result, err := s.SelectTenant.Execute(r.Context(), tenants.SelectTenantInput{
		UserID:   tenants.UserID(authenticated.User.ID),
		TenantID: tenants.TenantID(strings.TrimSpace(req.TenantID)),
	})
	if err != nil {
		writeError(w, http.StatusForbidden, "tenant_selection_failed")
		return
	}
	s.setCookie(w, selectedTenantCookieName, string(result.Tenant.ID), authenticated.Session.ExpiresAt, true)
	writeJSON(w, http.StatusOK, map[string]any{
		"tenant_id": string(result.Tenant.ID),
		"type":      string(result.Tenant.Type),
	})
}

func (s Server) handleProducts(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	if s.Products == nil {
		writeError(w, http.StatusNotImplemented, "products_not_configured")
		return
	}
	products, err := s.Products.ListProducts(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "product_lookup_failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"products": products})
}

func (s Server) handleEntitlements(w http.ResponseWriter, r *http.Request) {
	authenticated, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	if s.Entitlements == nil || s.SelectTenant == nil {
		writeError(w, http.StatusNotImplemented, "entitlements_not_configured")
		return
	}
	tenantID, ok := s.selectedTenant(w, r)
	if !ok {
		return
	}
	if _, err := s.SelectTenant.Execute(r.Context(), tenants.SelectTenantInput{
		UserID:   tenants.UserID(authenticated.User.ID),
		TenantID: tenants.TenantID(tenantID),
	}); err != nil {
		writeError(w, http.StatusForbidden, "tenant_membership_required")
		return
	}
	entitlements, err := s.Entitlements.ListEntitlements(r.Context(), access.TenantID(tenantID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "entitlement_lookup_failed")
		return
	}
	limits, err := s.Entitlements.ListFeatureLimits(r.Context(), access.TenantID(tenantID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "feature_limit_lookup_failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"tenant_id":    tenantID,
		"entitlements": entitlements,
		"limits":       limits,
	})
}

func (s Server) handleIssueToken(w http.ResponseWriter, r *http.Request) {
	authenticated, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	client, ok := s.requireProductClient(w, r)
	if !ok {
		return
	}
	tenantID, ok := s.selectedTenant(w, r)
	if !ok {
		return
	}
	if s.IssueToken == nil {
		writeError(w, http.StatusNotImplemented, "token_issuer_not_configured")
		return
	}
	token, err := s.IssueToken.Execute(r.Context(), access.IssueAccessTokenInput{
		SessionID: access.SessionID(authenticated.Session.ID),
		TenantID:  access.TenantID(tenantID),
		Audience:  client.Audience,
	})
	if err != nil {
		writeAccessError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"access_token": token.Value,
		"expires_at":   token.ExpiresAt,
		"jti":          token.JWTID,
	})
}

type checkProductRequest struct {
	Product string `json:"product"`
}

func (s Server) handleCheckProductAccess(w http.ResponseWriter, r *http.Request) {
	if !requireCSRF(w, r) {
		return
	}
	authenticated, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	if s.CheckProduct == nil {
		writeError(w, http.StatusNotImplemented, "product_access_not_configured")
		return
	}
	tenantID, ok := s.selectedTenant(w, r)
	if !ok {
		return
	}
	var req checkProductRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	decision, err := s.CheckProduct.Execute(r.Context(), access.CheckProductAccessInput{
		UserID:   access.UserID(authenticated.User.ID),
		TenantID: access.TenantID(tenantID),
		Product:  access.ProductKey(strings.TrimSpace(req.Product)),
	})
	writeDecision(w, decision, err)
}

type checkFeatureRequest struct {
	Product            string `json:"product"`
	Feature            string `json:"feature"`
	RequiredPermission string `json:"required_permission"`
}

func (s Server) handleCheckFeatureAccess(w http.ResponseWriter, r *http.Request) {
	if !requireCSRF(w, r) {
		return
	}
	authenticated, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	if s.CheckFeature == nil {
		writeError(w, http.StatusNotImplemented, "feature_access_not_configured")
		return
	}
	tenantID, ok := s.selectedTenant(w, r)
	if !ok {
		return
	}
	var req checkFeatureRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	decision, err := s.CheckFeature.Execute(r.Context(), access.CheckFeatureAccessInput{
		UserID:             access.UserID(authenticated.User.ID),
		TenantID:           access.TenantID(tenantID),
		Product:            access.ProductKey(strings.TrimSpace(req.Product)),
		Feature:            access.FeatureKey(strings.TrimSpace(req.Feature)),
		RequiredPermission: access.PermissionKey(strings.TrimSpace(req.RequiredPermission)),
	})
	writeDecision(w, decision, err)
}

func (s Server) handleAdminMe(w http.ResponseWriter, r *http.Request) {
	_, me, ok := s.requireAdmin(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, me)
}

func (s Server) handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	result, err := s.Admin.Dashboard(r.Context())
	writeAdminResult(w, result, err)
}

func (s Server) handleAdminSearch(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	result, err := s.Admin.Search(r.Context(), r.URL.Query().Get("q"), queryInt(r, "limit"), queryInt(r, "offset"))
	writeAdminResult(w, result, err)
}

func (s Server) handleAdminTenants(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	result, err := s.Admin.ListTenants(r.Context(), queryInt(r, "limit"), queryInt(r, "offset"))
	writeAdminResult(w, result, err)
}

func (s Server) handleAdminTenant(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	result, err := s.Admin.GetTenant(r.Context(), r.PathValue("tenantId"))
	writeAdminResult(w, result, err)
}

func (s Server) handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	result, err := s.Admin.ListUsers(r.Context(), queryInt(r, "limit"), queryInt(r, "offset"))
	writeAdminResult(w, result, err)
}

func (s Server) handleAdminUser(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	result, err := s.Admin.GetUser(r.Context(), r.PathValue("userId"))
	writeAdminResult(w, result, err)
}

func (s Server) handleAdminProducts(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	products, err := s.Admin.ListProducts(r.Context())
	writeAdminResult(w, map[string]any{"products": products}, err)
}

func (s Server) handleAdminSubscriptions(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	result, err := s.Admin.ListSubscriptions(r.Context(), queryInt(r, "limit"), queryInt(r, "offset"))
	writeAdminResult(w, result, err)
}

func (s Server) handleAdminAudit(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	result, err := s.Admin.ListAudit(r.Context(), queryInt(r, "limit"), queryInt(r, "offset"))
	writeAdminResult(w, result, err)
}

func (s Server) handleAdminNotifications(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	result, err := s.Admin.ListNotifications(r.Context())
	writeAdminResult(w, result, err)
}

func (s Server) handleAdminSecurity(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	result, err := s.Admin.ListSecurity(r.Context(), queryInt(r, "limit"), queryInt(r, "offset"))
	writeAdminResult(w, result, err)
}

func (s Server) handleAdminSettings(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	result, err := s.Admin.ListSettings(r.Context(), queryInt(r, "limit"), queryInt(r, "offset"))
	writeAdminResult(w, result, err)
}

func (s Server) requireAdmin(w http.ResponseWriter, r *http.Request) (auth.AuthenticatedSession, admin.Me, bool) {
	authenticated, ok := s.requireSession(w, r)
	if !ok {
		return auth.AuthenticatedSession{}, admin.Me{}, false
	}
	if s.Admin == nil {
		writeError(w, http.StatusForbidden, "superadmin_required")
		return auth.AuthenticatedSession{}, admin.Me{}, false
	}
	me, err := s.Admin.Me(r.Context(), authenticated.User)
	if err != nil {
		writeAdminError(w, err)
		return auth.AuthenticatedSession{}, admin.Me{}, false
	}
	return authenticated, me, true
}

func (s Server) requireSession(w http.ResponseWriter, r *http.Request) (auth.AuthenticatedSession, bool) {
	if s.Authenticate == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return auth.AuthenticatedSession{}, false
	}
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return auth.AuthenticatedSession{}, false
	}
	authenticated, err := s.Authenticate.Execute(r.Context(), cookie.Value)
	if err != nil {
		writeAuthError(w, err)
		return auth.AuthenticatedSession{}, false
	}
	return authenticated, true
}

func (s Server) requireProductClient(w http.ResponseWriter, r *http.Request) (productclients.ProductClient, bool) {
	if s.ProductClients == nil {
		writeError(w, http.StatusUnauthorized, "product_client_required")
		return productclients.ProductClient{}, false
	}
	productKey := strings.TrimSpace(r.Header.Get("X-Product-Key"))
	token := bearerToken(r.Header.Get("Authorization"))
	if token == "" {
		token = strings.TrimSpace(r.Header.Get("X-Product-Client-Token"))
	}
	client, err := s.ProductClients.Authenticate(r.Context(), productKey, token)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "product_client_invalid")
		return productclients.ProductClient{}, false
	}
	return client, true
}

func (s Server) selectedTenant(w http.ResponseWriter, r *http.Request) (string, bool) {
	if cookie, err := r.Cookie(selectedTenantCookieName); err == nil && strings.TrimSpace(cookie.Value) != "" {
		return strings.TrimSpace(cookie.Value), true
	}
	if header := strings.TrimSpace(r.Header.Get("X-Tenant-ID")); header != "" {
		return header, true
	}
	writeError(w, http.StatusBadRequest, "tenant_not_selected")
	return "", false
}

func writeDecision(w http.ResponseWriter, decision access.AccessDecision, err error) {
	if err != nil {
		writeAccessError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"allowed": decision.Allowed,
		"reason":  decision.Reason,
	})
}

func (s Server) setCookie(w http.ResponseWriter, name, value string, expires time.Time, httpOnly bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   s.Config.CookieDomain,
		Expires:  expires,
		HttpOnly: httpOnly,
		Secure:   s.Config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (s Server) clearCookie(w http.ResponseWriter, name string, httpOnly bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Domain:   s.Config.CookieDomain,
		MaxAge:   -1,
		HttpOnly: httpOnly,
		Secure:   s.Config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func requireCSRF(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie(csrfCookieName)
	if err != nil || cookie.Value == "" {
		writeError(w, http.StatusForbidden, "csrf_token_missing")
		return false
	}
	header := r.Header.Get("X-CSRF-Token")
	if subtle.ConstantTimeCompare([]byte(header), []byte(cookie.Value)) != 1 {
		writeError(w, http.StatusForbidden, "csrf_token_invalid")
		return false
	}
	return true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return false
	}
	return true
}

func writeAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, auth.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, "invalid_credentials")
	case errors.Is(err, auth.ErrInactiveUser):
		writeError(w, http.StatusForbidden, "account_inactive")
	case errors.Is(err, auth.ErrTokenRevoked), errors.Is(err, auth.ErrSessionExpired):
		writeError(w, http.StatusUnauthorized, "unauthorized")
	default:
		writeError(w, http.StatusUnauthorized, "unauthorized")
	}
}

func writeAccessError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, access.ErrTenantMembership):
		writeError(w, http.StatusForbidden, "tenant_membership_required")
	case errors.Is(err, access.ErrProductAccessDenied):
		writeError(w, http.StatusForbidden, "product_access_required")
	case errors.Is(err, access.ErrFeatureAccessDenied):
		writeError(w, http.StatusForbidden, "feature_access_required")
	case errors.Is(err, access.ErrPermissionDenied):
		writeError(w, http.StatusForbidden, "permission_required")
	case errors.Is(err, access.ErrSessionInactive), errors.Is(err, access.ErrStaleTokenVersion), errors.Is(err, access.ErrStaleEntitlements):
		writeError(w, http.StatusUnauthorized, "stale_authorization")
	default:
		writeError(w, http.StatusInternalServerError, "access_check_failed")
	}
}

func writeAdminResult(w http.ResponseWriter, result any, err error) {
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func writeAdminError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, admin.ErrSuperadminRequired):
		writeError(w, http.StatusForbidden, "superadmin_required")
	case errors.Is(err, admin.ErrNotFound):
		writeError(w, http.StatusNotFound, "not_found")
	default:
		writeError(w, http.StatusInternalServerError, "admin_request_failed")
	}
}

func writeError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, map[string]string{"error": code})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-content-type-options", "nosniff")
		w.Header().Set("referrer-policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	return r.RemoteAddr
}

func bearerToken(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}
	const prefix = "Bearer "
	if strings.HasPrefix(header, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(header, prefix))
	}
	return ""
}

func queryInt(r *http.Request, key string) int {
	value := strings.TrimSpace(r.URL.Query().Get(key))
	if value == "" {
		return 0
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return parsed
}
