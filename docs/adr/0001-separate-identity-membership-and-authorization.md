# ADR 0001: Separate identity, product membership, and authorization

- Status: Accepted
- Date: 2026-08-09

## Context

Nav used Identity for login and Foundation Authorization for grants, but had no product-local member relationship. The management copy consequently described a possible role grant as user registration. Reusing grants as a member directory would omit ordinary authenticated users, while copying Identity users would create a second account system and stale profile truth.

OIDC subjects may also be pairwise per client. Nav therefore needs an explicit stable cross-service identity reference rather than assuming every `sub` is a platform public user key.

## Decision

Nav will model a `Nav Membership` keyed by Identity Public User Key.

- Anonymous reads never create membership state.
- The first authenticated Nav request idempotently creates an active membership; later requests throttle `lastSeenAt` updates.
- Identity remains authoritative for credentials, security, global account status, and public profile.
- Nav may retain a non-authoritative public-profile snapshot for resilient display and search.
- Roles, grants, applications, and policy remain in Nav Authorization and are joined into member views as summaries.
- A disabled new-member rule grants a role only when a membership is first created. It is not replayed during every permission decision.
- Suspending a membership blocks authenticated Nav operations but preserves public browsing and audit history.
- OIDC exposes the Identity Public User Key as an additive claim; consumers retain a compatibility fallback while environments migrate.

## Consequences

The admin console can list all participating users without pretending to own their accounts. Permission and membership lifecycles can be understood and operated independently. Nav adds a membership table, a small membership module, an Identity public-profile adapter, and middleware on authenticated routes. Consumers that adopt this pattern need a stable public identity key and must define their own join and suspension semantics.

Profile display can briefly lag Identity when the cached snapshot is stale or Identity is unavailable. The snapshot must never be used for authentication or authorization.

## Rejected alternatives

- Copying Identity users into Nav, because it duplicates account truth and security responsibilities.
- Treating subjects with grants or applications as members, because ordinary authenticated participants disappear.
- Creating a member on anonymous visits, because there is no stable authenticated identity and the records have no product meaning.
- Keying membership by email, handle, or display name, because those values are mutable and not authorization-safe.
