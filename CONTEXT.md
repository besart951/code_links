# Context

## Glossary

### Admin Activity Event

A time-ordered operational event shown to an Admin Actor. It can describe authentication activity, a security signal, a notification, an admin audit entry, or a runtime process log.

### Runtime Process Log

A file-backed Auth Service process log produced through Go `log.*`. Fatal process exits are persisted as `fatal:` entries so `/admin/logs` can show the crash after restart.

### Current User

The request identity a product backend reads from a validated access token. It is optimized for hot request paths and can be stale until the access token expires.

### User Info API

The Auth Service HTTP API for fresh user profile projections. Product apps call it through `packages/authuser` instead of importing Auth Service internals or reading the Auth database.

### Refresh Session

A long-lived Auth Service session backed by an opaque `refresh_token` httpOnly cookie. It can mint short-lived Access Tokens and is consumed atomically during refresh rotation.

### Admin Command

An Auth Service Admin API mutation. Admin Commands require a Bearer Access Token; refresh-cookie fallback is only for SSR Admin Actor reads.

### Login Protection

Auth Service login defence that counts recent failed password attempts by normalized email and trusted client IP hash before password comparison. When a threshold is exceeded it records `too_many_attempts`, emits a Security Event, and may set a temporary `locked_until` for a known User.

### Refresh Token Reuse

A second use of an already consumed Refresh Session token. Reuse is treated as a session theft signal: Auth Service records a Security Event, revokes active Refresh Sessions for that User, clears the refresh cookie, and returns 401.

### Product Origin Allowlist

Per-product backend CORS configuration loaded from `PRODUCT_ALLOWED_ORIGINS`. Product APIs use Bearer tokens, not cookies, so they only emit CORS allow headers for explicit Origins and do not allow credentials.

### Trusted Proxy

An upstream reverse proxy whose CIDR is listed in `TRUSTED_PROXY_CIDRS`. Auth Service trusts `X-Forwarded-For` and `X-Real-IP` only when the direct peer is a Trusted Proxy; otherwise it uses `RemoteAddr`.

### Production Config Gate

Startup validation for production deployments. It fails fast when key material, cookie security, public Origins, mock flags, or proxy trust settings would weaken Auth Service or Admin-Link security.

### Structured Auth Error

Auth Service JSON error response containing both a human-readable `error` message and a stable machine-readable `code`. Admin-Link parses these responses so loads and actions can preserve useful debug and UX context.
