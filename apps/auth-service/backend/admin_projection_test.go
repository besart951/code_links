package main

import "testing"

func TestAdminProjectionMasksIPMetadataWithoutRawPermission(t *testing.T) {
	policy := AdminProjectionPolicy{}
	actor := AdminActor{
		Permissions: []AdminPermission{PermissionDashboardRead},
	}
	rawIP := "192.0.2.77"

	users := policy.ProjectUserList(actor, AdminUserListResult{
		Items: []AdminUserListItem{{LastKnownIPAddress: &rawIP}},
	})
	if got := *users.Items[0].LastKnownIPAddress; got != "192.0.x.x" {
		t.Fatalf("expected masked user IP, got %q", got)
	}

	attempts := policy.ProjectLoginAttempts(actor, LoginAttemptListResult{
		Items: []LoginAttempt{{IPAddress: "198.51.100.44"}},
	})
	if got := attempts.Items[0].IPAddress; got != "198.51.x.x" {
		t.Fatalf("expected masked login attempt IP, got %q", got)
	}

	events := policy.ProjectSecurityEvents(actor, []SecurityEvent{{SourceIPAddress: "203.0.113.18"}})
	if got := events[0].SourceIPAddress; got != "203.0.x.x" {
		t.Fatalf("expected masked security event IP, got %q", got)
	}

	stats := policy.ProjectDashboardStats(actor, DashboardStats{
		TopIPAddresses: []CountStat{{Key: "203.0.113.81", Count: 2}},
	})
	if got := stats.TopIPAddresses[0].Key; got != "203.0.x.x" {
		t.Fatalf("expected masked dashboard IP, got %q", got)
	}
}

func TestAdminProjectionKeepsRawIPWithRawPermission(t *testing.T) {
	policy := AdminProjectionPolicy{}
	actor := AdminActor{
		Permissions: []AdminPermission{PermissionUsersUpdate},
	}
	rawIP := "192.0.2.77"

	users := policy.ProjectUserList(actor, AdminUserListResult{
		Items: []AdminUserListItem{{LastKnownIPAddress: &rawIP}},
	})
	if got := *users.Items[0].LastKnownIPAddress; got != rawIP {
		t.Fatalf("expected raw user IP, got %q", got)
	}
}
