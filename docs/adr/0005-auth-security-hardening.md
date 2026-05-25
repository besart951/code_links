# Auth Security Hardening

Accepted.

Auth Service owns Login Protection, Refresh Token Reuse handling, Trusted Proxy resolution, Structured Auth Errors, and the Production Config Gate. Product backends own their Product Origin Allowlist and consume Auth Service keys through a rotating JWKS cache.

## Decisions

- Login Protection is a deep Auth Service module with balanced defaults: 5 failed attempts per normalized email, 10 failed attempts per trusted client IP hash, a 15 minute window, and a 15 minute temporary lockout for known Users.
- Refresh Sessions are no longer deleted on refresh. They are atomically marked consumed so a second use can be detected as Refresh Token Reuse. On reuse, Auth Service records a Security Event, revokes active Refresh Sessions for that User, clears the cookie, and returns 401.
- Product APIs do not reflect arbitrary Origins. They load `PRODUCT_ALLOWED_ORIGINS` as a CSV allowlist and omit `Access-Control-Allow-Credentials` because product calls use Bearer tokens, not cookies.
- Product token validators refresh JWKS on unknown `kid` and after cache TTL expiry. JWKs must be RSA signing keys and must not contradict RS256 token validation.
- Auth Service trusts forwarding headers only when the direct peer is inside `TRUSTED_PROXY_CIDRS`. Empty trust means no forwarding headers are trusted.
- Auth Service JSON errors include `error` and `code`. Admin-Link parses Structured Auth Errors so loads and actions can preserve status, code, and message.
- Production startup fails if JWT key material, SMTP secret, secure cookie settings, public Origins, mock flags, or Trusted Proxy CIDRs are unsafe.

## Consequences

The Store interface gains explicit methods for failed-login counting, temporary User locks, Refresh Session reuse/revocation, and Security Event recording. Product backend startup now requires explicit browser Origins when CORS is needed. Admin-Link routes can show cleaner failure messages without depending on raw `Admin API request failed: status` strings.

This extends ADR 0001: Refresh Sessions remain opaque cookie-backed sessions, but rotation now keeps consumed markers for theft detection, and product validators no longer require restart for normal JWKS rotation. This extends ADR 0002: Admin API errors are part of the Admin-Link/Auth Service contract and carry stable machine codes.
