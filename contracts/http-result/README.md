# Nav HTTP Result

From `api/`, run the Foundation Project v1 generator with
`-project ../contracts/http-result/project.json`; add `-check` for freshness.
Project paths are relative to this working directory.

The released Foundation v0.4.1 CLI does not yet implement Project v1. CI checks
canonical OpenAPI, catalog generation and compatibility separately until the
Foundation Project CLI is released. Local Project validation is an additional gate.

The catalog and structure responses are tree snapshots containing nested
categories/groups and, for the public catalog, site configuration and statistics.
They are not flat collection pages. Flat links, checks, members and tags use items.

Favicon delivery is a dedicated binary/conditional-request adapter. Its primary
result is 200 binary; matching If-None-Match returns 304 with no body. Foundation's
current operation schema rejects 304 empty as an additional success, so that
native HTTP result is verified in Nav's HTTP tests rather than misdeclared as JSON.
