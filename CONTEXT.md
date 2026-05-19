# CodeLinks Domain Context

## Core Terms

- User: a human identity authenticated by the central platform.
- Tenant: an account scope. A tenant can be a personal account, team or company.
- TenantMember: the membership connecting a user to a tenant with a role.
- Product: an independently deployable CodeLinks application such as
  `infra_link`, `planer_link` or `loka_link`.
- Subscription: a tenant's plan for one product or bundle.
- Entitlement: the platform decision that a tenant may use a product feature.
- FeatureLimit: a quota or numeric limit attached to an entitlement.

## Architecture Terms

- Platform: the central modular Go monolith for auth, tenants, RBAC, billing
  primitives and entitlements.
- Product backend: the backend owned by one product. It must enforce platform
  entitlements before executing gated behavior.
