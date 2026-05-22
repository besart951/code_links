package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type server struct {
	config      config
	store       Store
	signer      *tokenSigner
	emailSender EmailSender
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
	return &server{config: config, store: store, signer: signer, emailSender: NoopEmailSender{}}
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /.well-known/jwks.json", s.handleJWKS)
	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("POST /api/auth/login", s.handleLogin)
	mux.HandleFunc("POST /api/auth/signup", s.handleSignup)
	mux.HandleFunc("POST /api/auth/forgot-password", s.handleForgotPassword)
	mux.HandleFunc("POST /api/auth/reset-password", s.handleResetPassword)
	mux.HandleFunc("GET /api/auth/verify-email", s.handleVerifyEmail)
	mux.HandleFunc("POST /api/auth/refresh", s.handleRefresh)
	mux.HandleFunc("POST /api/auth/logout", s.handleLogout)
	mux.HandleFunc("POST /api/licenses/mock-purchase", s.handleMockPurchase)
	mux.HandleFunc("GET /api/admin/me", s.handleAdminMe)
	mux.HandleFunc("GET /api/admin/dashboard/stats", s.handleAdminDashboardStats)
	mux.HandleFunc("GET /api/admin/dashboard-summary", s.handleAdminDashboardStats)
	mux.HandleFunc("GET /api/admin/users", s.handleAdminUsers)
	mux.HandleFunc("GET /api/admin/users/{id}", s.handleAdminUserDetail)
	mux.HandleFunc("PATCH /api/admin/users/{id}", s.handleAdminUpdateUser)
	mux.HandleFunc("PATCH /api/admin/users/{id}/status", s.handleAdminUpdateUser)
	mux.HandleFunc("POST /api/admin/users/{id}/lock", s.handleAdminLockUser)
	mux.HandleFunc("POST /api/admin/users/{id}/unlock", s.handleAdminUnlockUser)
	mux.HandleFunc("PATCH /api/admin/users/{id}/role", s.handleAdminUpdateUserRole)
	mux.HandleFunc("GET /api/admin/login-attempts", s.handleAdminLoginAttempts)
	mux.HandleFunc("GET /api/admin/security-events", s.handleAdminSecurityEvents)
	mux.HandleFunc("GET /api/admin/notifications", s.handleAdminNotifications)
	mux.HandleFunc("GET /api/admin/settings/smtp", s.handleAdminGetSMTPSettings)
	mux.HandleFunc("PUT /api/admin/settings/smtp", s.handleAdminUpdateSMTPSettings)
	mux.HandleFunc("POST /api/admin/settings/smtp/test-email", s.handleAdminSendTestEmail)

	return s.withCORS(mux)
}

func (s *server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) handleJWKS(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.signer.jwks())
}

func (s *server) handleSignup(w http.ResponseWriter, r *http.Request) {
	var request signupRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	if _, err := mail.ParseAddress(request.Email); err != nil {
		writeError(w, http.StatusBadRequest, "invalid email")
		return
	}
	if len(request.Password) < 10 {
		writeError(w, http.StatusBadRequest, "password must be at least 10 characters")
		return
	}
	if !request.AcceptedTerms && !request.AcceptedTermsAndPrivacy {
		writeError(w, http.StatusBadRequest, "terms must be accepted")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	user, err := s.store.CreateUser(r.Context(), request.Name, request.Email, string(passwordHash))
	if errors.Is(err, errConflict) {
		writeError(w, http.StatusConflict, "email already exists")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	token, err := newOpaqueToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create verification token")
		return
	}
	if err := s.store.CreateEmailVerificationToken(r.Context(), user.ID, hashOpaqueToken(token), time.Now().UTC().Add(24*time.Hour)); err != nil {
		writeError(w, http.StatusInternalServerError, "could not create verification token")
		return
	}

	verificationURL := s.config.PublicFrontendURL + "/verify-email?token=" + token
	_ = s.store.CreateNotification(r.Context(), Notification{
		ID:        uuid.New(),
		Type:      "email_verification",
		Channel:   "email",
		Recipient: user.Email,
		Subject:   "CodeLinks E-Mail bestätigen",
		Status:    "queued",
		CreatedAt: time.Now().UTC(),
	})

	response := map[string]string{
		"message":         "Account created. Please verify your email before logging in.",
		"verificationUrl": verificationURL,
	}
	if s.config.Environment != "production" {
		response["debugVerificationToken"] = token
	}
	writeJSON(w, http.StatusCreated, response)
}

func (s *server) handleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		writeError(w, http.StatusBadRequest, "missing token")
		return
	}

	user, err := s.store.VerifyEmailToken(r.Context(), hashOpaqueToken(token), time.Now().UTC())
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid or expired token")
		return
	}

	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]string{"email": user.Email, "status": "verified"})
		return
	}
	http.Redirect(w, r, s.config.PublicFrontendURL+"/login?verified=1", http.StatusSeeOther)
}

