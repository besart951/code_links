import type {
  AdminDashboardSummary,
  AdminMeResponse,
  AdminSearchResponse,
  AdminSetting,
  AdminTenantSummary,
  AdminUserSummary,
  AuditLogEntry,
  EntitlementsResponse,
  FeatureKey,
  MeResponse,
  NotificationDeliverySummary,
  NotificationTemplateSummary,
  PageResponse,
  ProductKey,
  SecurityEventSummary
} from '@codelinks/contracts';

export interface PlatformClientOptions {
  baseUrl?: string;
  fetch?: typeof fetch;
  credentials?: RequestCredentials;
}

export class PlatformClient {
  private readonly baseUrl: string;
  private readonly fetchImpl: typeof fetch;
  private readonly credentials: RequestCredentials;

  constructor(options: PlatformClientOptions = {}) {
    this.baseUrl = options.baseUrl ?? '/platform';
    this.fetchImpl = options.fetch ?? fetch;
    this.credentials = options.credentials ?? 'include';
  }

  async me(): Promise<MeResponse> {
    return this.getJson<MeResponse>('/api/v1/auth/me');
  }

  async entitlements(): Promise<EntitlementsResponse> {
    return this.getJson<EntitlementsResponse>('/api/v1/auth/entitlements');
  }

  async adminMe(): Promise<AdminMeResponse> {
    return this.getJson<AdminMeResponse>('/api/v1/admin/me');
  }

  async adminDashboard(): Promise<AdminDashboardSummary> {
    return this.getJson<AdminDashboardSummary>('/api/v1/admin/dashboard');
  }

  async adminSearch(query: string): Promise<AdminSearchResponse> {
    return this.getJson<AdminSearchResponse>(`/api/v1/admin/search?q=${encodeURIComponent(query)}`);
  }

  async adminTenants(): Promise<PageResponse<AdminTenantSummary>> {
    return this.getJson<PageResponse<AdminTenantSummary>>('/api/v1/admin/tenants');
  }

  async adminTenant(tenantId: string): Promise<AdminTenantSummary> {
    return this.getJson<AdminTenantSummary>(`/api/v1/admin/tenants/${encodeURIComponent(tenantId)}`);
  }

  async adminUsers(): Promise<PageResponse<AdminUserSummary>> {
    return this.getJson<PageResponse<AdminUserSummary>>('/api/v1/admin/users');
  }

  async adminUser(userId: string): Promise<AdminUserSummary> {
    return this.getJson<AdminUserSummary>(`/api/v1/admin/users/${encodeURIComponent(userId)}`);
  }

  async adminAudit(): Promise<PageResponse<AuditLogEntry>> {
    return this.getJson<PageResponse<AuditLogEntry>>('/api/v1/admin/audit');
  }

  async adminNotifications(): Promise<{
    templates: NotificationTemplateSummary[];
    deliveries: NotificationDeliverySummary[];
  }> {
    return this.getJson('/api/v1/admin/notifications');
  }

  async adminSecurity(): Promise<PageResponse<SecurityEventSummary>> {
    return this.getJson<PageResponse<SecurityEventSummary>>('/api/v1/admin/security');
  }

  async adminSettings(): Promise<PageResponse<AdminSetting>> {
    return this.getJson<PageResponse<AdminSetting>>('/api/v1/admin/settings');
  }

  async hasFeature(tenantId: string, productKey: ProductKey, featureKey: FeatureKey): Promise<boolean> {
    const response = await this.entitlements();
    return response.entitlements.some((entitlement) => {
      if (!entitlement.enabled) return false;
      if (entitlement.tenant_id !== tenantId) return false;
      if (entitlement.product_key !== productKey) return false;
      if (entitlement.feature_key !== featureKey) return false;
      if (!entitlement.expires_at) return true;
      return new Date(entitlement.expires_at).getTime() > Date.now();
    });
  }

  private async getJson<T>(path: string): Promise<T> {
    const response = await this.fetchImpl(`${this.baseUrl}${path}`, {
      credentials: this.credentials,
      headers: { accept: 'application/json' }
    });

    if (!response.ok) {
      throw new Error(`Platform request failed with ${response.status}`);
    }

    return (await response.json()) as T;
  }
}

export function createPlatformClient(options: PlatformClientOptions = {}): PlatformClient {
  return new PlatformClient(options);
}
