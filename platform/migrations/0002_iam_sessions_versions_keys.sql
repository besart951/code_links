alter table users
  add column if not exists token_version integer not null default 1;

create table if not exists sessions (
  id text primary key,
  user_id uuid not null references users(id) on delete cascade,
  token_hash text not null unique,
  token_version integer not null default 1,
  user_agent text,
  ip text,
  created_at timestamptz not null default now(),
  last_seen_at timestamptz,
  revoked_at timestamptz,
  expires_at timestamptz not null
);

create index if not exists sessions_user_id_idx on sessions(user_id);
create index if not exists sessions_token_hash_idx on sessions(token_hash);
create index if not exists sessions_expires_at_idx on sessions(expires_at);

alter table refresh_tokens
  add column if not exists session_id text references sessions(id) on delete cascade;

create index if not exists refresh_tokens_session_id_idx on refresh_tokens(session_id);

create table if not exists tenant_product_versions (
  tenant_id uuid not null references tenants(id) on delete cascade,
  product_key text not null references products(key),
  entitlements_version integer not null default 1,
  updated_at timestamptz not null default now(),
  primary key (tenant_id, product_key)
);

create table if not exists audience_keys (
  audience text not null,
  kid text not null,
  public_jwk jsonb not null,
  algorithm text not null default 'ECDH-ES',
  status text not null default 'active' check (status in ('active', 'retiring', 'retired')),
  not_before timestamptz,
  not_after timestamptz,
  created_at timestamptz not null default now(),
  primary key (audience, kid)
);

create index if not exists audience_keys_active_idx
  on audience_keys(audience, status)
  where status in ('active', 'retiring');