func (s *server) handleForgotPassword(w http.ResponseWriter, r *http.Request) {
	var request forgotPasswordRequest
	if !decodeJSON(w, r, &request) {
		return
	}

	response := map[string]string{"message": "Falls ein Konto mit dieser E-Mail existiert, wurde ein Reset-Link versendet."}
	user, _, err := s.store.FindUserByEmail(r.Context(), request.Email)
	if err == nil {
		token, tokenErr := newOpaqueToken()
		if tokenErr != nil {
			writeError(w, http.StatusInternalServerError, "could not create reset token")
			return
		}
		if tokenErr := s.store.CreatePasswordResetToken(r.Context(), user.ID, hashOpaqueToken(token), time.Now().UTC().Add(time.Hour)); tokenErr != nil {
			writeError(w, http.StatusInternalServerError, "could not create reset token")
			return
		}
		resetURL := s.config.PublicFrontendURL + "/reset-password?token=" + token
		_ = s.store.CreateNotification(r.Context(), Notification{
			ID:        uuid.New(),
			Type:      "password_reset",
			Channel:   "email",
			Recipient: user.Email,
			Subject:   "CodeLinks Passwort zurücksetzen",
			Status:    "queued",
			CreatedAt: time.Now().UTC(),
		})
		if s.config.Environment != "production" {
			response["debugResetToken"] = token
			response["debugResetUrl"] = resetURL
		}
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *server) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	var request resetPasswordRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	if request.Token == "" {
		writeError(w, http.StatusBadRequest, "missing token")
		return
	}
	if len(request.Password) < 10 {
		writeError(w, http.StatusBadRequest, "password must be at least 10 characters")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	if _, err := s.store.ResetPasswordByToken(r.Context(), hashOpaqueToken(request.Token), string(passwordHash), time.Now().UTC()); err != nil {
		writeError(w, http.StatusBadRequest, "invalid or expired token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "password updated"})
}

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var request loginRequest
	if !decodeJSON(w, r, &request) {
		return
	}

	email := normalizeEmail(request.Email)
	user, licenses, err := s.store.FindUserByEmail(r.Context(), email)
	if errors.Is(err, errNotFound) {
		s.recordLoginAttempt(r, nil, email, false, LoginFailureUnknownEmail)
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load user")
		return
	}

	if user.Status == UserStatusLocked || user.Status == UserStatusDisabled {
		s.recordLoginAttempt(r, &user.ID, email, false, LoginFailureAccountLocked)
		writeError(w, http.StatusForbidden, "account is locked")
		return
	}
	if user.EmailVerifiedAt == nil {
		s.recordLoginAttempt(r, &user.ID, email, false, LoginFailureEmailNotConfirmed)
		writeError(w, http.StatusForbidden, "email_not_confirmed")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)) != nil {
		s.recordLoginAttempt(r, &user.ID, email, false, LoginFailureWrongPassword)
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	s.recordLoginAttempt(r, &user.ID, email, true, "")
	s.writeSession(w, r, user, licenses, true)
}

func (s *server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		writeError(w, http.StatusUnauthorized, "missing refresh token")
		return
	}

	tokenHash := hashRefreshToken(cookie.Value)
	userID, err := s.store.FindRefreshSession(r.Context(), tokenHash)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	_ = s.store.DeleteRefreshSession(r.Context(), tokenHash)

	user, licenses, err := s.store.FindUserByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	s.writeSession(w, r, user, licenses, true)
}

func (s *server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("refresh_token"); err == nil {
		_ = s.store.DeleteRefreshSession(r.Context(), hashRefreshToken(cookie.Value))
	}

	http.SetCookie(w, s.refreshCookie("", time.Unix(0, 0)))
	w.WriteHeader(http.StatusNoContent)
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

	licenses, err := s.store.GrantLicense(r.Context(), userID, request.ProductID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not grant license")
		return
	}

	user, _, err := s.store.FindUserByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load user")
		return
	}

	s.writeSession(w, r, user, licenses, false)
}

