# Platform as a modular Go monolith

The Platform owns cross-product identity, Tenants, roles, permissions, billing
primitives and Entitlements, but it will start as one modular Go monolith rather
than separate auth, tenant, billing and entitlement services. Platform Go code
lives under `platform/internal` so Products cannot import Platform internals.
Product backends and the Superadmin UI integrate through explicit Platform HTTP
APIs and shared contracts instead of Go imports or direct database access. This
keeps local development, transactions and schema evolution simple while the
domain is still settling.
