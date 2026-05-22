alter table users add column if not exists email_verified_at timestamptz;

create table if not exists password_reset_tokens (
	token_hash text primary key,
	user_id uuid not null references users (id) on delete cascade,
	expires_at timestamptz not null,
	used_at timestamptz,
	created_at timestamptz not null default now()
);

create table if not exists email_verification_tokens (
	token_hash text primary key,
	user_id uuid not null references users (id) on delete cascade,
	expires_at timestamptz not null,
	used_at timestamptz,
	created_at timestamptz not null default now()
);

create table if not exists smtp_settings (
	id text primary key default 'default',
	host text not null default '',
	port integer not null default 587,
	username text not null default '',
	password_encrypted text,
	encryption text not null default 'starttls',
	from_email text not null default 'no-reply@codelinks.dev',
	from_name text not null default 'CodeLinks',
	reply_to_email text not null default 'support@codelinks.dev',
	active boolean not null default false,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now(),
	constraint smtp_settings_encryption_check check (encryption in ('none', 'ssl', 'tls', 'starttls'))
);

create table if not exists notifications (
	id uuid primary key default gen_random_uuid(),
	type text not null,
	channel text not null,
	recipient text not null,
	subject text not null,
	status text not null,
	created_at timestamptz not null default now(),
	sent_at timestamptz,
	constraint notifications_channel_check check (channel in ('email')),
	constraint notifications_status_check check (status in ('queued', 'pending', 'sent', 'failed'))
);

create table if not exists notification_deliveries (
	id uuid primary key default gen_random_uuid(),
	notification_id uuid not null references notifications (id) on delete cascade,
	channel text not null,
	status text not null,
	provider_message_id text,
	error_message text,
	created_at timestamptz not null default now(),
	delivered_at timestamptz
);

create index if not exists password_reset_tokens_user_id_idx on password_reset_tokens (user_id);
create index if not exists password_reset_tokens_expires_at_idx on password_reset_tokens (expires_at);
create index if not exists email_verification_tokens_user_id_idx on email_verification_tokens (user_id);
create index if not exists email_verification_tokens_expires_at_idx on email_verification_tokens (expires_at);
create index if not exists notifications_created_at_idx on notifications (created_at desc);

alter table login_attempts drop constraint if exists login_attempts_failure_reason_check;
alter table login_attempts add constraint login_attempts_failure_reason_check check (
	failure_reason is null or failure_reason in (
		'wrong_password',
		'unknown_email',
		'account_locked',
		'too_many_attempts',
		'invalid_token',
		'email_not_confirmed'
	)
);

alter table security_events drop constraint if exists security_events_type_check;
alter table security_events add constraint security_events_type_check check (
	type in (
		'many_failed_logins',
		'unusual_country',
		'many_ips_for_user',
		'locked_account_attempt',
		'locked_account_login',
		'many_failures_from_ip',
		'new_ip',
		'new_country',
		'password_reset_spam'
	)
);

insert into roles (id, name)
values
	('user', 'User'),
	('admin', 'Admin'),
	('support', 'Support'),
	('auditor', 'Auditor')
on conflict (id) do nothing;

insert into permissions (id, name)
values
	('admin.dashboard.read', 'Read admin dashboard'),
	('admin.users.read', 'Read users'),
	('admin.users.update', 'Update users'),
	('admin.users.lock', 'Lock and unlock users'),
	('admin.users.change_role', 'Change user roles'),
	('admin.auth_logs.read', 'Read authentication logs'),
	('admin.security_events.read', 'Read security events'),
	('admin.smtp_settings.read', 'Read SMTP settings'),
	('admin.smtp_settings.update', 'Update SMTP settings'),
	('admin.notifications.read', 'Read notifications'),
	('admin.audit_entries.read', 'Read admin audit entries')
on conflict (id) do nothing;

insert into role_permissions (role_id, permission_id)
values
	('admin', 'admin.dashboard.read'),
	('admin', 'admin.users.read'),
	('admin', 'admin.users.update'),
	('admin', 'admin.users.lock'),
	('admin', 'admin.users.change_role'),
	('admin', 'admin.auth_logs.read'),
	('admin', 'admin.security_events.read'),
	('admin', 'admin.smtp_settings.read'),
	('admin', 'admin.smtp_settings.update'),
	('admin', 'admin.notifications.read'),
	('admin', 'admin.audit_entries.read'),
	('support', 'admin.dashboard.read'),
	('support', 'admin.users.read'),
	('support', 'admin.users.update'),
	('support', 'admin.users.lock'),
	('support', 'admin.auth_logs.read'),
	('support', 'admin.security_events.read'),
	('support', 'admin.notifications.read'),
	('auditor', 'admin.dashboard.read'),
	('auditor', 'admin.auth_logs.read'),
	('auditor', 'admin.security_events.read'),
	('auditor', 'admin.notifications.read'),
	('auditor', 'admin.audit_entries.read')
on conflict (role_id, permission_id) do nothing;

update users
set email_verified_at = coalesce(email_verified_at, now())
where email = 'demo@codelinks.dev';

insert into user_roles (user_id, role_id)
select id, 'admin'
from users
where email = 'demo@codelinks.dev'
on conflict (user_id, role_id) do nothing;
