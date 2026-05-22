export type SmtpEncryption = 'none' | 'ssl' | 'tls' | 'starttls';

export interface SmtpSettings {
	host: string;
	port: number;
	username: string;
	hasPassword: boolean;
	encryption: SmtpEncryption;
	fromEmail: string;
	fromName: string;
	replyToEmail: string;
	active: boolean;
	updatedAt: string;
}

export interface UpdateSmtpSettingsInput {
	host: string;
	port: number;
	username: string;
	password?: string;
	encryption: SmtpEncryption;
	fromEmail: string;
	fromName: string;
	replyToEmail: string;
	active: boolean;
}
