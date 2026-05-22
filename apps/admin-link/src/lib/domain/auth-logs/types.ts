export type LoginFailureReason =
	| 'wrong_password'
	| 'unknown_email'
	| 'account_locked'
	| 'too_many_attempts'
	| 'invalid_token'
	| 'email_not_confirmed';

export type AuthMethod = 'password' | 'refresh_token';

export interface LoginAttemptDevice {
	userAgent: string;
	browser: string;
	os: string;
}

export interface LoginAttempt {
	id: string;
	userId: string | null;
	emailAttempted: string;
	occurredAt: string;
	ipAddress: string;
	countryCode: string;
	city: string | null;
	success: boolean;
	failureReason: LoginFailureReason | null;
	authMethod: AuthMethod;
	device: LoginAttemptDevice;
	riskScore: number;
}

export interface LoginAttemptQuery {
	userId?: string;
	success?: boolean;
	query?: string;
	page: number;
	pageSize: number;
}