func (s *server) handleAdminMe(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, PermissionDashboardRead)
	if !ok {
		return
	}

	writeJSON(w, http.StatusOK, actor)
}

func (s *server) handleAdminDashboardStats(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r, PermissionDashboardRead); !ok {
		return
	}

	stats, err := s.store.GetDashboardStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load dashboard stats")
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

func (s *server) handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r, PermissionUsersRead); !ok {
		return
	}

	result, err := s.store.ListAdminUsers(r.Context(), AdminUserListQuery{
		Query:     r.URL.Query().Get("query"),
		Role:      r.URL.Query().Get("role"),
		Status:    UserStatus(r.URL.Query().Get("status")),
		Page:      queryInt(r, "page", 1),
		PageSize:  queryInt(r, "pageSize", 20),
		Sort:      r.URL.Query().Get("sort"),
		Direction: r.URL.Query().Get("direction"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load users")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *server) handleAdminUserDetail(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r, PermissionUsersRead); !ok {
		return
	}

	userID, ok := pathUUID(w, r, "id")
	if !ok {
		return
	}
	user, err := s.store.GetManagedUserDetail(r.Context(), userID)
	if errors.Is(err, errNotFound) {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load user")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (s *server) handleAdminUpdateUser(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, PermissionUsersUpdate)
	if !ok {
		return
	}

	userID, ok := pathUUID(w, r, "id")
	if !ok {
		return
	}
	var request updateUserRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	user, err := s.store.SetUserStatus(r.Context(), userID, request.Status)
	if errors.Is(err, errNotFound) {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not update user")
		return
	}
	s.audit(r, actor, "admin.users.update_status", "user", userID.String(), string(request.Status))
	writeJSON(w, http.StatusOK, user)
}

func (s *server) handleAdminLockUser(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, PermissionUsersLock)
	if !ok {
		return
	}
	userID, ok := pathUUID(w, r, "id")
	if !ok {
		return
	}
	user, err := s.store.SetUserStatus(r.Context(), userID, UserStatusLocked)
	if errors.Is(err, errNotFound) {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not lock user")
		return
	}
	s.audit(r, actor, "admin.users.lock", "user", userID.String(), "")
	writeJSON(w, http.StatusOK, user)
}

func (s *server) handleAdminUnlockUser(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, PermissionUsersLock)
	if !ok {
		return
	}
	userID, ok := pathUUID(w, r, "id")
	if !ok {
		return
	}
	user, err := s.store.SetUserStatus(r.Context(), userID, UserStatusActive)
	if errors.Is(err, errNotFound) {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not unlock user")
		return
	}
	s.audit(r, actor, "admin.users.unlock", "user", userID.String(), "")
	writeJSON(w, http.StatusOK, user)
}

func (s *server) handleAdminUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, PermissionUsersChangeRole)
	if !ok {
		return
	}
	userID, ok := pathUUID(w, r, "id")
	if !ok {
		return
	}
	var request updateRoleRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	user, err := s.store.SetUserRole(r.Context(), userID, request.Role)
	if errors.Is(err, errNotFound) {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not change role")
		return
	}
	s.audit(r, actor, "admin.users.change_role", "user", userID.String(), string(request.Role))
	writeJSON(w, http.StatusOK, user)
}

func (s *server) handleAdminLoginAttempts(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, PermissionAuthLogsRead)
	if !ok {
		return
	}

	var userID *uuid.UUID
	if rawUserID := r.URL.Query().Get("userId"); rawUserID != "" {
		parsed, err := uuid.Parse(rawUserID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid userId")
			return
		}
		userID = &parsed
	}

	var success *bool
	if rawSuccess := r.URL.Query().Get("success"); rawSuccess != "" {
		parsed, err := strconv.ParseBool(rawSuccess)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid success")
			return
		}
		success = &parsed
	}

	result, err := s.store.ListLoginAttempts(r.Context(), LoginAttemptListQuery{
		UserID:   userID,
		Success:  success,
		Query:    r.URL.Query().Get("query"),
		Page:     queryInt(r, "page", 1),
		PageSize: queryInt(r, "pageSize", 20),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load login attempts")
		return
	}
	if !hasAdminPermission(actor, PermissionUsersUpdate) {
		for index := range result.Items {
			result.Items[index].IPAddress = maskIP(result.Items[index].IPAddress)
		}
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *server) handleAdminSecurityEvents(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, PermissionSecurityEventsRead)
	if !ok {
		return
	}

	events, err := s.store.ListSecurityEvents(r.Context(), queryInt(r, "limit", 50))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load security events")
		return
	}
	if !hasAdminPermission(actor, PermissionUsersUpdate) {
		for index := range events {
			events[index].SourceIPAddress = maskIP(events[index].SourceIPAddress)
		}
	}

	writeJSON(w, http.StatusOK, events)
}

