# CodeLinks Monorepo Structure Target

This is the near-term architecture baseline for CodeLinks. CI/CD, image
registries, Kubernetes, observability stacks and complex billing webhooks are
intentionally out of scope.

## Current Assessment

The monorepo shape is already directionally correct:

- `apps/` contains the independently runnable products.
- `platform/` exists as a Go module and already owns users, tenants,
  subscriptions and entitlements at the database level.
- `packages/` contains shared TypeScript contracts and clients.
- `deploy/` contains Docker and Caddy assets.
- `go.work` already links `platform` and `apps/infra_link/backend`.
- Root `pnpm-lock.yaml` already contains `apps/startseite`,
  `apps/infra_link/frontend` and `apps/planer_link`.

The main remaining issues are migration leftovers:

- `apps/infra_link/backend/go.mod` still uses
  `github.com/besart951/go_infra_link/backend`.
- `apps/infra_link/.github`, `apps/infra_link/.codex`, app-local compose files
  and app-local lockfiles still reflect old standalone repos.
- `apps/planer_link` still contains `/besmir` naming and routing assumptions.
- `platform/` currently exposes packages at the module root instead of the
  intended `internal/<context>` structure.
- Docker has a working integrated stack, but dev/prod concerns are mixed in one
  compose file.

## Target Tree For This Phase

```text
code_links/
  apps/
    startseite/              # SvelteKit SSR marketing site
    infra_link/
      backend/               # InfraLink Go backend, independent module
      frontend/              # InfraLink SvelteKit/static frontend
    planer_link/             # PlanerLink SvelteKit/Tauri app
    loka_link/               # placeholder until product starts
  platform/
    cmd/
      platform-api/
      platform-migrate/
      platform-seed/
    internal/
      auth/
      tenants/
      users/
      permissions/
      billing/
      products/
      entitlements/
      http/
      postgres/
      config/
      tx/
    migrations/
    go.mod
  packages/
    contracts/
    auth-client/
    config/
    ui/
    utils/
  deploy/
    docker/
      docker-compose.dev.yml
      docker-compose.prod.yml
      docker-compose.override.yml
      infra-link-frontend.Dockerfile
      planer-link.Dockerfile
    caddy/
      Caddyfile.dev
      Caddyfile.prod
    env/
      platform.env.example
      startseite.env.example
      infra_link.env.example
      planer_link.env.example
      postgres.env.example
    scripts/
      dev-up.sh
      dev-down.sh
      db-reset.sh
  docs/
    adr/
    monorepo-migration.md
    monorepo-structure-target.md
  compose.yaml
  pnpm-workspace.yaml
  pnpm-lock.yaml
  go.work
```

## Platform Go Structure

Use a modular monolith with domain packages hidden under `internal/`.

```text
platform/internal/auth/
  domain.go
  ports.go
  usecase_login.go
  usecase_refresh.go
platform/internal/tenants/
  domain.go
  ports.go
  usecase_select_tenant.go
platform/internal/permissions/
  domain.go
  policy.go
platform/internal/entitlements/
  domain.go
  policy.go
  usecase_authorize.go
platform/internal/postgres/
  auth_repository.go
  tenant_repository.go
  entitlement_repository.go
  tx.go
platform/internal/http/
  server.go
  auth_handlers.go
  entitlement_handlers.go
platform/internal/config/
  config.go
platform/internal/tx/
  runner.go
```

Rules:

- Domain packages do not import HTTP, Gin, SvelteKit, GORM or pgx.
- Use cases depend on ports/interfaces owned by the domain package.
- Postgres adapters implement those ports under `internal/postgres`.
- HTTP handlers only translate request/response DTOs and call use cases.
- Transaction boundaries are explicit through a `tx.Runner` port.
- Keep `platform` on pgx/manual SQL for now; avoid adding GORM here.

## Auth And Access Flow

1. User logs in on `auth.codelinks.localhost` or `auth.codelinks.ch`.
2. Platform validates credentials and sets:
   - `access_token`: short-lived HttpOnly cookie.
   - `refresh_token`: rotating long-lived HttpOnly cookie.
   - `csrf_token`: readable double-submit cookie for unsafe requests.
3. User selects a Tenant. The chosen Tenant is carried as:
   - `codelinks_tenant_id` cookie for browser UX, and/or
   - `X-Tenant-ID` header for API calls.
4. Product frontend calls its own Product backend.
5. Product backend validates the User/session and asks Platform for
   authorization before gated behavior.
