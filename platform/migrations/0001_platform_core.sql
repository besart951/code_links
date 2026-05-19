create extension if not exists pgcrypto;

create table if not exists users (
  id uuid primary key default gen_random_uuid(),
  email text not null unique,
  password_hash text not null,
  display_name text not null,
  status text not null default 'active',
  created_at timestamptz not null default now(),
  last_login_at timestamptz
);

create table if not exists refresh_tokens (
  id text primary key,
  user_id uuid not null references users(id) on delete cascade,
  token_hash text not null unique,
  user_agent text,
  ip text,
  expires_at timestamptz not null,
  created_at timestamptz not null default now(),
  revoked_at timestamptz
);

create index if not exists refresh_tokens_user_id_idx on refresh_tokens(user_id);
create index if not exists refresh_tokens_expires_at_idx on refresh_tokens(expires_at);

create table if not exists tenants (
  id uuid primary key default gen_random_uuid(),
  type text not null check (type in ('personal', 'team', 'company')),
  name text not null,
  slug text not null unique,
  owner_user_id uuid not null references users(id),
  status text not null default 'active',
  billing_email text,
  created_at timestamptz not null default now()
);

create table if not exists roles (
  id uuid primary key default gen_random_uuid(),
  key text not null unique,
  name text not null,
  product_key text
);

create table if not exists permissions (
  id uuid primary key default gen_random_uuid(),
  key text not null unique,
  description text not null default '',
  product_key text
);

create table if not exists role_permissions (
  role_id uuid not null references roles(id) on delete cascade,
  permission_id uuid not null references permissions(id) on delete cascade,
  primary key (role_id, permission_id)
);

create table if not exists tenant_members (
  tenant_id uuid not null references tenants(id) on delete cascade,
  user_id uuid not null references users(id) on delete cascade,
  role_id uuid not null references roles(id),
  status text not null default 'active',
  joined_at timestamptz not null default now(),
  primary key (tenant_id, user_id)
);

create index if not exists tenant_members_user_id_idx on tenant_members(user_id);

create table if not exists products (
  key text primary key,
  name text not null,
  status text not null default 'active'
);

create table if not exists plans (
  id uuid primary key default gen_random_uuid(),
  product_key text references products(key),
  bundle_key text,
  name text not null,
  interval text not null check (interval in ('month', 'year', 'one_time')),
  status text not null default 'active',
  check ((product_key is not null and bundle_key is null) or (product_key is null and bundle_key is not null))
);

create table if not exists subscriptions (
  id uuid primary key default gen_random_uuid(),
  tenant_id uuid not null references tenants(id) on delete cascade,
  plan_id uuid not null references plans(id),
  provider text not null default 'manual',
  provider_customer_id text,
  provider_subscription_id text,
  status text not null,
  current_period_start timestamptz,
  current_period_end timestamptz,
  cancel_at timestamptz,
  created_at timestamptz not null default now()
);

create index if not exists subscriptions_tenant_id_idx on subscriptions(tenant_id);

create table if not exists entitlements (
  tenant_id uuid not null references tenants(id) on delete cascade,
  product_key text not null references products(key),
  feature_key text not null,
  source text not null check (source in ('subscription', 'manual', 'trial')),
  enabled boolean not null default true,
  expires_at timestamptz,
  created_at timestamptz not null default now(),
  primary key (tenant_id, product_key, feature_key, source)
);

create index if not exists entitlements_lookup_idx
  on entitlements(tenant_id, product_key, feature_key)
  where enabled = true;

create table if not exists feature_limits (
  tenant_id uuid not null references tenants(id) on delete cascade,
  product_key text not null references products(key),
  feature_key text not null,
  limit_key text not null,
  value bigint not null,
  period text not null default 'none' check (period in ('none', 'day', 'month', 'year')),
  reset_at timestamptz,
  primary key (tenant_id, product_key, feature_key, limit_key)
);

insert into products (key, name, status) values
  ('infra_link', 'infra_link', 'active'),
  ('planer_link', 'planer_link', 'active'),
  ('loka_link', 'loka_link', 'active')
on conflict (key) do update set name = excluded.name, status = excluded.status;

insert into roles (key, name, product_key) values
  ('owner', 'Owner', null),
  ('admin', 'Admin', null),
  ('member', 'Member', null)
on conflict (key) do nothing;

insert into permissions (key, description, product_key) values
  ('tenant.manage', 'Manage tenant settings and members', null),
  ('billing.manage', 'Manage tenant plans and subscriptions', null),
  ('planer.pdf_export', 'Use planer_link PDF export', 'planer_link'),
  ('planer.excel_export', 'Use planer_link Excel export', 'planer_link'),
  ('infra.module_bacnet', 'Use infra_link BACnet module', 'infra_link'),
  ('infra.module_sps', 'Use infra_link SPS module', 'infra_link'),
  ('infra.module_field_devices', 'Use infra_link field device module', 'infra_link')
on conflict (key) do nothing;
