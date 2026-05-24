package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *server) registerAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/login", s.handleLogin)
	mux.HandleFunc("POST /api/auth/signup", s.handleSignup)
	mux.HandleFunc("POST /api/auth/forgot-password", s.handleForgotPassword)
	mux.HandleFunc("POST /api/auth/reset-password", s.handleResetPassword)
	mux.HandleFunc("GET /api/auth/verify-email", s.handleVerifyEmail)
	mux.HandleFunc("POST /api/auth/refresh", s.handleRefresh)
	mux.HandleFunc("POST /api/auth/logout", s.handleLogout)
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

	user, err := s.authStore.CreateUser(r.Context(), request.Name, request.Email, string(passwordHash))
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
	if err := s.authStore.CreateEmailVerificationToken(r.Context(), user.ID, hashOpaqueToken(token), time.Now().UTC().Add(24*time.Hour)); err != nil {
		writeError(w, http.StatusInternalServerError, "could not create verification token")
		return
	}

	verificationURL := s.config.PublicFrontendURL + "/verify-email?token=" + token
	_ = s.notificationStore.CreateNotification(r.Context(), Notification{
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

	user, err := s.authStore.VerifyEmailToken(r.Context(), hashOpaqueToken(token), time.Now().UTC())
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
	user, _, err := s.authStore.FindUserByEmail(r.Context(), request.Email)
	if err == nil {
		token, tokenErr := newOpaqueToken()
		if tokenErr != nil {
			writeError(w, http.StatusInternalServerError, "could not create reset token")
			return
		}
		if tokenErr := s.authStore.CreatePasswordResetToken(r.Context(), user.ID, hashOpaqueToken(token), time.Now().UTC().Add(time.Hour)); tokenErr != nil {
			writeError(w, http.StatusInternalServerError, "could not create reset token")
			return
		}
		resetURL := s.config.PublicFrontendURL + "/reset-password?token=" + token
		_ = s.notificationStore.CreateNotification(r.Context(), Notification{
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

	if _, err := s.authStore.ResetPasswordByToken(r.Context(), hashOpaqueToken(request.Token), string(passwordHash), time.Now().UTC()); err != nil {
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
	user, licenses, err := s.authStore.FindUserByEmail(r.Context(), email)
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
	userID, err := s.authStore.FindRefreshSession(r.Context(), tokenHash)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	_ = s.authStore.DeleteRefreshSession(r.Context(), tokenHash)

	user, licenses, err := s.authStore.FindUserByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	s.writeSession(w, r, user, licenses, true)
}

func (s *server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("refresh_token"); err == nil {
		_ = s.authStore.DeleteRefreshSession(r.Context(), hashRefreshToken(cookie.Value))
	}

	http.SetCookie(w, s.refreshCookie("", time.Unix(0, 0)))
	w.WriteHeader(http.StatusNoContent)
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

	roles, permissions, err := s.adminStore.ListUserRoles(r.Context(), user.ID)
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
	if err := s.authStore.CreateRefreshSession(r.Context(), hashRefreshToken(token), userID, expiresAt); err != nil {
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

func (s *server) recordLoginAttempt(r *http.Request, userID *uuid.UUID, email string, success bool, failureReason LoginFailureReason) {
	var reason *LoginFailureReason
	if failureReason != "" {
		reason = &failureReason
	}
	browser, operatingSystem := deviceFromUserAgent(r.UserAgent())
	ipAddress := clientIPAddress(r)
	ipHash := sha256.Sum256([]byte(ipAddress))

	_ = s.authStore.RecordLoginAttempt(r.Context(), LoginAttempt{
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
