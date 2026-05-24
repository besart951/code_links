import { describe, expect, it } from 'vitest';
import { hasPermission } from './permissions';
import type { AdminActor } from './types';

describe('admin permissions', () => {
	it('allows support user locks but blocks SMTP settings', () => {
		expect.assertions(3);

		const admin: AdminActor = {
			id: 'admin',
			name: 'Admin',
			email: 'admin@example.com',
			roles: ['admin'],
			permissions: ['admin.users.update']
		};
		const support: AdminActor = {
			id: 'support',
			name: 'Support',
			email: 'support@example.com',
			roles: ['support'],
			permissions: ['admin.users.lock']
		};

		expect(hasPermission(admin, 'admin.users.update')).toBe(true);
		expect(hasPermission(support, 'admin.users.lock')).toBe(true);
		expect(hasPermission(support, 'admin.smtp_settings.update')).toBe(false);
	});
});
