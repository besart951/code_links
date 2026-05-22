import type { LoginAttempt } from '$lib/domain/auth-logs/types';
import type { SecurityEvent } from '$lib/domain/security/types';

export interface DashboardSummary {
	users: {
		total: number;
		active: number;
		locked: number;
		newLast7Days: number;
		newLast30Days: number;
	};
	loginAttempts: {
		total: number;
		successful: number;
		failed: number;
	};
	passwordResetRequests: number;
	notifications: number;
	security: {
		openEvents: number;
		suspiciousAttempts: number;
	};
	topCountries: { countryCode: string; count: number }[];
	topIpAddresses: { ipAddress: string; count: number }[];
	trend: { date: string; successful: number; failed: number }[];
	recentActivity: LoginAttempt[];
	highlightedEvents: SecurityEvent[];
}
