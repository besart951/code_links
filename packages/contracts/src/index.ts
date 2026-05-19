export const productKeys = ['infra_link', 'planer_link', 'loka_link'] as const;
export type ProductKey = (typeof productKeys)[number];

export const featureKeys = [
  'planer.pdf_export',
  'planer.excel_export',
  'planer.sync',
  'infra.module_bacnet',
  'infra.module_sps',
  'infra.module_field_devices',
  'loka.core'
] as const;
export type FeatureKey = (typeof featureKeys)[number];

export type TenantType = 'personal' | 'team' | 'company';
export type TenantStatus = 'active' | 'disabled';

export interface AuthUser {
  id: string;
  email: string;
  display_name: string;
  status: string;
}

export interface Tenant {
  id: string;
  type: TenantType;
  name: string;
  slug: string;
  status: TenantStatus;
  role_key?: string;
}

export interface Entitlement {
  tenant_id: string;
  product_key: ProductKey;
  feature_key: FeatureKey;
  enabled: boolean;
  source: 'subscription' | 'manual' | 'trial';
  expires_at: string | null;
}

export interface FeatureLimit {
  tenant_id: string;
  product_key: ProductKey;
  feature_key: FeatureKey;
  limit_key: string;
  value: number;
  period: 'none' | 'day' | 'month' | 'year';
  reset_at: string | null;
}

export interface MeResponse {
  user: AuthUser;
  tenants: Tenant[];
}

export interface EntitlementsResponse {
  tenant_id: string;
  entitlements: Entitlement[];
  limits: FeatureLimit[];
}

export interface AuthorizeRequest {
  user_id: string;
  tenant_id: string;
  product_key: ProductKey;
  feature_key: FeatureKey;
}

export interface AuthorizeResponse {
  allowed: boolean;
  reason?: 'unauthorized' | 'tenant_membership_required' | 'entitlement_required' | 'feature_limit_exceeded';
}
