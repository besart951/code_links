# ADR 0003: Generated Contracts And Test Schema

## Status

Accepted.

## Context

Product IDs, Admin Roles, and Admin Permissions are used by TypeScript UI code, Go services, SQL seed data, and Docker topology. These strings are security-sensitive and license-sensitive.

The Auth Service schema now evolves through migrations, and the Postgres store is the production path.

## Decision

`packages/config/*.json` is the source of truth for the Product Catalog and Admin Access Contract.

`pnpm generate:contracts` generates TypeScript exports, Go Product Catalog code, and Go Admin Access code. Contract checks fail when generated files drift from JSON.

Auth Service tests include memory store coverage and a compose-backed Postgres store contract through `pnpm go:test:postgres`. The root `pnpm test` runs the Postgres path by default.

`/api/licenses/mock-purchase` is a dev/test route only. It is registered only when `ENABLE_MOCK_PURCHASE=true`.

## Consequences

- Product IDs and Admin Permissions have one edit point.
- SQL seed/check constraints are parity-tested against contracts.
- Local test runs require Docker for the default Postgres store contract.
- Demo license grants cannot appear in production by default.
