insert into users (id, email, name, password_hash)
values (
	'00000000-0000-4000-8000-000000000001',
	'demo@codelinks.dev',
	'Demo User',
	'$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'
)
on conflict (email) do nothing;

insert into user_licenses (user_id, product_id)
values ('00000000-0000-4000-8000-000000000001', 'infra-link')
on conflict (user_id, product_id) do nothing;
