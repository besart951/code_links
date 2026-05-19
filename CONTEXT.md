# CodeLinks Domain Context

CodeLinks is a product family with independent applications and one shared
platform identity. This glossary keeps the business language stable across the
monorepo.

## Identity And Tenancy

**User**:
A human identity authenticated by the central Platform. A User can belong to
many Tenants.
_Avoid_: account, login when referring to the person.

**Tenant**:
The account scope that owns product access, billing state and membership. A
Tenant can represent one person, a team, a company or a customer mandate.
_Avoid_: customer, company, team when the exact legal or organizational shape is
not relevant.

**TenantMember**:
The membership connecting one User to one Tenant. The membership carries the
Role used for authorization inside that Tenant.
_Avoid_: user role when the membership relationship is meant.

**Role**:
A named set of Permissions assigned to a TenantMember. Global Roles apply across
the Tenant, while product Roles can be scoped to one Product.

**Permission**:
A named capability checked by Platform or a Product backend. Permissions answer
"may this User do this action in this Tenant?"

## Products And Access

**Product**:
An independently runnable CodeLinks application such as `infra_link`,
`planer_link` or `loka_link`.

**Product backend**:
The backend owned by one Product. It must enforce Platform Entitlements before
executing gated behavior.

**ProductAccess**:
The Platform decision that a Tenant may open a Product at all. It is derived
from active Subscriptions, manual grants or trials.
_Avoid_: paid check, app access.

**FeatureAccess**:
The Platform decision that a Tenant may use a specific feature inside a Product.
It is derived from Entitlements and FeatureLimits.
_Avoid_: frontend gate when server-side enforcement is meant.

## Billing Primitives

**Plan**:
The commercial package for one Product or bundle. A Product can have many Plans.

**Subscription**:
A Tenant's active, trialing, cancelled or expired assignment to one Plan. A
Tenant can have Subscriptions for multiple Products.

**Entitlement**:
The grant that allows a Tenant to use a Product feature. Entitlements can come
from a Subscription, a manual override or a trial.

**FeatureLimit**:
A quota or numeric limit attached to a feature grant, such as
`max_employees = 10`.

## Platform

**Platform**:
The central modular Go monolith for authentication, Tenants, RBAC, billing
primitives and Entitlements.

## Flagged Ambiguities

**Auth app vs Platform API**:
`auth.codelinks.ch` is a user-facing route/domain. The Platform API owns the
auth behavior behind it.

**Tenant vs customer/company/team**:
Use Tenant for access and billing ownership. Use company or team only when the
legal or organizational distinction matters to a Product.

## Example Dialogue

Developer: "Can this User open PlanerLink?"

Domain expert: "Check ProductAccess for the selected Tenant and `planer_link`."

Developer: "Can they export the monthly plan as Excel?"

Domain expert: "Check FeatureAccess for `planer.excel_export`, then apply any
FeatureLimit attached to the Tenant's Plan."
