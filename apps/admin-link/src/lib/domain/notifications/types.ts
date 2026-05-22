export type NotificationStatus = 'queued' | 'pending' | 'sent' | 'failed';

export interface Notification {
	id: string;
	type: string;
	channel: 'email';
	recipient: string;
	subject: string;
	status: NotificationStatus;
	createdAt: string;
	sentAt: string | null;
}
