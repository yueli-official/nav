# 统一错误与 HTTP Result 合同

## Goal

让 Nav API、application mapping 与 Nuxt 消费者采用 Foundation Project v1、声明式错误目录、标准成功 DTO 和完整结构化失败反馈，并通过真实公开导航、投稿与管理流程验收。

## Status

Open

## Current

Nav 已实现 42-operation Project v1、声明式 catalog/operation-errors、Go/TS/i18n 生成物与迁移期 CI。领域错误由 navcause 拥有，在 application middleware 单次映射并保留 errors.Is/As；Site Profile 多字段诊断转换为 RFC 6901 pointers。集合统一 items，创建/删除 handler 显式返回 201/204，检查任务为顶层 id/status 的 202，旧 Envelope decoder 已移除。表单保留字段、未映射摘要和技术详情；review 的 Nav 自有问题已修复。

Go test/race/vet/build、govulncheck、Web 31 tests/typecheck 通过；Web 生产构建通过。显式 Vue/Vue Router 依赖及 Vite dedupe 修复本地组合的多 Vue 实例水合问题。CLI Playwright 最新结果 3 passed / 1 failed：公开桌面/移动、后台刷新及页面、设置字段反馈通过；真实 CRUD 201/204/DTO/Problem 均通过，仅 group 创建响应的 Location 被 Foundation API BFF 过滤（保留 soft assertion，以便仍验证后续真实写入和清理）。

Nav 使用 API 8190 / Web 3006，避免默认 8090 与另一组织 Registry 冲突。基础 Provider 未停止；完成最终构建后须通过 Workspace CLI 查询 Nav Session 的实时状态。

## Next

取得跨仓修复授权后，按 [BFF Location 修复方案](references/bff-location.md) 修复 Foundation 201/202 的安全 Location 透传；然后重建 Nav local 组合并通过完整 Playwright。保留本 Work Open，不能将当前 3/4 判为验收通过。Foundation Project v1 正式发布与制品升级仍按上游 Work 单独授权。

## References

- [稳定上下文](context.md)
- [Foundation HTTP Result Contract](../../../../foundation/flightdeck/knowledge/errors/http-result-contract.md)
- [Foundation Error Catalog](../../../../foundation/flightdeck/knowledge/errors/error-catalog.md)
- [Foundation Failure Feedback](../../../../foundation/flightdeck/knowledge/errors/failure-feedback.md)
