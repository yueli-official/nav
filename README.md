# Nav product

- Lifecycle: active reusable product type
- Authority: Catalog product type `nav`, `api/` migrations/OpenAPI, `web/` UI
- Consumers: navigation site instances such as `nav-yueli`
- Verify: `pnpm platformctl verify product --file catalog/overlays/local.yaml --root . nav`

Nav owns curated link groups, topics, featured links, search and click metrics.
`api/` is the durable domain service and `web/` provides public discovery plus
management. Research notes under `research/` inform the product but are not
runtime authority.
