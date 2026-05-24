import type { AdminCommandRepository } from './AdminCommandRepository';
import type { AdminQueryRepository } from './AdminQueryRepository';

export interface AdminRepository extends AdminQueryRepository, AdminCommandRepository {}
