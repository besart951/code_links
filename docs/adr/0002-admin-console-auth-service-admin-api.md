# ADR 0002: Admin Console Uses Auth Service Admin API

## Status

Accepted.

## Context

The Admin Console needs to manage Users, Admin Roles, Login Attempts, Security Events, SMTP Settings, Notifications, and Admin Audit Entries. These concepts are security-sensitive and belong next to the Auth Service data model.

The SvelteKit Admin Console must stay SSR-safe. Browser code must not receive database credentials or internal service secrets.

## Decision

`apps/admin-link` is the SvelteKit Admin Console. It renders the shadcn-svelte dashboard shell and calls application use cases from server `load` functions and server actions.

The Auth Service remains the system of record for admin data. The Admin Console depends on repository interfaces, with a server-only mock repository for local UI work and an `AuthServiceAdminApiRepository` seam for the real Go Admin API.

Admin access is controlled by Admin Roles and canonical `admin.*` Permissions, not by Product Licenses. The UI can hide unavailable actions, but permission checks must also run in use cases and in the Auth Service Admin API.

The Auth Service frontend owns public auth screens (`/login`, `/signup`, `/forgot-password`, `/reset-password`) and proxies form submissions to the Go API through SvelteKit server actions. `admin-link` stays an SSR-protected Admin Console and resolves the current Admin Actor by calling `GET /api/admin/me` with the httpOnly refresh cookie.

SMTP passwords are encrypted by the Auth Service with AES-GCM using `SMTP_SECRET_KEY`; a missing key is only tolerated for local development.

## Consequences

- Admin terminology is separate from Product License terminology.
- Admin UI code stays thin and receives view models from server-side use cases.
- Sensitive IP and login metadata can be masked per permission before crossing the SSR boundary.
- SMTP writes and test-email requests create Admin Audit Entries.
- Forgot-password responses stay neutral so account existence is not leaked.