6. Platform returns a decision from Tenant membership, ProductAccess,
   FeatureAccess and optional FeatureLimits.

Keep billing decisions out of Product code. Products may cache short-lived
authorization decisions, but Platform remains the source of truth.

### Product Access

```go
func CanUseProduct(access []ProductAccess, tenantID, productKey string, now time.Time) bool {
	for _, item := range access {
		if item.TenantID == tenantID &&
			item.ProductKey == productKey &&
			item.Enabled &&
			(item.ExpiresAt == nil || item.ExpiresAt.After(now)) {
			return true
		}
	}
	return false
}
```

### Feature Access

```go
func CanUseFeature(entitlements []Entitlement, tenantID, productKey, featureKey string, now time.Time) bool {
	for _, item := range entitlements {
		if item.TenantID == tenantID &&
			item.ProductKey == productKey &&
			item.FeatureKey == featureKey &&
			item.Enabled &&
			(item.ExpiresAt == nil || item.ExpiresAt.After(now)) {
			return true
		}
	}
	return false
}
```

## Domain Data Model

The current SQL model is close. The next refinement is to distinguish product
opening from feature use.

```text
users
  id, email, password_hash, display_name, status, created_at, last_login_at

tenants
  id, type(personal|team|company), name, slug, owner_user_id, status, billing_email

tenant_members
  tenant_id, user_id, role_id, status, joined_at

roles
  id, key, name, product_key nullable

permissions
  id, key, description, product_key nullable

role_permissions
  role_id, permission_id

products
  key, name, status

plans
  id, product_key nullable, bundle_key nullable, key, name, interval, status

subscriptions
  id, tenant_id, plan_id, provider, provider_customer_id,
  provider_subscription_id, status, current_period_start, current_period_end,
  cancel_at, created_at

plan_entitlements
  plan_id, product_key, feature_key, enabled, limit_key nullable, limit_value nullable

entitlements
  tenant_id, product_key, feature_key, source, enabled, expires_at

feature_limits
  tenant_id, product_key, feature_key, limit_key, value, period, reset_at
```

Example Plan:

```text
Product: planer_link
Plan: starter
FeatureAccess:
  planer.pdf_export = true
  planer.excel_export = true
  planer.multi_user = false
FeatureLimits:
  planer.employees.max = 10
```

## Startseite SvelteKit SSR

Keep `apps/startseite` a small SvelteKit SSR site with `adapter-node`.

Recommended route shape:

```text
apps/startseite/src/routes/
  +layout.ts
  +layout.svelte
  +server.ts                     # redirects / to /de
  [locale=locale]/+page.ts
  [locale=locale]/+page.svelte
  [locale=locale]/produkte/[slug]/+page.server.ts
  [locale=locale]/produkte/[slug]/+page.svelte
  produkte/[slug]/+server.ts     # redirects legacy URLs to /de/...
  robots.txt/+server.ts
  sitemap.xml/+server.ts
apps/startseite/src/params/
  locale.ts
```

SEO baseline:

- `<title>`, description, canonical URL per route.
- OpenGraph and Twitter card tags.
- JSON-LD for Organization and SoftwareApplication where useful.
- Server-rendered product pages for CodeLinks, InfraLink, PlanerLink and later
  LokaLink.
- `robots.txt` and `sitemap.xml` generated from the same route/product data.
- German copy uses correct umlauts (`ä`, `ö`, `ü`) instead of `ae`, `oe`, `ue`
  in visible text.
- The marketing site supports German, English, French, Italian and Spanish
  copy from one typed content source.
- Language URLs use real path prefixes, not query parameters:
  `/de`, `/en`, `/fr`, `/it`, `/es`.
- Product URLs use the same locale prefix:
  `/de/produkte/infra-link`, `/en/produkte/infra-link`, and equivalent
  language variants.
- Every localized page emits a self-referencing canonical URL, all alternate
  `hreflang` URLs and an `x-default` alternate pointing at German.
- `sitemap.xml` includes every localized URL plus matching `xhtml:link`
  alternates for all supported languages.
- Language selection and light/dark mode are available in the top-right corner
  on every marketing page.
- No authenticated app state, no app shell and no pure SPA mode. Small client
  interactivity is allowed for theme selection.

## Docker Development

Local integrated stack:

```powershell
docker compose up --build
```

or explicitly:

```powershell
docker compose -f deploy/docker/docker-compose.dev.yml up --build
```

Services:

