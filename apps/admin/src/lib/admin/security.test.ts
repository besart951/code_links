import { describe, expect, it } from 'vitest';
import { hasActionReason, maskEmail, requiresReason } from './security.js';

describe('admin security helpers', () => {
	it('requires reasons for sensitive actions', () => {
		expect(requiresReason('tenant.suspend')).toBe(true);
		expect(requiresReason('dashboard.read')).toBe(false);
	});

	it('rejects short reasons for sensitive actions', () => {
		expect(hasActionReason('tenant.suspend', 'oops')).toBe(false);
		expect(hasActionReason('tenant.suspend', 'Fraud investigation')).toBe(true);
	});

	it('masks email local parts for list views', () => {
		expect(maskEmail('besart@example.com')).toBe('be****@example.com');
	});
});
