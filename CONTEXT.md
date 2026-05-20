# CodeLinks Domain Context

CodeLinks is a product family with independent applications and one shared
platform identity. This glossary keeps the business language stable across the
monorepo.

## Identity And Tenancy

**User**:
A human identity authenticated by the central Platform. A User can belong to
many Tenants.
_Avoid_: account, login when referring to the person.

**Session**:
Authenticated continuity for one User login. A Session can be revoked and is
used to derive the User for Platform requests.

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

**Product Client**:
The server-side identity of a Product backend when it asks Platform for
authorization decisions or audience-scoped Access Tokens.
_Avoid_: frontend app, browser client.

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

**AuthorizationSnapshot**:
A short-lived, audience-scoped snapshot of the authorization facts a Product
backend needs for one User in one Tenant. It can include Roles, Permissions,
ProductAccess, Entitlements, FeatureLimits and version fields, but not profile,
billing or payment data.
_Avoid_: LicenseSnapshot, user profile token.

**Audience-scoped Access Token**:
A token issued for exactly one Product backend after Platform validates the
Session, Tenant membership, ProductAccess and Product Client audience. The
receiving Product can only decrypt and use the token when its audience and key
match.
_Avoid_: global access token.

**Token Version**:
The User or Session version used to invalidate existing Access Tokens after
security-sensitive account or session changes.

**Entitlements Version**:
The Tenant and Product version used to detect stale AuthorizationSnapshots after
Role, Permission, Subscription or Entitlement changes.

## Platform

**Platform**:
The central modular Go monolith for authentication, Tenants, RBAC, billing
primitives and Entitlements.

**Superadmin**:
A Platform operator role with cross-Tenant access for support, security and
backoffice administration. Superadmin access is separate from Tenant admin
access and must be treated as privileged.
_Avoid_: owner, admin when the scope is only one Tenant.

**Admin Action**:
A privileged operation performed in the Superadmin area, such as suspending a
User, changing a Tenant setting or updating Entitlements. Admin Actions must be
audited.

**Reason**:
The human explanation supplied for a sensitive Admin Action. Reason text is part
of the audit trail and prevents silent changes.

**Audit Log**:
An append-only record of important authentication, authorization, Superadmin,
security and system events. It is searchable and filterable, but not editable.

**Global Search**:
The Superadmin search surface over safe summaries of Platform records, such as
Users, Tenants, Subscriptions, Audit Logs and Notifications. It must avoid
indexing or returning unnecessary sensitive data.

**Notification Rule**:
A typed configuration that decides when a Notification Event should be delivered
and through which channel.

**System Setting**:
A typed Platform configuration value. Sensitive System Setting changes require a
Reason and write an Audit Log entry.

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

Developer: "Can InfraLink decide from the token alone?"

Domain expert: "It may use the AuthorizationSnapshot for quick checks, but must
ask Platform again when the snapshot version is stale or the action needs a
fresh decision."

Developer: "Who can request an InfraLink token?"

Domain expert: "The request must combine the User's active Session with the
InfraLink Product Client. The request body cannot choose the User or audience."
