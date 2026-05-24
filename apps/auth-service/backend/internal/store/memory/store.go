package memory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Store struct {
	mu                  sync.RWMutex
	usersByID           map[uuid.UUID]User
	usersByEmail        map[string]uuid.UUID
	roles               map[uuid.UUID][]AdminRole
	licenses            map[uuid.UUID]map[string]bool
	sessions            map[string]refreshSession
	passwordResetTokens map[string]oneTimeToken
	verificationTokens  map[string]oneTimeToken
	loginAttempts       []LoginAttempt
	securityEvents      []SecurityEvent
	notifications       []Notification
	auditEntries        []AdminAuditEntry
	smtpSettings        SmtpSettings
}

type refreshSession struct {
	userID    uuid.UUID
	expiresAt time.Time
}

type oneTimeToken struct {
	userID    uuid.UUID
	expiresAt time.Time
	usedAt    *time.Time
}

func New() (*Store, error) {
	now := time.Now().UTC().Truncate(time.Second)
	adminID := uuid.MustParse("00000000-0000-4000-8000-000000000001")
	supportID := uuid.MustParse("00000000-0000-4000-8000-000000000002")
	auditorID := uuid.MustParse("00000000-0000-4000-8000-000000000003")

	adminPasswordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	supportPasswordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	auditorPasswordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	users := []User{
		{
			ID:              adminID,
			Email:           "demo@codelinks.dev",
			Name:            "Demo Admin",
			PasswordHash:    string(adminPasswordHash),
			Status:          UserStatusActive,
			EmailVerifiedAt: &now,
			CreatedAt:       now.Add(-20 * 24 * time.Hour),
			UpdatedAt:       now,
		},
		{
			ID:              supportID,
			Email:           "support@codelinks.dev",
			Name:            "Support User",
			PasswordHash:    string(supportPasswordHash),
			Status:          UserStatusActive,
			EmailVerifiedAt: &now,
			CreatedAt:       now.Add(-14 * 24 * time.Hour),
			UpdatedAt:       now,
		},
		{
			ID:              auditorID,
			Email:           "auditor@codelinks.dev",
			Name:            "Audit Reader",
			PasswordHash:    string(auditorPasswordHash),
			Status:          UserStatusActive,
			EmailVerifiedAt: &now,
			CreatedAt:       now.Add(-10 * 24 * time.Hour),
			UpdatedAt:       now,
		},
	}

	store := &Store{
		usersByID:           map[uuid.UUID]User{},
		usersByEmail:        map[string]uuid.UUID{},
		roles:               map[uuid.UUID][]AdminRole{adminID: {AdminRoleAdmin}, supportID: {AdminRoleSupport}, auditorID: {AdminRoleAuditor}},
		licenses:            map[uuid.UUID]map[string]bool{adminID: {"infra-link": true}},
		sessions:            map[string]refreshSession{},
		passwordResetTokens: map[string]oneTimeToken{},
		verificationTokens:  map[string]oneTimeToken{},
		loginAttempts:       []LoginAttempt{},
		securityEvents: []SecurityEvent{{
			ID:              uuid.New(),
			Type:            "many_failed_logins",
			Severity:        "medium",
			Status:          "open",
			Summary:         "Mehrere fehlgeschlagene Logins wurden erkannt",
			DetectedAt:      now.Add(-2 * time.Hour),
			SourceIPAddress: "192.0.2.77",
			CountryCode:     "DE",
		}},
		notifications: []Notification{{
			ID:        uuid.New(),
			Type:      "system",
			Channel:   "email",
			Recipient: "admin@codelinks.dev",
			Subject:   "SMTP noch nicht konfiguriert",
			Status:    "pending",
			CreatedAt: now.Add(-1 * time.Hour),
		}},
		auditEntries: []AdminAuditEntry{},
		smtpSettings: SmtpSettings{
			Host:         "",
			Port:         587,
			Encryption:   SmtpEncryptionStartTLS,
			FromEmail:    "no-reply@codelinks.dev",
			FromName:     "CodeLinks",
			ReplyToEmail: "support@codelinks.dev",
			Active:       false,
			UpdatedAt:    now,
		},
	}

	for _, user := range users {
		store.putUserLocked(user)
	}

	return store, nil
}

