# Product

<!-- impeccable:product-schema 1 -->
<!-- Inferred from the repository and the user-approved membership brief on 2026-08-09. -->

## Platform

web

## Users

- Visitors use the public catalog to discover curated websites without signing in.
- Authenticated contributors submit and maintain navigation content.
- Administrators operate the catalog, members, and product-local authorization policy.

## Product Purpose

月离导航 turns a deliberately curated set of websites into a calm, task-oriented working surface. It should help people reach useful resources quickly while keeping governance understandable to administrators.

## Positioning

Nav is a consumer product in the 月离 platform. Identity owns the person and account; Nav owns the person's participation in Nav; product authorization owns roles, grants, applications, and policy.

## Operating Context

The public catalog is optimized for frequent browsing. The authenticated management console is an operational tool used on desktop and narrow screens. Account and security tasks remain in the central Identity experience.

## Capabilities and Constraints

- Anonymous browsing never creates a Nav member.
- The first authenticated Nav session creates or refreshes an idempotent Nav membership.
- A membership is not a credential, global profile, role, grant, or application.
- Identity public user keys are the stable cross-service reference.
- Management collections use the existing shared collection shell and keep controls, column headers, and pagination stable while rows scroll.
- Frontend acceptance is performed with CLI Playwright against the real OIDC journey.

## Brand Commitments

The Chinese interface is concise, calm, direct, and operational. It favors clear relationships and explicit consequences over promotional language.

## Evidence on Hand

- Existing public catalog and management routes.
- Existing Identity OIDC integration and public-user API.
- Existing Foundation authorization and shared collection components.
- Seeded local catalog data and the current Nav acceptance environment.

## Product Principles

- One concept has one owner and one name.
- Public reading stays frictionless; authenticated participation is explicit.
- Membership, permission, and application state are shown separately.
- Administrative actions expose scope and consequence before mutation.
- Empty, filtered-empty, unavailable, and loading states are distinct.

## Accessibility & Inclusion

Use semantic controls and tables, visible keyboard focus, meaningful status text that does not rely on color, responsive layouts without hidden required actions, and concise recovery actions for errors and empty results.
