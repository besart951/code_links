# CodeLinks

CodeLinks is a real monorepo for multiple independent products with one shared
platform for identity, tenants, billing primitives and entitlements.

## Workspace

```text
apps/
  startseite/
  infra_link/
  planer_link/
  loka_link/
platform/
packages/
deploy/
```

The products must stay fachlich independent. Shared packages are only for stable
cross-cutting contracts, auth helpers, config and small utilities.

## Databases

The first deployment uses one PostgreSQL server with separate databases:

- `codelinks_platform`
- `codelinks_infra_link`
- `codelinks_planer_link`
- `codelinks_loka_link`

There are no cross-database foreign keys or cross-database transactions.
Platform owns users, tenants, roles, subscriptions and entitlements. Product
databases store product data with platform IDs such as `tenant_id`.

## Local Commands

```powershell
pnpm install
pnpm check
go test ./platform/...
go test ./apps/infra_link/backend/...
```

Run the integrated Docker stack:

```powershell
Copy-Item deploy/docker/.env.example deploy/docker/.env
.\deploy\scripts\dev.ps1
```

Optional development seed after the stack has created the platform schema:

```powershell
docker compose -f deploy/docker/docker-compose.yml --profile seed up platform-seed
```

## Platform API

- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/me`
- `GET /api/v1/tenants`
- `GET /api/v1/entitlements?tenant_id=...`
- `POST /internal/authorize`

Auth cookies are `HttpOnly` for access and refresh tokens. CSRF uses a readable
double-submit `csrf_token` cookie and `X-CSRF-Token` header.

## Product Gates

- `planer_link` server export routes call the platform before producing Excel
  output.
- `infra_link` can enforce module entitlements through
  `PLATFORM_ENTITLEMENTS_REQUIRED=true`.

UI gating is allowed for ergonomics, but product backends remain the enforcement
point.