- `caddy` on port `80`.
- `startseite` on internal `3000`.
- `platform-api` on internal `8080`, local `8081`.
- `infra-backend` on internal `8080`, local `8082`.
- `infra-frontend` on internal `80`, local `5174`.
- `planer-web` on internal `3000`, local `3002`.
- `postgres` on local `5432`.

Redis is intentionally omitted for now. Refresh tokens and sessions are in
Postgres; add Redis only when rate limiting, job queues or high-write ephemeral
state become real needs.

Local domains:

```text
http://codelinks.localhost
http://auth.codelinks.localhost
http://platform.codelinks.localhost
http://api.codelinks.localhost
http://infra.codelinks.localhost
http://planer.codelinks.localhost
http://loka.codelinks.localhost
```

Use `COOKIE_DOMAIN=.codelinks.localhost` and `COOKIE_SECURE=false` locally.
Production should use `COOKIE_DOMAIN=.codelinks.ch` and `COOKIE_SECURE=true`.

## PNPM Workspace

Keep exactly one root lockfile:

```text
pnpm-lock.yaml
```

Remove these after verifying a clean install/build from root:

```powershell
Remove-Item apps/infra_link/frontend/pnpm-lock.yaml
Remove-Item apps/planer_link/pnpm-lock.yaml
```

Workspace stays:

```yaml
packages:
  - "apps/*"
  - "apps/*/frontend"
  - "packages/*"
```

Useful commands:

```powershell
pnpm install
pnpm --filter startseite dev
pnpm --filter startseite build
pnpm --filter @codelinks/infra-link-frontend dev
pnpm --filter @codelinks/planer-link dev
pnpm -r check
```

Docker builds should use the root context when the app depends on workspace
packages. Product Dockerfiles inside old app directories can remain temporarily
for standalone fallback, but integrated monorepo builds should live under
`deploy/docker/`.

## Go Workspace

Keep two Go modules for now:

```text
github.com/besart951/code_links/platform
github.com/besart951/code_links/apps/infra_link/backend
```

Target `go.work`:

```text
go 1.26.0

use (
  ./platform
  ./apps/infra_link/backend
)
```

Migration commands for InfraLink module path:

```powershell
Set-Location apps/infra_link/backend
go mod edit -module github.com/besart951/code_links/apps/infra_link/backend
$files = rg -l "github.com/besart951/go_infra_link/backend" -g "*.go"
$files | ForEach-Object {
  $content = Get-Content -Raw $_
  $content.Replace("github.com/besart951/go_infra_link/backend", "github.com/besart951/code_links/apps/infra_link/backend") | Set-Content -NoNewline $_
}
gofmt -w $files
go mod tidy
Set-Location ../../..
go work sync
go test ./apps/infra_link/backend/...
```

Do this as its own commit because it touches many files.

## Migration Cleanup Plan

Delete after verification:

- `apps/infra_link/.github/workflows/ci.yml`
- `apps/infra_link/.codex/`
- `apps/infra_link/frontend/pnpm-lock.yaml`
- `apps/planer_link/pnpm-lock.yaml`

Move or replace:

- `apps/infra_link/docker-compose.yml` and `docker-compose.prod.yml` after the
  integrated stack covers the same standalone use case.
- `apps/infra_link/frontend/Dockerfile` after the root-context Dockerfile is the
  default.
- `apps/planer_link/Dockerfile` after `deploy/docker/planer-link.Dockerfile`
  is verified.

Fix in place:

- `go_infra_link` import paths and README titles.
- `apps/planer_link` `/besmir` base path, package metadata, Tauri identifiers
  and docs if PlanerLink is no longer branded as Besmir.
- Old domains such as `cloud.besi94.ch`.

Review later:

- Bruno collections and generated Swagger docs.
- Historical migration notes that intentionally preserve old source repo names.
- Tauri bundle identifiers, because changing them can affect desktop updates.

## Next Five Steps

1. Commit the documentation and Docker/SEO baseline from this phase.
2. Run `pnpm install`, `pnpm --filter startseite build`, `pnpm -r check`.
3. Run `docker compose -f deploy/docker/docker-compose.dev.yml config`, then
   bring up the dev stack.
4. Migrate the InfraLink Go module path in one isolated commit.
5. Remove nested PNPM lockfiles and old app-local CI/Docker leftovers in a
   cleanup commit after root builds pass.

## Later

- CI/CD and image registries.
- Forgejo, GitHub Actions or Woodpecker.
- Kubernetes.
- Observability stack.
- Stripe webhook complexity.
- Microservice extraction.
- Full design system work.
