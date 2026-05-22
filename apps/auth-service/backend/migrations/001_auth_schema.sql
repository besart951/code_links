create extension if not exists pgcrypto;

create table if not exists users (
	id uuid primary key default gen_random_uuid(),
	email text not null unique,
	name text not null,
	password_hash text not null,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now()
);

create table if not exists user_licenses (
	user_id uuid not null references users (id) on delete cascade,
	product_id text not null,
	granted_at timestamptz not null default now(),
	primary key (user_id, product_id),
	constraint user_licenses_product_id_check check (product_id in ('infra-link', 'planer-link', 'loka-link'))
);

create table if not exists refresh_sessions (
	token_hash text primary key,
	user_id uuid not null references users (id) on delete cascade,
	expires_at timestamptz not null,
	created_at timestamptz not null default now()
);

create index if not exists refresh_sessions_user_id_idx on refresh_sessions (user_id);
create index if not exists refresh_sessions_expires_at_idx on refresh_sessions (expires_at);
