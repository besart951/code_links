import { describe, expect, it } from 'vitest';
import { Button } from './button/index.js';
import * as Card from './card/index.js';
import { Input } from './input/index.js';

describe('shared UI primitive exports', () => {
	it('exports promoted primitives from shared package paths', () => {
		expect(Button).toBeDefined();
		expect(Card.Root).toBeDefined();
		expect(Input).toBeDefined();
	});
});
