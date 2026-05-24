import { describe, expect, it } from 'vitest';
import { ProductAppShell } from '@codelinks/ui-library';

describe('loka-link shell contract', () => {
	it('uses the shared product app shell', () => {
		expect(ProductAppShell).toBeDefined();
	});
});
