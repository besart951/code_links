const sensitiveActions = new Set([
	'user.lock',
	'user.unlock',
	'user.session.revoke',
	'tenant.suspend',
	'tenant.archive',
	'tenant.reactivate',
	'subscription.change',
	'entitlement.change',
	'setting.change',
	'impersonation.start'
]);

export function requiresReason(action: string): boolean {
	return sensitiveActions.has(action);
}

export function hasActionReason(action: string, reason: string): boolean {
	if (!requiresReason(action)) return true;
	return reason.trim().length >= 8;
}

export function maskEmail(email: string): string {
	const [name, domain] = email.split('@');
	if (!name || !domain) return email;
	const visible = name.slice(0, 2);
	return `${visible}${'*'.repeat(Math.max(name.length - visible.length, 3))}@${domain}`;
}
