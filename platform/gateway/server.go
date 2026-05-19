package gateway

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/besart951/code_links/platform/auth"
	"github.com/besart951/code_links/platform/billing"
	"github.com/besart951/code_links/platform/entitlements"
	"github.com/besart951/code_links/platform/tenants"
)

type tenantStore interface {
	ListTenantsForUser(ctx context.Context, userID string) ([]tenants.Tenant, error)
	IsTenantMember(ctx context.Context, userID, tenantID string) (bool, error)
	ListEntitlements(ctx context.Context, tenantID string) ([]entitlements.Entitlement, error)
	ListFeatureLimits(ctx context.Context, tenantID string) ([]billing.FeatureLimit, error)
	HasEntitlement(ctx context.Context, tenantID, productKey, featureKey string) (bool, error)
}

type Server struct {
	cfg      Config
	auth     *auth.Service
	store    tenantStore
	authRepo auth.Repository
}

func NewServer(cfg Config, authService *auth.Service, store tenantStore, authRepo auth.Repository) http.Handler {
	server := &Server{cfg: cfg, auth: authService, store: store, authRepo: authRepo}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", server.handleHealth)
	mux.HandleFunc("POST /api/v1/auth/login", server.handleLogin)
	mux.HandleFunc("POST /api/v1/auth/refresh", server.handleRefresh)
	mux.HandleFunc("POST /api/v1/auth/logout", server.handleLogout)
	mux.HandleFunc("GET /api/v1/me", server.handleMe)
	mux.HandleFunc("GET /api/v1/tenants", server.handleTenants)
	mux.HandleFunc("GET /api/v1/entitlements", server.handleEntitlements)
	mux.HandleFunc("POST /internal/authorize", server.handleAuthorize)
	return withSecurityHeaders(mux)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	User                 auth.User `json:"user"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	RefreshExpiresAt     time.Time `json:"refresh_token_expires_at"`
	CSRFToken            string    `json:"csrf_token"`
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	userAgent := r.UserAgent()
	ip := clientIP(r)
	pair, err := s.auth.Login(r.Context(), req.Email, req.Password, &userAgent, &ip)
	if err != nil {
		writeAuthError(w, err)
		return
	}
	s.setAuthCookies(w, pair)
	writeJSON(w, http.StatusOK, authResponse{
		User:                 pair.User,
		AccessTokenExpiresAt: pair.AccessExpiresAt,
		RefreshExpiresAt:     pair.RefreshExpiresAt,
		CSRFToken:            pair.CSRFToken,
	})
}

func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	if !s.requireCSRF(w, r) {
		return
	}
	cookie, err := r.Cookie("refresh_token")
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	userAgent := r.UserAgent()
	ip := clientIP(r)
	pair, err := s.auth.Refresh(r.Context(), cookie.Value, &userAgent, &ip)
	if err != nil {
		writeAuthError(w, err)
		return
	}
	s.setAuthCookies(w, pair)
	writeJSON(w, http.StatusOK, authResponse{
		User:                 pair.User,
		AccessTokenExpiresAt: pair.AccessExpiresAt,
		RefreshExpiresAt:     pair.RefreshExpiresAt,
		CSRFToken:            pair.CSRFToken,
	})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if !s.requireCSRF(w, r) {
		return
	}
	if cookie, err := r.Cookie("refresh_token"); err == nil {
		_ = s.auth.Logout(r.Context(), cookie.Value)
	}
	s.clearAuthCookies(w)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	user, ok := s.requireUser(w, r)
	if !ok {
		return
	}
	tenants, err := s.store.ListTenantsForUser(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "tenant_lookup_failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": user, "tenants": tenants})
}

func (s *Server) handleTenants(w http.ResponseWriter, r *http.Request) {
	user, ok := s.requireUser(w, r)
	if !ok {
		return
	}
	tenants, err := s.store.ListTenantsForUser(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "tenant_lookup_failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tenants": tenants})
}

func (s *Server) handleEntitlements(w http.ResponseWriter, r *http.Request) {
	user, ok := s.requireUser(w, r)
	if !ok {
		return
	}
	tenantID := strings.TrimSpace(r.URL.Query().Get("tenant_id"))
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id_required")
		return
	}
	member, err := s.store.IsTenantMember(r.Context(), user.ID, tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "tenant_membership_check_failed")
		return
	}
	if !member {
		writeError(w, http.StatusForbidden, "tenant_membership_required")
		return
	}
	items, err := s.store.ListEntitlements(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "entitlement_lookup_failed")
		return
	}
	limits, err := s.store.ListFeatureLimits(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "feature_limit_lookup_failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"tenant_id":    tenantID,
		"entitlements": items,
		"limits":       limits,
	})
}

func (s *Server) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	if !s.requireInternalToken(w, r) {
		return
	}
	var req entitlements.AuthorizeRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	member, err := s.store.IsTenantMember(r.Context(), req.UserID, req.TenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "tenant_membership_check_failed")
		return
	}
	if !member {
		writeJSON(w, http.StatusOK, entitlements.AuthorizeResponse{
			Allowed: false,
			Reason:  strPtr("tenant_membership_required"),
		})
		return
	}
	allowed, err := s.store.HasEntitlement(r.Context(), req.TenantID, req.ProductKey, req.FeatureKey)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "entitlement_lookup_failed")
		return
	}
	if !allowed {
		writeJSON(w, http.StatusOK, entitlements.AuthorizeResponse{
			Allowed: false,
			Reason:  strPtr("entitlement_required"),
		})
		return
	}
	writeJSON(w, http.StatusOK, entitlements.AuthorizeResponse{Allowed: true})
}

func (s *Server) requireUser(w http.ResponseWriter, r *http.Request) (auth.User, bool) {
	cookie, err := r.Cookie("access_token")
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return auth.User{}, false
	}
	claims, err := s.auth.ValidateAccessToken(cookie.Value)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return auth.User{}, false
	}
	user, err := s.authRepo.FindUserByID(r.Context(), claims.Subject)
	if err != nil || user.Status != "active" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return auth.User{}, false
	}
	return user, true
}

func (s *Server) requireCSRF(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("csrf_token")
	if err != nil || cookie.Value == "" {
		writeError(w, http.StatusForbidden, "csrf_token_missing")
		return false
	}
	header := r.Header.Get("X-CSRF-Token")
	if header == "" || subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(header)) != 1 {
		writeError(w, http.StatusForbidden, "csrf_token_invalid")
		return false
	}
	return true
}

func (s *Server) requireInternalToken(w http.ResponseWriter, r *http.Request) bool {
	if s.cfg.InternalToken == "" {
		return true
	}
	header := r.Header.Get("X-Internal-Token")
	if subtle.ConstantTimeCompare([]byte(header), []byte(s.cfg.InternalToken)) != 1 {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return false
	}
	return true
}

func (s *Server) setAuthCookies(w http.ResponseWriter, pair auth.TokenPair) {
	s.setCookie(w, "access_token", pair.AccessToken, int(time.Until(pair.AccessExpiresAt).Seconds()), true)
	s.setCookie(w, "refresh_token", pair.RefreshToken, int(time.Until(pair.RefreshExpiresAt).Seconds()), true)
	s.setCookie(w, "csrf_token", pair.CSRFToken, int((24 * time.Hour).Seconds()), false)
}

func (s *Server) clearAuthCookies(w http.ResponseWriter) {
	s.setCookie(w, "access_token", "", -1, true)
	s.setCookie(w, "refresh_token", "", -1, true)
	s.setCookie(w, "csrf_token", "", -1, false)
}

func (s *Server) setCookie(w http.ResponseWriter, name, value string, maxAge int, httpOnly bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   s.cfg.CookieDomain,
		MaxAge:   maxAge,
		HttpOnly: httpOnly,
		Secure:   s.cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
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
	case errors.Is(err, auth.ErrTokenRevoked):
		writeError(w, http.StatusUnauthorized, "token_revoked")
	default:
		writeError(w, http.StatusUnauthorized, "unauthorized")
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

func withSecurityHeaders(next http.Handler) http.Handler {
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

func strPtr(value string) *string {
	return &value
}
