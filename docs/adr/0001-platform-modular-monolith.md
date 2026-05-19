# Platform as a modular Go monolith

The Platform owns cross-product identity, Tenants, roles, permissions, billing
primitives and Entitlements, but it will start as one modular Go monolith rather
than separate auth, tenant, billing and entitlement services. This keeps local
development, transactions and schema evolution simple while the domain is still
settling; Product backends integrate through explicit Platform APIs and shared
contracts instead of reaching into Platform tables.