func (s *Store) FindUserByEmail(_ context.Context, email string) (User, []string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, ok := s.usersByEmail[normalizeEmail(email)]
	if !ok {
		return User{}, nil, errNotFound
	}

	user := s.usersByID[userID]
	return user, s.userLicensesLocked(user.ID), nil
}

func (s *Store) FindUserByID(_ context.Context, userID uuid.UUID) (User, []string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.usersByID[userID]
	if !ok {
		return User{}, nil, errNotFound
	}

	return user, s.userLicensesLocked(user.ID), nil
}

func (s *Store) CreateUser(_ context.Context, name string, email string, passwordHash string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	email = normalizeEmail(email)
	if _, ok := s.usersByEmail[email]; ok {
		return User{}, errConflict
	}

	now := time.Now().UTC()
	user := User{
		ID:           uuid.New(),
		Email:        email,
		Name:         strings.TrimSpace(name),
		PasswordHash: passwordHash,
		Status:       UserStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if user.Name == "" {
		user.Name = email
	}

	s.putUserLocked(user)
	s.roles[user.ID] = []AdminRole{AdminRoleUser}
	return user, nil
}

func (s *Store) GrantLicense(_ context.Context, userID uuid.UUID, productID string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.usersByID[userID]; !ok {
		return nil, errNotFound
	}
	if s.licenses[userID] == nil {
		s.licenses[userID] = map[string]bool{}
	}
	s.licenses[userID][productID] = true

	return s.userLicensesLocked(userID), nil
}

func (s *Store) CreateRefreshSession(_ context.Context, tokenHash string, userID uuid.UUID, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[tokenHash] = refreshSession{userID: userID, expiresAt: expiresAt}
	return nil
}

func (s *Store) FindRefreshSession(_ context.Context, tokenHash string) (uuid.UUID, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[tokenHash]
	if !ok || time.Now().UTC().After(session.expiresAt) {
		return uuid.Nil, errNotFound
	}

	return session.userID, nil
}

func (s *Store) DeleteRefreshSession(_ context.Context, tokenHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, tokenHash)
	return nil
}

func (s *Store) CreateEmailVerificationToken(_ context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.usersByID[userID]; !ok {
		return errNotFound
	}
	s.verificationTokens[tokenHash] = oneTimeToken{userID: userID, expiresAt: expiresAt}
	return nil
}

func (s *Store) VerifyEmailToken(_ context.Context, tokenHash string, now time.Time) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.verificationTokens[tokenHash]
	if !ok || record.usedAt != nil || now.After(record.expiresAt) {
		return User{}, errNotFound
	}

	user, ok := s.usersByID[record.userID]
	if !ok {
		return User{}, errNotFound
	}

	record.usedAt = &now
	s.verificationTokens[tokenHash] = record
	user.EmailVerifiedAt = &now
	user.UpdatedAt = now
	s.putUserLocked(user)

	return user, nil
}

func (s *Store) CreatePasswordResetToken(_ context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.usersByID[userID]; !ok {
		return errNotFound
	}
	s.passwordResetTokens[tokenHash] = oneTimeToken{userID: userID, expiresAt: expiresAt}
	return nil
}

func (s *Store) ResetPasswordByToken(_ context.Context, tokenHash string, passwordHash string, now time.Time) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.passwordResetTokens[tokenHash]
	if !ok || record.usedAt != nil || now.After(record.expiresAt) {
		return User{}, errNotFound
	}

	user, ok := s.usersByID[record.userID]
	if !ok {
		return User{}, errNotFound
	}

	record.usedAt = &now
	s.passwordResetTokens[tokenHash] = record
	user.PasswordHash = passwordHash
	user.UpdatedAt = now
	s.putUserLocked(user)

	return user, nil
}

