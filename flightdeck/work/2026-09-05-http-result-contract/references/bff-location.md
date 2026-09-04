# Foundation BFF creation Location

## Reproducer

Run Nav's CLI Playwright `http-result.spec.ts` using the Shared local checkout.
The real category/group/link lifecycle reaches 201, reads raw DTOs, receives a
safe structured validation error and deletes with empty 204. Its sole failing
assertion is `group.headers().location`: the browser-facing response omits it.
Nav's group controller explicitly sets `/api/v1/nav/groups/{id}`.

## Root cause

Foundation `js/packages/nuxt-runtime/src/server/index.ts` only forwards Location
for the asset profile. The normal API profile used by Identity's inherited
`/api/v1` BFF drops the header. Current CreateBffHandlerOptions has no response
projection hook; changing all Nav API requests to asset profile would widen
redirect and download-header behavior.

## Proposed shared fix

Within Foundation createBffHandler, after copying permitted response headers:

1. For API-profile 201/202 only, read the upstream Location.
2. Accept a safe root-relative path beneath the configured target.pathPrefix.
   Reject external/protocol-relative URLs, backslashes, controls, encoded
   separators, traversal, and similar-but-distinct prefix paths.
3. Replace target.pathPrefix with mountPath. Preserve public resource/query
   semantics without disclosing an internal upstream origin.
4. Omit malformed Location without turning a successful write into failure.
5. Leave ordinary API 3xx filtering and asset-profile behavior unchanged.

Tests: 201/202 forwarding, differing mount/target prefix rewrite, invalid paths
and external origins rejected, existing API 302 filtering and Asset 302 behavior.
Then rerun Nav's real browser lifecycle and job Location tests.

This requires a Foundation-owned change. No Foundation source was modified in
this Nav session. Tag/package release remains a separate action.
