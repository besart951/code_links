export interface AdminAuditEntry {
	id: string;
	actorUserId: string;
	action: string;
	targetType: string;
	targetId: string;
	reason: string;
	ipAddress: string;
	createdAt: string;
}
