create extension if not exists pg_trgm;
create extension if not exists unaccent;

alter table users add column if not exists email_verified boolean not null default false;
alter table users add column if not exists mfa_enabled boolean not null default false;
alter table users add column if not exists failed_login_count integer not null default 0;
alter table users add column if not exists locked_until timestamptz;

alter table tenants add column if not exists updated_at timestamptz;
alter table tenants add column if not exists country text;
alter table tenants add column if not exists locale text;
alter table tenants add column if not exists timezone text;

create table if not exists platform_admin_permissions (
  user_id uuid not null references users(id) on delete cascade,
  permission_key text not null,
  granted_at timestamptz not null default now(),
  granted_by uuid references users(id),
  revoked_at timestamptz,
  primary key (user_id, permission_key)
);

create index if not exists platform_admin_permissions_active_idx
  on platform_admin_permissions (permission_key, user_id)
  where revoked_at is null;

insert into permissions (key, description, product_key) values
  ('platform.admin.read', 'Read Platform Superadmin backoffice data', null),
  ('platform.admin.write', 'Mutate Platform Superadmin backoffice data', null)
on conflict (key) do nothing;

create table if not exists audit_logs (
  id uuid primary key default gen_random_uuid(),
  tenant_id uuid references tenants(id),
  actor_user_id uuid not null references users(id),
  target_type text not null,
  target_id text not null,
  action text not null,
  reason text,
  ip_address text,
  user_agent text,
  metadata jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now()
);

create index if not exists audit_logs_created_at_idx on audit_logs (created_at desc);
create index if not exists audit_logs_actor_idx on audit_logs (actor_user_id, created_at desc);
create index if not exists audit_logs_target_idx on audit_logs (target_type, target_id, created_at desc);

create table if not exists security_events (
  id uuid primary key default gen_random_uuid(),
  event_type text not null,
  severity text not null check (severity in ('info', 'warning', 'critical')),
  user_id uuid references users(id),
  tenant_id uuid references tenants(id),
  ip_address text,
  summary text not null,
  metadata jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now(),
  resolved_at timestamptz
);

create index if not exists security_events_open_idx
  on security_events (severity, created_at desc)
  where resolved_at is null;

create table if not exists admin_settings (
  key text primary key,
  label text not null,
  value_json jsonb not null,
  value_type text not null check (value_type in ('string', 'number', 'boolean', 'duration')),
  sensitive boolean not null default false,
  requires_reason boolean not null default true,
  updated_at timestamptz not null default now(),
  updated_by uuid references users(id)
);

insert into admin_settings (key, label, value_json, value_type, sensitive, requires_reason) values
  ('platform.name', 'Platform name', '"CodeLinks"', 'string', false, true),
  ('registration.enabled', 'Registration enabled', 'true', 'boolean', false, true),
  ('maintenance.enabled', 'Maintenance mode', 'false', 'boolean', false, true),
  ('session.lifetime', 'Session lifetime', '"720h"', 'duration', false, true),
  ('access_token.lifetime', 'Access token lifetime', '"10m"', 'duration', false, true)
on conflict (key) do nothing;

create table if not exists notification_templates (
  id uuid primary key default gen_random_uuid(),
  key text not null,
  channel text not null check (channel in ('in_app', 'email', 'webhook', 'sms')),
  subject text not null default '',
  body text not null default '',
  enabled boolean not null default true,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (key, channel)
);

create table if not exists notification_deliveries (
  id uuid primary key default gen_random_uuid(),
  event_key text not null,
  channel text not null check (channel in ('in_app', 'email', 'webhook', 'sms')),
  status text not null check (status in ('queued', 'sent', 'failed', 'retrying')),
  recipient text not null,
  attempts integer not null default 0,
  created_at timestamptz not null default now(),
  last_attempt_at timestamptz,
  metadata jsonb not null default '{}'::jsonb
);

create index if not exists notification_deliveries_status_idx
  on notification_deliveries (status, created_at desc);

create table if not exists search_documents (
  id uuid primary key default gen_random_uuid(),
  source_type text not null,
  source_id text not null,
  title text not null,
  subtitle text not null default '',
  body text not null default '',
  search_vector tsvector generated always as (
    setweight(to_tsvector('simple', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('simple', coalesce(subtitle, '')), 'B') ||
    setweight(to_tsvector('simple', coalesce(body, '')), 'C')
  ) stored,
  updated_at timestamptz not null default now(),
  unique (source_type, source_id)
);

create index if not exists search_documents_vector_idx
  on search_documents using gin (search_vector);

create index if not exists search_documents_title_trgm_idx
  on search_documents using gin (title gin_trgm_ops);