func (s *Store) RecordLoginAttempt(_ context.Context, attempt LoginAttempt) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if attempt.ID == uuid.Nil {
		attempt.ID = uuid.New()
	}
	if attempt.OccurredAt.IsZero() {
		attempt.OccurredAt = time.Now().UTC()
	}

	s.loginAttempts = append(s.loginAttempts, attempt)

	if attempt.Success && attempt.UserID != nil {
		user, ok := s.usersByID[*attempt.UserID]
		if ok {
			user.LastLoginAt = &attempt.OccurredAt
			user.LastLoginIP = attempt.IPAddress
			user.LastLoginCountryCode = attempt.CountryCode
			user.UpdatedAt = attempt.OccurredAt
			s.putUserLocked(user)
		}
	}

	s.detectSecurityEventsLocked(attempt)
	return nil
}

func (s *Store) GetAdminActor(_ context.Context, userID uuid.UUID) (AdminActor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.usersByID[userID]
	if !ok {
		return AdminActor{}, errNotFound
	}

	roles := s.roles[userID]
	permissions := permissionsForRoles(roles)
	if len(permissions) == 0 {
		return AdminActor{}, errForbidden
	}

	return AdminActor{
		ID:          user.ID.String(),
		Email:       user.Email,
		Name:        user.Name,
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

func (s *Store) ListUserRoles(_ context.Context, userID uuid.UUID) ([]AdminRole, []AdminPermission, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.usersByID[userID]; !ok {
		return nil, nil, errNotFound
	}

	roles := append([]AdminRole{}, s.roles[userID]...)
	return roles, permissionsForRoles(roles), nil
}

func (s *Store) ListAdminUsers(_ context.Context, query AdminUserListQuery) (AdminUserListResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]AdminUserListItem, 0, len(s.usersByID))
	for _, user := range s.usersByID {
		item := s.adminUserListItemLocked(user)
		if !matchesUserQuery(item, query) {
			continue
		}
		items = append(items, item)
	}

	sortAdminUsers(items, query.Sort, query.Direction)
	total := len(items)
	page, pageSize := normalizedPage(query.Page, query.PageSize)
	start := (page - 1) * pageSize
	if start > len(items) {
		start = len(items)
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}

	return AdminUserListResult{
		Items:    items[start:end],
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *Store) GetManagedUserDetail(_ context.Context, userID uuid.UUID) (ManagedUserDetail, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.usersByID[userID]
	if !ok {
		return ManagedUserDetail{}, errNotFound
	}

	item := s.adminUserListItemLocked(user)
	attempts := s.loginAttemptsForUserLocked(userID)
	knownIPs := uniqueStringsFromAttempts(attempts, func(attempt LoginAttempt) string { return attempt.IPAddress })
	countries := uniqueStringsFromAttempts(attempts, func(attempt LoginAttempt) string { return attempt.CountryCode })
	devices := uniqueStringsFromAttempts(attempts, func(attempt LoginAttempt) string {
		return strings.TrimSpace(attempt.Browser + " / " + attempt.OperatingSystem)
	})

	roles := []UserPermissionGrant{}
	for _, role := range s.roles[userID] {
		roles = append(roles, UserPermissionGrant{
			Role:      role,
			GrantedAt: user.CreatedAt,
			GrantedBy: "system",
		})
	}

	return ManagedUserDetail{
		AdminUserListItem: item,
		Roles:             roles,
		ProductLicenses:   s.userLicensesLocked(userID),
		KnownIPAddresses:  knownIPs,
		LoginCountries:    countries,
		UsedDevices:       devices,
	}, nil
}

func (s *Store) SetUserStatus(_ context.Context, userID uuid.UUID, status UserStatus) (User, error) {
	if !validUserStatus(status) {
		return User{}, errNotFound
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.usersByID[userID]
	if !ok {
		return User{}, errNotFound
	}
	user.Status = status
	user.UpdatedAt = time.Now().UTC()
	s.putUserLocked(user)
	return user, nil
}

func (s *Store) SetUserRole(_ context.Context, userID uuid.UUID, role AdminRole) (User, error) {
	if !validAdminRole(role) {
		return User{}, errNotFound
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.usersByID[userID]
	if !ok {
		return User{}, errNotFound
	}

	s.roles[userID] = []AdminRole{role}
	user.UpdatedAt = time.Now().UTC()
	s.putUserLocked(user)
	return user, nil
}

func (s *Store) ListLoginAttempts(_ context.Context, query LoginAttemptListQuery) (LoginAttemptListResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := []LoginAttempt{}
	for _, attempt := range s.loginAttempts {
		if query.UserID != nil && (attempt.UserID == nil || *attempt.UserID != *query.UserID) {
			continue
		}
		if query.Success != nil && attempt.Success != *query.Success {
			continue
		}
		if query.Query != "" {
			search := strings.ToLower(query.Query)
			if !strings.Contains(strings.ToLower(attempt.EmailAttempted), search) && !strings.Contains(attempt.IPAddress, search) {
				continue
			}
		}
		items = append(items, attempt)
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].OccurredAt.After(items[j].OccurredAt)
	})

	total := len(items)
	page, pageSize := normalizedPage(query.Page, query.PageSize)
	start := (page - 1) * pageSize
	if start > len(items) {
		start = len(items)
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}

	return LoginAttemptListResult{Items: items[start:end], Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Store) ListSecurityEvents(_ context.Context, limit int) ([]SecurityEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := append([]SecurityEvent{}, s.securityEvents...)
	sort.SliceStable(events, func(i, j int) bool {
		return events[i].DetectedAt.After(events[j].DetectedAt)
	})
	if limit > 0 && len(events) > limit {
		events = events[:limit]
	}

	return events, nil
}

func (s *Store) GetDashboardStats(_ context.Context) (DashboardStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var stats DashboardStats
	now := time.Now().UTC()
	stats.Users.Total = len(s.usersByID)
	for _, user := range s.usersByID {
		switch user.Status {
		case UserStatusActive:
			stats.Users.Active++
		case UserStatusLocked:
			stats.Users.Locked++
		}
		if user.CreatedAt.After(now.Add(-7 * 24 * time.Hour)) {
			stats.Users.NewLast7Days++
		}
		if user.CreatedAt.After(now.Add(-30 * 24 * time.Hour)) {
			stats.Users.NewLast30Days++
		}
	}

	countryCounts := map[string]int{}
	ipCounts := map[string]int{}
	for _, attempt := range s.loginAttempts {
		stats.LoginAttempts.Total++
		if attempt.Success {
			stats.LoginAttempts.Successful++
		} else {
			stats.LoginAttempts.Failed++
		}
		if attempt.CountryCode != "" {
			countryCounts[attempt.CountryCode]++
		}
		if attempt.IPAddress != "" {
			ipCounts[attempt.IPAddress]++
		}
	}

	for _, token := range s.passwordResetTokens {
		if token.expiresAt.After(now.Add(-24 * time.Hour)) {
			stats.PasswordResetRequests++
		}
	}
	stats.Notifications = len(s.notifications)
	for _, event := range s.securityEvents {
		if event.Status == "open" {
			stats.OpenSecurityEvents++
		}
	}
	stats.TopCountries = topCounts(countryCounts, 5)
	stats.TopIPAddresses = topCounts(ipCounts, 5)

	return stats, nil
}

func (s *Store) GetSmtpSettings(_ context.Context) (SmtpSettings, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	settings := s.smtpSettings
	settings.HasPassword = settings.PasswordEncrypted != ""
	return settings, nil
}

func (s *Store) SaveSmtpSettings(_ context.Context, settings SmtpSettings) (SmtpSettings, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if settings.PasswordEncrypted == "" {
		settings.PasswordEncrypted = s.smtpSettings.PasswordEncrypted
	}
	settings.HasPassword = settings.PasswordEncrypted != ""
	settings.UpdatedAt = time.Now().UTC()
	s.smtpSettings = settings

	clean := s.smtpSettings
	clean.PasswordEncrypted = ""
	return clean, nil
}

func (s *Store) ListNotifications(_ context.Context, limit int) ([]Notification, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	notifications := append([]Notification{}, s.notifications...)
	sort.SliceStable(notifications, func(i, j int) bool {
		return notifications[i].CreatedAt.After(notifications[j].CreatedAt)
	})
	if limit > 0 && len(notifications) > limit {
		notifications = notifications[:limit]
	}

	return notifications, nil
}

func (s *Store) CreateNotification(_ context.Context, notification Notification) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if notification.ID == uuid.Nil {
		notification.ID = uuid.New()
	}
	if notification.CreatedAt.IsZero() {
		notification.CreatedAt = time.Now().UTC()
	}
	s.notifications = append(s.notifications, notification)
	return nil
}

func (s *Store) RecordAdminAuditEntry(_ context.Context, entry AdminAuditEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	s.auditEntries = append(s.auditEntries, entry)
	return nil
}

func (s *Store) ListAdminAuditEntries(_ context.Context, limit int) ([]AdminAuditEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries := append([]AdminAuditEntry{}, s.auditEntries...)
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].CreatedAt.After(entries[j].CreatedAt)
	})
	if limit > 0 && len(entries) > limit {
		entries = entries[:limit]
	}

	return entries, nil
}

