import { error, type RequestEvent } from '@sveltejs/kit';
import type { EntitlementsResponse, FeatureKey, ProductKey } from '@codelinks/contracts';

const defaultPlatformUrl = 'http://platform-api:8080';

export async function requireEntitlement(
  event: RequestEvent,
  productKey: ProductKey,
  featureKey: FeatureKey
): Promise<void> {
  if (process.env.PLATFORM_ENTITLEMENTS_REQUIRED !== 'true') {
    return;
  }

  const tenantId = currentTenantId(event);
  const platformUrl = process.env.PLATFORM_URL ?? defaultPlatformUrl;
  const response = await event.fetch(
    `${platformUrl}/api/v1/entitlements?tenant_id=${encodeURIComponent(tenantId)}`,
    {
      headers: {
        accept: 'application/json',
        cookie: event.request.headers.get('cookie') ?? ''
      }
    }
  );

  if (response.status === 401) {
    throw error(401, 'unauthorized');
  }
  if (response.status === 403) {
    throw error(403, 'tenant_membership_required');
  }
  if (!response.ok) {
    throw error(502, 'platform_unavailable');
  }

  const body = (await response.json()) as EntitlementsResponse;
  const allowed = body.entitlements.some((entitlement) => {
    if (!entitlement.enabled) return false;
    if (entitlement.tenant_id !== tenantId) return false;
    if (entitlement.product_key !== productKey) return false;
    if (entitlement.feature_key !== featureKey) return false;
    if (!entitlement.expires_at) return true;
    return new Date(entitlement.expires_at).getTime() > Date.now();
  });

  if (!allowed) {
    throw error(403, 'entitlement_required');
  }
}

function currentTenantId(event: RequestEvent): string {
  const tenantId =
    event.request.headers.get('X-Tenant-ID') ??
    event.url.searchParams.get('tenant_id') ??
    event.cookies.get('codelinks_tenant_id');

  if (!tenantId) {
    throw error(400, 'tenant_id_required');
  }

  return tenantId;
}
