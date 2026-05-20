# Superadmin UI and shared Svelte UI package

The Superadmin surface is implemented as a separate SvelteKit SSR app in
`apps/admin`. It is a high-risk backoffice interface, so it gets its own route
protection, server-side Platform fetches and deployment boundary while still
living in the monorepo.

Reusable shadcn-svelte primitives live in `packages/ui` and are consumed through
subpath exports such as `@codelinks/ui/button` and `@codelinks/ui/table`.
Product-specific components remain inside their Product apps. A component only
moves into `packages/ui` when it is a neutral primitive or when at least two apps
need the same pattern.

The Admin UI talks to Platform through HTTP APIs under `/api/v1/admin/...` and
shared TypeScript contracts from `packages/contracts`. It does not import
Platform Go packages, query Platform tables directly or use Product JWE token
exchange. The browser uses the Platform HttpOnly session cookie, and Admin pages
load data through server-side SvelteKit code so cookies and secrets stay
server-side.
