alter table users add column if not exists locked_until timestamptz;
alter table refresh_sessions add column if not exists consumed_at timestamptz;

create index if not exists refresh_sessions_user_active_idx
	on refresh_sessions (user_id)
	where consumed_at is null;

create index if not exists login_attempts_email_failed_occurred_at_idx
	on login_attempts (lower(email_attempted), occurred_at desc)
	where success = false;

create index if not exists login_attempts_ip_failed_occurred_at_idx
	on login_attempts (ip_hash, occurred_at desc)
	where success = false;

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
		'password_reset_spam',
		'login_rate_limited',
		'refresh_token_reuse'
	)
);
