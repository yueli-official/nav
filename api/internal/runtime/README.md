# Nav API 运行时装配

本 module 只拥有 Nav API 的进程装配：环境配置、鉴权、HTTP 中间件、健康检查、PostgreSQL、OpenAPI 和遥测。

Nav 领域规则仍归 `internal/catalog`、`internal/navauthz` 等包；跨产品协议原语直接来自 Foundation。独立
Nav 不依赖 Platform 的 `gokit`，也不在这里重新发明领域 seam。
