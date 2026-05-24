package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	adminsvc "github.com/besart951/code-links/apps/auth-service/backend/internal/admin"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/domain"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/requestmeta"
	"github.com/google/uuid"
)

func (s *Server) registerAdminRoutes(mux *http.ServeMux) {
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
	mux.HandleFunc("GET /api/admin/audit-entries", s.handleAdminAuditEntries)
	mux.HandleFunc("GET /api/admin/runtime-logs", s.handleAdminRuntimeLogs)
	mux.HandleFunc("GET /api/admin/settings/smtp", s.handleAdminGetSMTPSettings)
	mux.HandleFunc("PUT /api/admin/settings/smtp", s.handleAdminUpdateSMTPSettings)
	mux.HandleFunc("POST /api/admin/settings/smtp/test-email", s.handleAdminSendTestEmail)
}

func (s *Server) handleAdminMe(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, domain.PermissionDashboardRead)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, actor)
}

func (s *Server) handleAdminDashboardStats(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, domain.PermissionDashboardRead)
	if !ok {
		return
	}
	stats, err := s.admin.DashboardStats(r.Context(), actor)
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (s *Server) handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, domain.PermissionUsersRead)
	if !ok {
		return
	}
	result, err := s.admin.ListUsers(r.Context(), actor, domain.AdminUserListQuery{
		Query:     r.URL.Query().Get("query"),
		Role:      r.URL.Query().Get("role"),
		Status:    domain.UserStatus(r.URL.Query().Get("status")),
		Page:      queryInt(r, "page", 1),
		PageSize:  queryInt(r, "pageSize", 20),
		Sort:      r.URL.Query().Get("sort"),
		Direction: r.URL.Query().Get("direction"),
	})
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleAdminUserDetail(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, domain.PermissionUsersRead)
	if !ok {
		return
	}
	userID, ok := pathUUID(w, r, "id")
	if !ok {
		return
	}
	user, err := s.admin.UserDetail(r.Context(), actor, userID)
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (s *Server) handleAdminUpdateUser(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, domain.PermissionUsersUpdate)
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
	user, err := s.admin.SetUserStatus(r.Context(), actor, userID, request.Status, adminMeta(r))
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (s *Server) handleAdminLockUser(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, domain.PermissionUsersLock)
	if !ok {
		return
	}
	userID, ok := pathUUID(w, r, "id")
	if !ok {
		return
	}
	user, err := s.admin.LockUser(r.Context(), actor, userID, adminMeta(r))
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (s *Server) handleAdminUnlockUser(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, domain.PermissionUsersLock)
	if !ok {
		return
	}
	userID, ok := pathUUID(w, r, "id")
	if !ok {
		return
	}
	user, err := s.admin.UnlockUser(r.Context(), actor, userID, adminMeta(r))
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (s *Server) handleAdminUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, domain.PermissionUsersChangeRole)
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
	user, err := s.admin.SetUserRole(r.Context(), actor, userID, request.Role, adminMeta(r))
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (s *Server) handleAdminLoginAttempts(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, domain.PermissionAuthLogsRead)
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

	result, err := s.admin.ListLoginAttempts(r.Context(), actor, domain.LoginAttemptListQuery{
		UserID:   userID,
		Success:  success,
		Query:    r.URL.Query().Get("query"),
		Page:     queryInt(r, "page", 1),
		PageSize: queryInt(r, "pageSize", 20),
	})
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleAdminSecurityEvents(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, domain.PermissionSecurityEventsRead)
	if !ok {
		return
	}
	events, err := s.admin.ListSecurityEvents(r.Context(), actor, queryInt(r, "limit", 50))
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (s *Server) handleAdminNotifications(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r, domain.PermissionNotificationsRead); !ok {
		return
	}
	notifications, err := s.admin.ListNotifications(r.Context(), queryInt(r, "limit", 50))
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, notifications)
}

func (s *Server) handleAdminAuditEntries(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r, domain.PermissionAuditEntriesRead); !ok {
		return
	}
	entries, err := s.admin.ListAuditEntries(r.Context(), queryInt(r, "limit", 50))
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

func (s *Server) handleAdminRuntimeLogs(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r, domain.PermissionAuditEntriesRead); !ok {
		return
	}
	entries, err := s.admin.ListRuntimeLogs(r.Context(), queryInt(r, "limit", 100))
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

func (s *Server) handleAdminGetSMTPSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r, domain.PermissionSMTPSettingsRead); !ok {
		return
	}
	settings, err := s.admin.SMTPSettings(r.Context())
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func (s *Server) handleAdminUpdateSMTPSettings(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, domain.PermissionSMTPSettingsUpdate)
	if !ok {
		return
	}
	var request smtpSettingsRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	settings, err := s.admin.UpdateSMTPSettings(r.Context(), actor, adminsvc.SMTPSettingsInput{
		Host:         request.Host,
		Port:         request.Port,
		Username:     request.Username,
		Password:     request.Password,
		Encryption:   request.Encryption,
		FromEmail:    request.FromEmail,
		FromName:     request.FromName,
		ReplyToEmail: request.ReplyToEmail,
		Active:       request.Active,
	}, adminMeta(r))
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func (s *Server) handleAdminSendTestEmail(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, domain.PermissionSMTPSettingsUpdate)
	if !ok {
		return
	}
	var request testEmailRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	if err := s.admin.SendTestEmail(r.Context(), actor, request.Recipient, adminMeta(r)); err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "test email sent"})
}

func (s *Server) requireAdmin(w http.ResponseWriter, r *http.Request, permission domain.AdminPermission) (domain.AdminActor, bool) {
	actor, err := s.admin.ResolveActor(r.Context(), adminsvc.Authn{
		BearerToken:  bearerToken(r),
		RefreshToken: refreshCookieValue(r),
	}, permission)
	if err != nil {
		writeAdminError(w, err)
		return domain.AdminActor{}, false
	}
	return actor, true
}

func refreshCookieValue(r *http.Request) string {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func adminMeta(r *http.Request) adminsvc.RequestMeta {
	return adminsvc.RequestMeta{IPAddress: requestmeta.ClientIPAddress(r)}
}

func writeAdminError(w http.ResponseWriter, err error) {
	var serviceError *adminsvc.Error
	if !errors.As(err, &serviceError) {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	status := http.StatusInternalServerError
	switch serviceError.Kind {
	case adminsvc.KindBadRequest:
		status = http.StatusBadRequest
	case adminsvc.KindUnauthorized:
		status = http.StatusUnauthorized
	case adminsvc.KindForbidden:
		status = http.StatusForbidden
	case adminsvc.KindNotFound:
		status = http.StatusNotFound
	case adminsvc.KindBadGateway:
		status = http.StatusBadGateway
	}
	writeError(w, status, serviceError.Message)
}
