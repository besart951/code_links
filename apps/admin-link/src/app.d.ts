import type { AdminActor } from '$lib/domain/admin-access/types';

declare global {
	namespace App {
		interface Locals {
			admin: AdminActor | null;
		}
	}
}

export {};
