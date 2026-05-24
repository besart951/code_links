# CodeLinks Ubiquitous Language

## User

A person who can authenticate with CodeLinks and receive product licenses. A User owns licenses directly in this initial scaffold.

## Product

A licensed CodeLinks application, identified by a canonical hyphenated product ID such as `infra-link`, `planer-link`, or `loka-link`.

## Product Catalog

The generated contract listing every canonical Product ID and launch URL. TypeScript UI code and Go Product backend validation consume the Product Catalog instead of hard-coded Product IDs.

## Product License

A grant that allows a User to open and call one Product. Product Licenses are embedded into Access Tokens so Product backends can authorize requests locally.

## Access Token

A short-lived RS256 JWT issued by the central Auth Service. Product backends validate Access Tokens in memory with the Auth Service JWKS and do not query the Auth Service database.

## Refresh Session

A long-lived server-side session represented in the browser by an httpOnly refresh cookie. A Refresh Session is used to mint new Access Tokens without storing a long-lived JWT in browser JavaScript.

## Admin Console

The internal CodeLinks backoffice UI for user management, login statistics, security events, and audit review. The Admin Console is not a licensed Product; access is controlled by Admin Roles and Permissions.

## Admin Actor

A User acting inside the Admin Console. An Admin Actor has one or more Admin Roles and a derived set of Permissions for the current request.

## Admin Role

A coarse-grained role such as `admin`, `support`, or `auditor`. Roles are mapped to Permissions by the Auth Service.

## Permission

A fine-grained capability such as `admin.users.read`, `admin.users.lock`, or `admin.smtp_settings.update`. Admin Console Permissions use the canonical `admin.*` prefix. UI visibility can use Permissions, but the Auth Service and Admin API must enforce them server-side.

## Admin Access Contract

The generated contract mapping Admin Roles to canonical Admin Console Permissions. The Auth Service owns enforcement, while the Admin Console consumes the current Admin Actor Permissions returned by the Auth Service.

## Login Attempt

A recorded authentication attempt with attempted email, timestamp, IP metadata, device metadata, auth method, success/failure state, failure reason, and risk score.

## Security Event

A derived security signal created from Login Attempts or account state, such as many failed logins, unusual countries, or attempts against locked accounts.

## Admin Audit Entry

An immutable record of an Admin Actor changing or inspecting sensitive administrative state. It captures actor, action, target, before/after values, reason, IP address, and timestamp.

## SMTP Settings

The Auth Service owned mail transport configuration used for verification emails, password-reset emails, and system notifications. SMTP passwords are stored encrypted with AES-GCM and are never returned to UI code.

## Notification

A system message delivery request. In v1 the only implemented Notification Channel is `email`; in-app, webhook, and SMS remain future extensions.

## Password Reset Token

A one-time opaque token for setting a new password. Only a hash of the token is stored, the token expires after 60 minutes, and a successful reset consumes it atomically.

## Email Verification Token

A one-time opaque token proving control of a signup email address. Only a hash of the token is stored, the token expires after 24 hours, and login is blocked until the User has verified their email.
