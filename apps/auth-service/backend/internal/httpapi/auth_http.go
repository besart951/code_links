package httpapi

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/besart951/code-links/apps/auth-service/backend/internal/auth"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/requestmeta"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/token"
)

func (s *Server) registerAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/login", s.handleLogin)
	mux.HandleFunc("POST /api/auth/signup", s.handleSignup)
	mux.HandleFunc("POST /api/auth/forgot-password", s.handleForgotPassword)
	mux.HandleFunc("POST /api/auth/reset-password", s.handleResetPassword)
	mux.HandleFunc("GET /api/auth/verify-email", s.handleVerifyEmail)
	mux.HandleFunc("POST /api/auth/refresh", s.handleRefresh)
	mux.HandleFunc("POST /api/auth/logout", s.handleLogout)
}

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	var request signupRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := s.auth.Signup(r.Context(), auth.SignupInput{
		Name:                    request.Name,
		Email:                   request.Email,
		Password:                request.Password,
		AcceptedTerms:           request.AcceptedTerms,
		AcceptedTermsAndPrivacy: request.AcceptedTermsAndPrivacy,
	})
	if err != nil {
		writeAuthError(w, err)
		return
	}

	response := map[string]string{
		"message":         result.Message,
		"verificationUrl": result.VerificationURL,
	}
	if result.DebugVerificationToken != "" {
		response["debugVerificationToken"] = result.DebugVerificationToken
	}
	writeJSON(w, http.StatusCreated, response)
}

func (s *Server) handleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	user, err := s.auth.VerifyEmail(r.Context(), r.URL.Query().Get("token"))
	if err != nil {
		writeAuthError(w, err)
		return
	}

	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]string{"email": user.Email, "status": "verified"})
		return
	}
	http.Redirect(w, r, s.config.PublicFrontendURL+"/login?verified=1", http.StatusSeeOther)
}

func (s *Server) handleForgotPassword(w http.ResponseWriter, r *http.Request) {
	var request forgotPasswordRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := s.auth.ForgotPassword(r.Context(), request.Email)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	response := map[string]string{"message": result.Message}
	if result.DebugResetToken != "" {
		response["debugResetToken"] = result.DebugResetToken
		response["debugResetUrl"] = result.DebugResetURL
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	var request resetPasswordRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	if err := s.auth.ResetPassword(r.Context(), request.Token, request.Password); err != nil {
		writeAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "password updated"})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var request loginRequest
	if !decodeJSON(w, r, &request) {
		return
	}

	session, err := s.auth.Login(r.Context(), auth.LoginInput{
		Email:    request.Email,
		Password: request.Password,
		Attempt:  s.loginAttemptMetadata(r),
	})
	if err != nil {
		writeAuthError(w, err)
		return
	}
	s.writeSession(w, session)
}

func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		writeAuthError(w, &auth.Error{Kind: auth.KindUnauthorized, Message: "missing refresh token"})
		return
	}
	session, err := s.auth.Refresh(r.Context(), cookie.Value)
	if err != nil {
		http.SetCookie(w, s.refreshCookie("", time.Unix(0, 0)))
		writeAuthError(w, err)
		return
	}
	s.writeSession(w, session)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("refresh_token"); err == nil {
		s.auth.Logout(r.Context(), cookie.Value)
	}

	http.SetCookie(w, s.refreshCookie("", time.Unix(0, 0)))
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) writeSession(w http.ResponseWriter, session auth.Session) {
	if session.RefreshToken != "" {
		http.SetCookie(w, s.refreshCookie(session.RefreshToken, session.RefreshExpiresAt))
	}

	writeJSON(w, http.StatusOK, sessionResponse{
		AccessToken: session.AccessToken,
		ExpiresAt:   session.AccessExpiresAt.Format(time.RFC3339),
		User: userResponse{
			ID:            session.User.ID.String(),
			Email:         session.User.Email,
			Name:          session.User.Name,
			Status:        session.User.Status,
			EmailVerified: session.User.EmailVerifiedAt != nil,
			Licenses:      session.Licenses,
			Roles:         session.Roles,
			Permissions:   session.Permissions,
		},
	})
}

func (s *Server) refreshCookie(value string, expiresAt time.Time) *http.Cookie {
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

func (s *Server) claimsFromBearer(r *http.Request) (token.Claims, error) {
	header := r.Header.Get("Authorization")
	rawToken, ok := strings.CutPrefix(header, "Bearer ")
	if !ok || rawToken == "" {
		return token.Claims{}, errors.New("missing bearer token")
	}

	return s.tokens.Parse(rawToken)
}

func bearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	rawToken, ok := strings.CutPrefix(header, "Bearer ")
	if !ok {
		return ""
	}
	return rawToken
}

func (s *Server) loginAttemptMetadata(r *http.Request) auth.LoginAttemptMetadata {
	browser, operatingSystem := requestmeta.DeviceFromUserAgent(r.UserAgent())
	return auth.LoginAttemptMetadata{
		IPAddress:       s.meta.ClientIPAddress(r),
		CountryCode:     requestmeta.CountryFromRequest(r),
		City:            requestmeta.CityFromRequest(r),
		UserAgent:       r.UserAgent(),
		Browser:         browser,
		OperatingSystem: operatingSystem,
		CorrelationID:   r.Header.Get("X-Request-Id"),
	}
}

func writeAuthError(w http.ResponseWriter, err error) {
	var serviceError *auth.Error
	if !errors.As(err, &serviceError) {
		writeErrorCode(w, http.StatusInternalServerError, string(auth.KindInternal), "internal server error")
		return
	}

	status := http.StatusInternalServerError
	switch serviceError.Kind {
	case auth.KindBadRequest:
		status = http.StatusBadRequest
	case auth.KindUnauthorized, auth.KindMissingToken, auth.KindInvalidToken:
		status = http.StatusUnauthorized
	case auth.KindForbidden:
		status = http.StatusForbidden
	case auth.KindConflict, auth.KindEmailConflict:
		status = http.StatusConflict
	case auth.KindBadGateway:
		status = http.StatusBadGateway
	case auth.KindRateLimited:
		status = http.StatusTooManyRequests
	case auth.KindRefreshReuse:
		status = http.StatusUnauthorized
	}
	writeErrorCode(w, status, string(serviceError.Kind), serviceError.Message)
}
