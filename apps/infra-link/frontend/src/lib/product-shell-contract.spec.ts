import { describe, expect, it } from 'vitest';
import { ProductAppShell } from '@codelinks/ui-library';

describe('infra-link shell contract', () => {
	it('uses the shared product app shell', () => {
		expect(ProductAppShell).toBeDefined();
	});
});