func (s *Store) putUserLocked(user User) {
	s.usersByID[user.ID] = user
	s.usersByEmail[normalizeEmail(user.Email)] = user.ID
}

func (s *Store) userLicensesLocked(userID uuid.UUID) []string {
	licenseMap := s.licenses[userID]
	licenses := make([]string, 0, len(licenseMap))
	for productID := range licenseMap {
		licenses = append(licenses, productID)
	}
	sort.Strings(licenses)

	return licenses
}

func (s *Store) adminUserListItemLocked(user User) AdminUserListItem {
	successful, failed := s.loginCountsLocked(user.ID)

	return AdminUserListItem{
		ID:                   user.ID.String(),
		Name:                 user.Name,
		Email:                user.Email,
		PrimaryRole:          primaryRole(s.roles[user.ID]),
		Status:               user.Status,
		EmailVerified:        user.EmailVerifiedAt != nil,
		CreatedAt:            user.CreatedAt,
		LastLoginAt:          user.LastLoginAt,
		SuccessfulLoginCount: successful,
		FailedLoginCount:     failed,
		LastKnownIPAddress:   emptyStringToNil(user.LastLoginIP),
		LastLoginCountryCode: emptyStringToNil(user.LastLoginCountryCode),
	}
}

func (s *Store) loginCountsLocked(userID uuid.UUID) (int, int) {
	successful := 0
	failed := 0
	for _, attempt := range s.loginAttempts {
		if attempt.UserID == nil || *attempt.UserID != userID {
			continue
		}
		if attempt.Success {
			successful++
		} else {
			failed++
		}
	}

	return successful, failed
}

