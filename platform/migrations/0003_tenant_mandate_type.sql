alter table tenants
  drop constraint if exists tenants_type_check;

alter table tenants
  add constraint tenants_type_check
  check (type in ('personal', 'team', 'company', 'mandate'));
