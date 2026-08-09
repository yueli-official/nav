# Nav Domain Context

## Core terms

### Identity User

The global person and account owned by Identity. Credentials, account security, global status, and public profile truth stay in Identity.

### Identity Public User Key

The stable opaque cross-service reference for an Identity User. Nav persists this key and never uses email, handle, display name, or an OIDC pairwise subject as a membership key.

### Nav Membership

The relationship between an Identity User and Nav, created when that user first enters Nav with an authenticated session. It records product-local lifecycle and activity facts; it does not grant a role.

### Nav Member

An Identity User with a Nav Membership. Product copy calls this person a “导航成员”.

### Membership Status

The member's product-local participation state. `active` permits authenticated Nav actions. `suspended` preserves history and public browsing but blocks authenticated Nav operations until another administrator reactivates the member.

### Authorization Grant

A product-local assignment of a role and scope to a subject. A grant may exist only after membership is ensured, but membership never implies a grant.

### Authorization Application

A request for a product-local role. It is workflow state, not a member record and not a grant.

## Ownership rules

- Identity owns User and public profile truth.
- Nav owns Membership, Nav status, joined/last-seen timestamps, and Nav activity summaries.
- Nav Authorization owns roles, grants, applications, policies, and scopes.
- Anonymous catalog access owns no membership lifecycle event.

## Lifecycle rules

- The first authenticated Nav request idempotently ensures a membership.
- Subsequent authenticated requests update `lastSeenAt` with write throttling.
- A new-member automatic authorization rule runs only for the membership creation event, not on every permission decision.
- Suspending a membership does not delete Identity data, grants, applications, or audit history.
