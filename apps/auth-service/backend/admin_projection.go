package main

type AdminProjectionPolicy struct{}

func (AdminProjectionPolicy) CanSeeRawIP(actor AdminActor) bool {
	return hasAdminPermission(actor, PermissionUsersUpdate)
}

func (policy AdminProjectionPolicy) ProjectDashboardStats(actor AdminActor, stats DashboardStats) DashboardStats {
	if policy.CanSeeRawIP(actor) {
		return stats
	}
	for index := range stats.TopIPAddresses {
		stats.TopIPAddresses[index].Key = maskIP(stats.TopIPAddresses[index].Key)
	}
	return stats
}

func (policy AdminProjectionPolicy) ProjectUserList(actor AdminActor, result AdminUserListResult) AdminUserListResult {
	if policy.CanSeeRawIP(actor) {
		return result
	}
	for index := range result.Items {
		result.Items[index] = policy.projectUserListItem(result.Items[index])
	}
	return result
}

func (policy AdminProjectionPolicy) ProjectManagedUserDetail(actor AdminActor, user ManagedUserDetail) ManagedUserDetail {
	if policy.CanSeeRawIP(actor) {
		return user
	}
	user.AdminUserListItem = policy.projectUserListItem(user.AdminUserListItem)
	for index := range user.KnownIPAddresses {
		user.KnownIPAddresses[index] = maskIP(user.KnownIPAddresses[index])
	}
	return user
}

func (policy AdminProjectionPolicy) ProjectLoginAttempts(actor AdminActor, result LoginAttemptListResult) LoginAttemptListResult {
	if policy.CanSeeRawIP(actor) {
		return result
	}
	for index := range result.Items {
		result.Items[index].IPAddress = maskIP(result.Items[index].IPAddress)
	}
	return result
}

func (policy AdminProjectionPolicy) ProjectSecurityEvents(actor AdminActor, events []SecurityEvent) []SecurityEvent {
	if policy.CanSeeRawIP(actor) {
		return events
	}
	for index := range events {
		events[index].SourceIPAddress = maskIP(events[index].SourceIPAddress)
	}
	return events
}

func (AdminProjectionPolicy) projectUserListItem(item AdminUserListItem) AdminUserListItem {
	if item.LastKnownIPAddress != nil {
		masked := maskIP(*item.LastKnownIPAddress)
		item.LastKnownIPAddress = &masked
	}
	return item
}
