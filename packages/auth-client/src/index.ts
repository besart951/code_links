import type {
  EntitlementsResponse,
  FeatureKey,
  MeResponse,
  ProductKey
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
    return this.getJson<MeResponse>('/api/v1/me');
  }

  async entitlements(tenantId: string): Promise<EntitlementsResponse> {
    return this.getJson<EntitlementsResponse>(`/api/v1/entitlements?tenant_id=${encodeURIComponent(tenantId)}`);
  }

  async hasFeature(tenantId: string, productKey: ProductKey, featureKey: FeatureKey): Promise<boolean> {
    const response = await this.entitlements(tenantId);
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