func (s *Store) loginAttemptsForUserLocked(userID uuid.UUID) []LoginAttempt {
	attempts := []LoginAttempt{}
	for _, attempt := range s.loginAttempts {
		if attempt.UserID != nil && *attempt.UserID == userID {
			attempts = append(attempts, attempt)
		}
	}

	return attempts
}

func (s *Store) detectSecurityEventsLocked(attempt LoginAttempt) {
	if attempt.Success {
		return
	}

	failuresFromIP := 0
	since := attempt.OccurredAt.Add(-15 * time.Minute)
	for _, candidate := range s.loginAttempts {
		if candidate.Success || candidate.IPAddress != attempt.IPAddress || candidate.OccurredAt.Before(since) {
			continue
		}
		failuresFromIP++
	}

	if failuresFromIP >= 5 {
		s.securityEvents = append(s.securityEvents, SecurityEvent{
			ID:              uuid.New(),
			UserID:          attempt.UserID,
			Type:            "many_failures_from_ip",
			Severity:        "high",
			Status:          "open",
			Summary:         "Viele fehlgeschlagene Logins von derselben IP-Adresse",
			DetectedAt:      attempt.OccurredAt,
			SourceIPAddress: attempt.IPAddress,
			CountryCode:     attempt.CountryCode,
		})
	}

	if attempt.FailureReason != nil && *attempt.FailureReason == LoginFailureAccountLocked {
		s.securityEvents = append(s.securityEvents, SecurityEvent{
			ID:              uuid.New(),
			UserID:          attempt.UserID,
			Type:            "locked_account_login",
			Severity:        "medium",
			Status:          "open",
			Summary:         "Login-Versuch auf gesperrten Account",
			DetectedAt:      attempt.OccurredAt,
			SourceIPAddress: attempt.IPAddress,
			CountryCode:     attempt.CountryCode,
		})
	}
}

