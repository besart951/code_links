package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/besart951/code-links/apps/auth-service/backend/internal/admin"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/auth"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/domain"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/token"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/userinfo"
	"github.com/besart951/code-links/packages/productcatalog"
	"github.com/google/uuid"
)

type Config struct {
	AllowedOrigins     []string
	EnableMockPurchase bool
	CookieDomain       string
	CookieSecure       bool
	PublicFrontendURL  string
}

type Server struct {
	config Config
	auth   *auth.Service
	admin  *admin.Service
	users  *userinfo.Service
	tokens *token.Signer
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signupRequest struct {
	Name                    string `json:"name"`
	Email                   string `json:"email"`
	Password                string `json:"password"`
	AcceptedTerms           bool   `json:"acceptedTerms"`
	AcceptedTermsAndPrivacy bool   `json:"acceptedTermsAndPrivacy"`
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

type resetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type purchaseRequest struct {
	ProductID string `json:"productId"`
}

type lookupUsersRequest struct {
	UserIDs []string `json:"userIds"`
}

type updateUserRequest struct {
	Status domain.UserStatus `json:"status"`
}

type updateRoleRequest struct {
	Role domain.AdminRole `json:"role"`
}

type smtpSettingsRequest struct {
	Host         string                `json:"host"`
	Port         int                   `json:"port"`
	Username     string                `json:"username"`
	Password     string                `json:"password"`
	Encryption   domain.SmtpEncryption `json:"encryption"`
	FromEmail    string                `json:"fromEmail"`
	FromName     string                `json:"fromName"`
	ReplyToEmail string                `json:"replyToEmail"`
	Active       bool                  `json:"active"`
}

type testEmailRequest struct {
	Recipient string `json:"recipient"`
}

type sessionResponse struct {
	AccessToken string       `json:"accessToken"`
	User        userResponse `json:"user"`
	ExpiresAt   string       `json:"expiresAt"`
}

type userResponse struct {
	ID            string                   `json:"id"`
	Email         string                   `json:"email"`
	Name          string                   `json:"name"`
	Status        domain.UserStatus        `json:"status"`
	EmailVerified bool                     `json:"emailVerified"`
	Licenses      []string                 `json:"licenses"`
	Roles         []domain.AdminRole       `json:"roles"`
	Permissions   []domain.AdminPermission `json:"permissions"`
}

func NewServer(config Config, authService *auth.Service, adminService *admin.Service, userInfoService *userinfo.Service, tokens *token.Signer) *Server {
	return &Server{config: config, auth: authService, admin: adminService, users: userInfoService, tokens: tokens}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	s.registerSystemRoutes(mux)
	s.registerAuthRoutes(mux)
	s.registerLicenseRoutes(mux)
	s.registerUserInfoRoutes(mux)
	s.registerAdminRoutes(mux)

	return s.withCORS(mux)
}

func (s *Server) registerSystemRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /.well-known/jwks.json", s.handleJWKS)
	mux.HandleFunc("GET /api/health", s.handleHealth)
}

func (s *Server) registerLicenseRoutes(mux *http.ServeMux) {
	if s.config.EnableMockPurchase {
		mux.HandleFunc("POST /api/licenses/mock-purchase", s.handleMockPurchase)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleJWKS(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.tokens.JWKS())
}

func (s *Server) handleMockPurchase(w http.ResponseWriter, r *http.Request) {
	claims, err := s.claimsFromBearer(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid access token")
		return
	}

	var request purchaseRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	if !validProductID(request.ProductID) {
		writeError(w, http.StatusBadRequest, "unknown product")
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid access token")
		return
	}

	session, err := s.auth.MockPurchase(r.Context(), userID, request.ProductID)
	if err != nil {
		writeAuthError(w, err)
		return
	}
	s.writeSession(w, session)
}

func (s *Server) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if s.isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type, x-sveltekit-action, x-request-id")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) isAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	for _, allowed := range s.config.AllowedOrigins {
		if origin == allowed {
			return true
		}
	}

	return false
}

func queryInt(r *http.Request, key string, fallback int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func pathUUID(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, bool) {
	value := r.PathValue(key)
	parsed, err := uuid.Parse(value)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid "+key)
		return uuid.Nil, false
	}

	return parsed, true
}

func wantsJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "application/json")
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}

	return true
}

func validProductID(productID string) bool {
	return productcatalog.IsValidID(productID)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
