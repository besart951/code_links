export type RuntimeLogLevel = 'info' | 'fatal';

export interface RuntimeLogEntry {
	id: string;
	occurredAt: string;
	level: RuntimeLogLevel;
	source: string;
	message: string;
	raw: string;
}
