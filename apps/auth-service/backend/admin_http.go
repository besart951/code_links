package main

import (
	"errors"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *server) registerAdminRoutes(mux *http.ServeMux) {
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
}

func (s *server) handleAdminMe(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, PermissionDashboardRead)
	if !ok {
		return
	}

	writeJSON(w, http.StatusOK, actor)
}

func (s *server) handleAdminDashboardStats(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, PermissionDashboardRead)
	if !ok {
		return
	}

	stats, err := s.adminStore.GetDashboardStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load dashboard stats")
		return
	}

	writeJSON(w, http.StatusOK, s.adminProjection.ProjectDashboardStats(actor, stats))
}

func (s *server) handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, PermissionUsersRead)
	if !ok {
		return
	}

	result, err := s.adminStore.ListAdminUsers(r.Context(), AdminUserListQuery{
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

	writeJSON(w, http.StatusOK, s.adminProjection.ProjectUserList(actor, result))
}

func (s *server) handleAdminUserDetail(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, PermissionUsersRead)
	if !ok {
		return
	}

	userID, ok := pathUUID(w, r, "id")
	if !ok {
		return
	}
	user, err := s.adminStore.GetManagedUserDetail(r.Context(), userID)
	if errors.Is(err, errNotFound) {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load user")
		return
	}

	writeJSON(w, http.StatusOK, s.adminProjection.ProjectManagedUserDetail(actor, user))
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
	user, err := s.adminStore.SetUserStatus(r.Context(), userID, request.Status)
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
	user, err := s.adminStore.SetUserStatus(r.Context(), userID, UserStatusLocked)
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
	user, err := s.adminStore.SetUserStatus(r.Context(), userID, UserStatusActive)
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
	user, err := s.adminStore.SetUserRole(r.Context(), userID, request.Role)
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

	result, err := s.adminStore.ListLoginAttempts(r.Context(), LoginAttemptListQuery{
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

	writeJSON(w, http.StatusOK, s.adminProjection.ProjectLoginAttempts(actor, result))
}

func (s *server) handleAdminSecurityEvents(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.requireAdmin(w, r, PermissionSecurityEventsRead)
	if !ok {
		return
	}

	events, err := s.adminStore.ListSecurityEvents(r.Context(), queryInt(r, "limit", 50))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load security events")
		return
	}

	writeJSON(w, http.StatusOK, s.adminProjection.ProjectSecurityEvents(actor, events))
}

func (s *server) handleAdminNotifications(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r, PermissionNotificationsRead); !ok {
		return
	}

	notifications, err := s.notificationStore.ListNotifications(r.Context(), queryInt(r, "limit", 50))
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

	settings, err := s.smtpStore.GetSmtpSettings(r.Context())
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

	settings, err := s.smtpStore.SaveSmtpSettings(r.Context(), SmtpSettings{
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

	settings, err := s.smtpStore.GetSmtpSettings(r.Context())
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
	_ = s.notificationStore.CreateNotification(r.Context(), Notification{
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

		return s.adminStore.GetAdminActor(r.Context(), userID)
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		return AdminActor{}, errNotFound
	}
	userID, err := s.authStore.FindRefreshSession(r.Context(), hashRefreshToken(cookie.Value))
	if err != nil {
		return AdminActor{}, err
	}

	return s.adminStore.GetAdminActor(r.Context(), userID)
}

func (s *server) audit(r *http.Request, actor AdminActor, action string, targetType string, targetID string, reason string) {
	actorID, err := uuid.Parse(actor.ID)
	if err != nil {
		return
	}
	_ = s.auditStore.RecordAdminAuditEntry(r.Context(), AdminAuditEntry{
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
