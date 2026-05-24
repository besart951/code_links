import type { SmtpSettings, UpdateSmtpSettingsInput } from '$lib/domain/smtp/types';
import type { UserStatus } from '$lib/domain/users/types';

export interface AdminCommandRepository {
	updateSmtpSettings(input: UpdateSmtpSettingsInput): Promise<SmtpSettings>;
	sendTestEmail(recipient: string): Promise<void>;
	setUserStatus(userId: string, status: UserStatus): Promise<void>;
	setUserRole(userId: string, role: string): Promise<void>;
}
