import type { AdminMeResponse } from '@codelinks/contracts';

declare global {
	namespace App {
		interface Locals {
			admin?: AdminMeResponse;
		}
	}
}

export {};
