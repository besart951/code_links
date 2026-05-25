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
