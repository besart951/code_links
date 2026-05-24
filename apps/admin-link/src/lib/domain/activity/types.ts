export type ActivityEventSource = 'auth' | 'security' | 'notification' | 'audit' | 'runtime';
export type ActivityEventTone = 'neutral' | 'success' | 'warning' | 'danger' | 'info';

export interface ActivityEventDetail {
	label: string;
	value: string;
}

export interface ActivityEvent {
	id: string;
	source: ActivityEventSource;
	tone: ActivityEventTone;
	title: string;
	message: string;
	occurredAt: string;
	details: ActivityEventDetail[];
}