func matchesUserQuery(item AdminUserListItem, query AdminUserListQuery) bool {
	search := strings.ToLower(strings.TrimSpace(query.Query))
	if search != "" && !strings.Contains(strings.ToLower(item.Name), search) && !strings.Contains(strings.ToLower(item.Email), search) {
		return false
	}
	if query.Role != "" && string(item.PrimaryRole) != query.Role {
		return false
	}
	if query.Status != "" && item.Status != query.Status {
		return false
	}

	return true
}

func sortAdminUsers(items []AdminUserListItem, field string, direction string) {
	if field == "" {
		field = "createdAt"
	}
	desc := direction != "asc"
	sort.SliceStable(items, func(i, j int) bool {
		left := items[i]
		right := items[j]
		result := 0
		switch field {
		case "name":
			result = strings.Compare(left.Name, right.Name)
		case "email":
			result = strings.Compare(left.Email, right.Email)
		case "primaryRole":
			result = strings.Compare(string(left.PrimaryRole), string(right.PrimaryRole))
		case "status":
			result = strings.Compare(string(left.Status), string(right.Status))
		case "lastLoginAt":
			result = compareOptionalTime(left.LastLoginAt, right.LastLoginAt)
		default:
			result = left.CreatedAt.Compare(right.CreatedAt)
		}
		if desc {
			return result > 0
		}

		return result < 0
	})
}

func normalizedPage(page int, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}

func compareOptionalTime(left *time.Time, right *time.Time) int {
	if left == nil && right == nil {
		return 0
	}
	if left == nil {
		return -1
	}
	if right == nil {
		return 1
	}

	return left.Compare(*right)
}

func uniqueStringsFromAttempts(attempts []LoginAttempt, getValue func(LoginAttempt) string) []string {
	seen := map[string]bool{}
	values := []string{}
	for _, attempt := range attempts {
		value := strings.TrimSpace(getValue(attempt))
		if value == "" || value == "/" || seen[value] {
			continue
		}
		seen[value] = true
		values = append(values, value)
	}
	sort.Strings(values)

	return values
}

func topCounts(counts map[string]int, limit int) []CountStat {
	stats := make([]CountStat, 0, len(counts))
	for key, count := range counts {
		stats = append(stats, CountStat{Key: key, Count: count})
	}
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Count == stats[j].Count {
			return stats[i].Key < stats[j].Key
		}

		return stats[i].Count > stats[j].Count
	})
	if limit > 0 && len(stats) > limit {
		stats = stats[:limit]
	}

	return stats
}

func emptyStringToNil(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
