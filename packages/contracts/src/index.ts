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

export type TenantType = 'personal' | 'team' | 'company' | 'mandate';
export type TenantStatus = 'active' | 'inactive' | 'suspended' | 'archived' | 'disabled';
export type UserStatus = 'active' | 'inactive' | 'locked' | 'disabled';
export type SubscriptionStatus = 'active' | 'trialing' | 'canceled' | 'expired';
export type RiskStatus = 'normal' | 'watch' | 'high';
export type SortDirection = 'asc' | 'desc';

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

export interface PageRequest {
  limit?: number;
  offset?: number;
  sort?: string;
  direction?: SortDirection;
}

export interface PageMeta {
  limit: number;
  offset: number;
  total: number;
}

export interface PageResponse<T> {
  items: T[];
  page: PageMeta;
}

export interface AdminUserSummary {
  id: string;
  email: string;
  display_name: string;
  status: UserStatus;
  email_verified: boolean;
  mfa_enabled: boolean;
  last_login_at: string | null;
  failed_login_count: number;
  locked_until: string | null;
  created_at: string;
  tenant_count: number;
  active_sessions: number;
}

export interface AdminTenantSummary {
  id: string;
  name: string;
  tenant_type: TenantType;
  status: TenantStatus;
  created_at: string;
  updated_at: string | null;
  owner_user_id: string;
  billing_email: string | null;
  country: string | null;
  locale: string | null;
  timezone: string | null;
  active_products: ProductKey[];
  subscription_status: SubscriptionStatus | 'none';
  risk_status: RiskStatus;
}

export interface AdminDashboardMetric {
  key: string;
  label: string;
  value: number;
  tone: 'neutral' | 'success' | 'warning' | 'danger' | 'info';
}

export interface AdminProductSummary {
  product_key: ProductKey;
  name: string;
  active_tenants: number;
  active_users: number;
  active_subscriptions: number;
  warning_count: number;
  last_access_at: string | null;
}

export interface AdminDashboardSummary {
  metrics: AdminDashboardMetric[];
  products: AdminProductSummary[];
  security_warnings: number;
  open_system_messages: number;
  generated_at: string;
}

export type SearchResultType =
  | 'tenant'
  | 'user'
  | 'product'
  | 'subscription'
  | 'audit_log'
  | 'notification';

export interface AdminSearchResult {
  type: SearchResultType;
  id: string;
  title: string;
  subtitle: string;
  matched_fields: string[];
  rank: number;
}

export interface AdminSearchResponse {
  query: string;
  results: AdminSearchResult[];
  facets: Record<SearchResultType, number>;
  page: PageMeta;
}

export interface AuditLogEntry {
  id: string;
  tenant_id: string | null;
  actor_user_id: string;
  target_type: string;
  target_id: string;
  action: string;
  reason: string | null;
  ip_address: string | null;
  user_agent: string | null;
  created_at: string;
  metadata: Record<string, unknown>;
}

export type NotificationChannel = 'in_app' | 'email' | 'webhook' | 'sms';

export interface NotificationTemplateSummary {
  id: string;
  key: string;
  channel: NotificationChannel;
  subject: string;
  enabled: boolean;
  updated_at: string;
}

export interface NotificationDeliverySummary {
  id: string;
  event_key: string;
  channel: NotificationChannel;
  status: 'queued' | 'sent' | 'failed' | 'retrying';
  recipient: string;
  created_at: string;
  last_attempt_at: string | null;
}

export interface SecurityEventSummary {
  id: string;
  event_type: string;
  severity: 'info' | 'warning' | 'critical';
  user_id: string | null;
  tenant_id: string | null;
  ip_address: string | null;
  created_at: string;
  summary: string;
}

export interface AdminSetting {
  key: string;
  label: string;
  value: string | number | boolean;
  value_type: 'string' | 'number' | 'boolean' | 'duration';
  sensitive: boolean;
  requires_reason: boolean;
  updated_at: string;
}

export interface AdminMeResponse {
  user: AuthUser;
  permissions: string[];
  superadmin: boolean;
}
