# ADR 0001: Central Auth With Local RS256 Product Validation

## Status

Accepted

## Context

CodeLinks has multiple product applications that need a shared login and license model. Product backends should avoid calling the central auth database on every request, but they still need a trustworthy way to authorize product access.

## Decision

Use a central Auth Service for login, refresh sessions, license grants, and RS256 Access Token issuance. Product backends fetch the Auth Service JWKS, cache the public key in memory, validate incoming Bearer tokens locally, and check their own product ID in the token `licenses` claim. Refresh sessions are opaque, httpOnly cookie-backed sessions and are atomically consumed during refresh rotation.

## Consequences

This keeps product request latency low and prevents product services from depending on the auth database. Access Tokens cannot be revoked immediately once issued, so they must stay short-lived and be refreshed through httpOnly Refresh Sessions.
