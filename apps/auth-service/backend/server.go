package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/besart951/code-links/packages/productcatalog"
	"github.com/google/uuid"
)

type server struct {
	config            config
	store             Store
	authStore         AuthStore
	adminStore        AdminStore
	smtpStore         SmtpSettingsStore
	notificationStore NotificationStore
	auditStore        AuditStore
	signer            *tokenSigner
	emailSender       EmailSender
	adminProjection   AdminProjectionPolicy
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

type updateUserRequest struct {
	Status UserStatus `json:"status"`
}

type updateRoleRequest struct {
	Role AdminRole `json:"role"`
}

type smtpSettingsRequest struct {
	Host         string         `json:"host"`
	Port         int            `json:"port"`
	Username     string         `json:"username"`
	Password     string         `json:"password"`
	Encryption   SmtpEncryption `json:"encryption"`
	FromEmail    string         `json:"fromEmail"`
	FromName     string         `json:"fromName"`
	ReplyToEmail string         `json:"replyToEmail"`
	Active       bool           `json:"active"`
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
	ID            string            `json:"id"`
	Email         string            `json:"email"`
	Name          string            `json:"name"`
	Status        UserStatus        `json:"status"`
	EmailVerified bool              `json:"emailVerified"`
	Licenses      []string          `json:"licenses"`
	Roles         []AdminRole       `json:"roles"`
	Permissions   []AdminPermission `json:"permissions"`
}

func newServer(config config, store Store, signer *tokenSigner) *server {
	return &server{
		config:            config,
		store:             store,
		authStore:         store,
		adminStore:        store,
		smtpStore:         store,
		notificationStore: store,
		auditStore:        store,
		signer:            signer,
		emailSender:       NoopEmailSender{},
		adminProjection:   AdminProjectionPolicy{},
	}
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	s.registerSystemRoutes(mux)
	s.registerAuthRoutes(mux)
	s.registerLicenseRoutes(mux)
	s.registerAdminRoutes(mux)

	return s.withCORS(mux)
}

func (s *server) registerSystemRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /.well-known/jwks.json", s.handleJWKS)
	mux.HandleFunc("GET /api/health", s.handleHealth)
}

func (s *server) registerLicenseRoutes(mux *http.ServeMux) {
	if s.config.EnableMockPurchase {
		mux.HandleFunc("POST /api/licenses/mock-purchase", s.handleMockPurchase)
	}
}

func (s *server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) handleJWKS(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.signer.jwks())
}

func (s *server) handleMockPurchase(w http.ResponseWriter, r *http.Request) {
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

	licenses, err := s.authStore.GrantLicense(r.Context(), userID, request.ProductID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not grant license")
		return
	}

	user, _, err := s.authStore.FindUserByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load user")
		return
	}

	s.writeSession(w, r, user, licenses, false)
}

func (s *server) withCORS(next http.Handler) http.Handler {
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

func (s *server) isAllowedOrigin(origin string) bool {
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

func maskIP(value string) string {
	parts := strings.Split(value, ".")
	if len(parts) == 4 {
		return parts[0] + "." + parts[1] + ".x.x"
	}
	if value == "" {
		return value
	}

	return "masked"
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