func (s *server) handleAdminNotifications(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r, PermissionNotificationsRead); !ok {
		return
	}

	notifications, err := s.store.ListNotifications(r.Context(), queryInt(r, "limit", 50))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load notifications")
		return
	}

	writeJSON(w, http.StatusOK, notifications)
}

func (s *server) handleAdminGetSMTPSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r, PermissionSMTPSettingsRead); !ok {
		return
	}

	settings, err := s.store.GetSmtpSettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load smtp settings")
		return
	}

	writeJSON(w, http.StatusOK, settings)
}

func (s *server) handleAdminUpdateSMTPSettings(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, PermissionSMTPSettingsUpdate)
	if !ok {
		return
	}

	var request smtpSettingsRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	encryptedPassword := ""
	if request.Password != "" {
		var err error
		encryptedPassword, err = encryptSecret(request.Password, s.config.SMTPSecretKey)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not encrypt smtp password")
			return
		}
	}

	settings, err := s.store.SaveSmtpSettings(r.Context(), SmtpSettings{
		Host:              strings.TrimSpace(request.Host),
		Port:              request.Port,
		Username:          strings.TrimSpace(request.Username),
		PasswordEncrypted: encryptedPassword,
		Encryption:        request.Encryption,
		FromEmail:         strings.TrimSpace(request.FromEmail),
		FromName:          strings.TrimSpace(request.FromName),
		ReplyToEmail:      strings.TrimSpace(request.ReplyToEmail),
		Active:            request.Active,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not save smtp settings")
		return
	}
	s.audit(r, actor, "admin.smtp_settings.update", "smtp_settings", "default", "")
	writeJSON(w, http.StatusOK, settings)
}

func (s *server) handleAdminSendTestEmail(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, PermissionSMTPSettingsUpdate)
	if !ok {
		return
	}

	var request testEmailRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	if _, err := mail.ParseAddress(request.Recipient); err != nil {
		writeError(w, http.StatusBadRequest, "invalid recipient")
		return
	}

	settings, err := s.store.GetSmtpSettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load smtp settings")
		return
	}
	password := ""
	if settings.PasswordEncrypted != "" {
		password, err = decryptSecret(settings.PasswordEncrypted, s.config.SMTPSecretKey)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not decrypt smtp password")
			return
		}
	}
	if err := s.emailSender.Send(r.Context(), settings, EmailMessage{
		To:      request.Recipient,
		Subject: "CodeLinks SMTP Test",
		Body:    "Diese Test-E-Mail wurde von CodeLinks gesendet.",
	}, password); err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	now := time.Now().UTC()
	_ = s.store.CreateNotification(r.Context(), Notification{
		ID:        uuid.New(),
		Type:      "smtp_test",
		Channel:   "email",
		Recipient: request.Recipient,
		Subject:   "CodeLinks SMTP Test",
		Status:    "sent",
		CreatedAt: now,
		SentAt:    &now,
	})
	s.audit(r, actor, "admin.smtp_settings.test_email", "smtp_settings", "default", request.Recipient)
	writeJSON(w, http.StatusOK, map[string]string{"message": "test email sent"})
}

func (s *server) writeSession(w http.ResponseWriter, r *http.Request, user User, licenses []string, rotateRefresh bool) {
	if rotateRefresh {
		refreshToken, expiresAt, err := s.createRefreshSession(r, user.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not create refresh session")
			return
		}
		http.SetCookie(w, s.refreshCookie(refreshToken, expiresAt))
	}

	roles, permissions, err := s.store.ListUserRoles(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load roles")
		return
	}
	accessToken, accessExpiresAt, err := s.signer.issue(user, licenses, roles, permissions)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not issue token")
		return
	}

	writeJSON(w, http.StatusOK, sessionResponse{
		AccessToken: accessToken,
		ExpiresAt:   accessExpiresAt.Format(time.RFC3339),
		User: userResponse{
			ID:            user.ID.String(),
			Email:         user.Email,
			Name:          user.Name,
			Status:        user.Status,
			EmailVerified: user.EmailVerifiedAt != nil,
			Licenses:      licenses,
			Roles:         roles,
			Permissions:   permissions,
		},
	})
}

