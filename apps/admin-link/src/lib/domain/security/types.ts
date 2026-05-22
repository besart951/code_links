export type SecurityEventType =
	| 'many_failed_logins'
	| 'unusual_country'
	| 'many_ips_for_user'
	| 'locked_account_attempt'
	| 'many_failures_from_ip';

export type SecuritySeverity = 'low' | 'medium' | 'high' | 'critical';
export type SecurityEventStatus = 'open' | 'resolved';

export interface SecurityEvent {
	id: string;
	userId: string | null;
	type: SecurityEventType;
	severity: SecuritySeverity;
	status: SecurityEventStatus;
	summary: string;
	detectedAt: string;
	resolvedAt: string | null;
	sourceIpAddress: string | null;
	countryCode: string | null;
}
