#!/bin/sh
set -eu

create_user_db() {
  db_user="$1"
  db_password="$2"
  db_name="$3"

  psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
    -v db_user="$db_user" \
    -v db_password="$db_password" \
    -v db_name="$db_name" <<'SQL'
select format('create role %I login password %L', :'db_user', :'db_password')
where not exists (select 1 from pg_roles where rolname = :'db_user')\gexec

select format('create database %I owner %I', :'db_name', :'db_user')
where not exists (select 1 from pg_database where datname = :'db_name')\gexec
SQL
}

create_user_db "codelinks_platform" "${CODELINKS_PLATFORM_DB_PASSWORD:-platform_dev_password}" "codelinks_platform"
create_user_db "codelinks_infra_link" "${CODELINKS_INFRA_LINK_DB_PASSWORD:-infra_dev_password}" "codelinks_infra_link"
create_user_db "codelinks_planer_link" "${CODELINKS_PLANER_LINK_DB_PASSWORD:-planer_dev_password}" "codelinks_planer_link"
create_user_db "codelinks_loka_link" "${CODELINKS_LOKA_LINK_DB_PASSWORD:-loka_dev_password}" "codelinks_loka_link"
