alter table users add column if not exists status text not null default 'active';
alter table users add column if not exists disabled_at timestamptz;
alter table users add column if not exists locked_until timestamptz;
alter table users add column if not exists last_login_at timestamptz;
alter table users add column if not exists last_login_ip inet;
alter table users add column if not exists last_login_country_code char(2);

do $$
begin
	alter table users add constraint users_status_check check (status in ('active', 'disabled', 'locked'));
exception
	when duplicate_object then null;
end $$;

create table if not exists roles (
	id text primary key,
	name text not null unique
);

create table if not exists permissions (
	id text primary key,
	name text not null unique
);

create table if not exists user_roles (
	user_id uuid not null references users (id) on delete cascade,
	role_id text not null references roles (id) on delete restrict,
	granted_at timestamptz not null default now(),
	granted_by uuid references users (id) on delete set null,
	primary key (user_id, role_id)
);

create table if not exists role_permissions (
	role_id text not null references roles (id) on delete cascade,
	permission_id text not null references permissions (id) on delete cascade,
	primary key (role_id, permission_id)
);

create table if not exists login_attempts (
	id uuid primary key default gen_random_uuid(),
	user_id uuid references users (id) on delete set null,
	email_attempted text not null,
	occurred_at timestamptz not null default now(),
	ip_address inet not null,
	ip_hash text not null,
	country_code char(2),
	city text,
	user_agent text,
	browser text,
	operating_system text,
	success boolean not null,
	failure_reason text,
	auth_method text not null,
	risk_score integer not null default 0,
	correlation_id text
);

do $$
begin
	alter table login_attempts add constraint login_attempts_failure_reason_check check (
		failure_reason is null or failure_reason in (
			'wrong_password',
			'unknown_email',
			'account_locked',
			'too_many_attempts',
			'invalid_token'
		)
	);
exception
	when duplicate_object then null;
end $$;

do $$
begin
	alter table login_attempts add constraint login_attempts_auth_method_check check (
		auth_method in ('password', 'refresh_token')
	);
exception
	when duplicate_object then null;
end $$;

create table if not exists security_events (
	id uuid primary key default gen_random_uuid(),
	user_id uuid references users (id) on delete set null,
	login_attempt_id uuid references login_attempts (id) on delete set null,
	type text not null,
	severity text not null,
	status text not null default 'open',
	summary text not null,
	metadata jsonb not null default '{}'::jsonb,
	detected_at timestamptz not null default now(),
	resolved_at timestamptz,
	resolved_by uuid references users (id) on delete set null
);

do $$
begin
	alter table security_events add constraint security_events_type_check check (
		type in (
			'many_failed_logins',
			'unusual_country',
			'many_ips_for_user',
			'locked_account_attempt',
			'many_failures_from_ip'
		)
	);
exception
	when duplicate_object then null;
end $$;

do $$
begin
	alter table security_events add constraint security_events_severity_check check (severity in ('low', 'medium', 'high', 'critical'));
exception
	when duplicate_object then null;
end $$;

do $$
begin
	alter table security_events add constraint security_events_status_check check (status in ('open', 'resolved'));
exception
	when duplicate_object then null;
end $$;

create table if not exists admin_audit_entries (
	id uuid primary key default gen_random_uuid(),
	actor_user_id uuid references users (id) on delete set null,
	action text not null,
	target_type text not null,
	target_id text not null,
	before_value jsonb,
	after_value jsonb,
	reason text,
	ip_address inet,
	created_at timestamptz not null default now()
);

create index if not exists users_status_idx on users (status);
create index if not exists user_roles_user_id_idx on user_roles (user_id);
create index if not exists login_attempts_user_id_occurred_at_idx on login_attempts (user_id, occurred_at desc);
create index if not exists login_attempts_ip_hash_occurred_at_idx on login_attempts (ip_hash, occurred_at desc);
create index if not exists login_attempts_success_occurred_at_idx on login_attempts (success, occurred_at desc);
create index if not exists security_events_status_detected_at_idx on security_events (status, detected_at desc);
create index if not exists admin_audit_entries_actor_created_at_idx on admin_audit_entries (actor_user_id, created_at desc);

insert into roles (id, name)
values
	('admin', 'Admin'),
	('support', 'Support'),
	('auditor', 'Auditor')
on conflict (id) do nothing;

insert into permissions (id, name)
values
	('admin.dashboard.read', 'Read admin dashboard'),
	('admin.users.read', 'Read users'),
	('admin.users.status.write', 'Change user status'),
	('admin.users.roles.write', 'Change user roles'),
	('admin.logs.read', 'Read login logs'),
	('admin.security.read', 'Read security events'),
	('admin.security.resolve', 'Resolve security events'),
	('admin.audit.read', 'Read admin audit log')
on conflict (id) do nothing;

insert into role_permissions (role_id, permission_id)
values
	('admin', 'admin.dashboard.read'),
	('admin', 'admin.users.read'),
	('admin', 'admin.users.status.write'),
	('admin', 'admin.users.roles.write'),
	('admin', 'admin.logs.read'),
	('admin', 'admin.security.read'),
	('admin', 'admin.security.resolve'),
	('admin', 'admin.audit.read'),
	('support', 'admin.dashboard.read'),
	('support', 'admin.users.read'),
	('support', 'admin.logs.read'),
	('support', 'admin.security.read'),
	('auditor', 'admin.dashboard.read'),
	('auditor', 'admin.logs.read'),
	('auditor', 'admin.security.read'),
	('auditor', 'admin.audit.read')
on conflict (role_id, permission_id) do nothing;