func (s *server) createRefreshSession(r *http.Request, userID uuid.UUID) (string, time.Time, error) {
	token, err := newOpaqueToken()
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := time.Now().UTC().Add(s.config.RefreshTokenLifetime)
	if err := s.store.CreateRefreshSession(r.Context(), hashRefreshToken(token), userID, expiresAt); err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func (s *server) refreshCookie(value string, expiresAt time.Time) *http.Cookie {
	maxAge := int(time.Until(expiresAt).Seconds())
	if value == "" {
		maxAge = -1
	}

	return &http.Cookie{
		Name:     "refresh_token",
		Value:    value,
		Path:     "/",
		Domain:   s.config.CookieDomain,
		Expires:  expiresAt,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   s.config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
}

func (s *server) claimsFromBearer(r *http.Request) (accessClaims, error) {
	header := r.Header.Get("Authorization")
	token, ok := strings.CutPrefix(header, "Bearer ")
	if !ok || token == "" {
		return accessClaims{}, errors.New("missing bearer token")
	}

	return s.signer.parse(token)
}

func (s *server) requireAdmin(w http.ResponseWriter, r *http.Request, permission AdminPermission) (AdminActor, bool) {
	actor, err := s.adminActorFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "admin authentication required")
		return AdminActor{}, false
	}
	if !hasAdminPermission(actor, permission) {
		writeError(w, http.StatusForbidden, "missing permission")
		return AdminActor{}, false
	}

	return actor, true
}

func (s *server) adminActorFromRequest(r *http.Request) (AdminActor, error) {
	if header := r.Header.Get("Authorization"); strings.HasPrefix(header, "Bearer ") {
		claims, err := s.claimsFromBearer(r)
		if err != nil {
			return AdminActor{}, err
		}
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return AdminActor{}, err
		}

		return s.store.GetAdminActor(r.Context(), userID)
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		return AdminActor{}, errNotFound
	}
	userID, err := s.store.FindRefreshSession(r.Context(), hashRefreshToken(cookie.Value))
	if err != nil {
		return AdminActor{}, err
	}

	return s.store.GetAdminActor(r.Context(), userID)
}

func (s *server) recordLoginAttempt(r *http.Request, userID *uuid.UUID, email string, success bool, failureReason LoginFailureReason) {
	var reason *LoginFailureReason
	if failureReason != "" {
		reason = &failureReason
	}
	browser, operatingSystem := deviceFromUserAgent(r.UserAgent())
	ipAddress := clientIPAddress(r)
	ipHash := sha256.Sum256([]byte(ipAddress))

	_ = s.store.RecordLoginAttempt(r.Context(), LoginAttempt{
		ID:              uuid.New(),
		UserID:          userID,
		EmailAttempted:  email,
		OccurredAt:      time.Now().UTC(),
		IPAddress:       ipAddress,
		IPHash:          base64.RawURLEncoding.EncodeToString(ipHash[:]),
		CountryCode:     countryFromRequest(r),
		City:            cityFromRequest(r),
		UserAgent:       r.UserAgent(),
		Browser:         browser,
		OperatingSystem: operatingSystem,
		Success:         success,
		FailureReason:   reason,
		AuthMethod:      "password",
		RiskScore:       riskScore(success, failureReason),
		CorrelationID:   r.Header.Get("X-Request-Id"),
	})
}

func (s *server) audit(r *http.Request, actor AdminActor, action string, targetType string, targetID string, reason string) {
	actorID, err := uuid.Parse(actor.ID)
	if err != nil {
		return
	}
	_ = s.store.RecordAdminAuditEntry(r.Context(), AdminAuditEntry{
		ID:          uuid.New(),
		ActorUserID: actorID,
		Action:      action,
		TargetType:  targetType,
		TargetID:    targetID,
		Reason:      reason,
		IPAddress:   clientIPAddress(r),
		CreatedAt:   time.Now().UTC(),
	})
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

func hasAdminPermission(actor AdminActor, permission AdminPermission) bool {
	for _, grant := range actor.Permissions {
		if grant == permission {
			return true
		}
	}

	return false
}

func riskScore(success bool, failureReason LoginFailureReason) int {
	if success {
		return 8
	}
	switch failureReason {
	case LoginFailureAccountLocked:
		return 85
	case LoginFailureUnknownEmail:
		return 55
	case LoginFailureEmailNotConfirmed:
		return 30
	default:
		return 65
	}
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
	switch productID {
	case "infra-link", "planer-link", "loka-link":
		return true
	default:
		return false
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
