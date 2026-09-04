# 统一错误与 HTTP Result 合同

## Goal

让 Nav API、application mapping 与 Nuxt 消费者采用 Foundation Project v1、声明式错误目录、标准成功 DTO 和完整结构化失败反馈，并通过真实公开导航、投稿与管理流程验收。

## Status

Finished

## Current

Nav 已实现 42-operation Project v1、声明式 catalog/operation-errors、Go/TS/i18n 生成物与迁移期 CI。领域错误由 navcause 拥有，在 application middleware 单次映射并保留 errors.Is/As；Site Profile 多字段诊断转换为 RFC 6901 pointers。集合统一 items，创建/删除 handler 显式返回 201/204，检查任务为顶层 id/status 的 202，旧 Envelope decoder 已移除。表单保留字段、未映射摘要和技术详情；review 的 Nav 自有问题已修复。

Go test/race/vet/build、govulncheck、Web 31 tests/typecheck/build 通过。显式 Vue/Vue Router 依赖及 Vite dedupe 修复本地组合的多 Vue 实例水合问题。Foundation nuxt-runtime 0.1.4 候选已修复 API BFF 的 201/202 Location 安全映射；Nav CLI Playwright 最新 5/5 通过，覆盖公开桌面/移动、后台连续刷新、真实 CRUD 201/204/DTO/Problem、创建资源和异步任务 Location 后续 GET、设置字段反馈和技术详情。针对变更的 Standards/Spec 复审无阻断发现。

Nav 使用 API 8190 / Web 3006，避免默认 8090 与另一组织 Registry 冲突；验收 Session 为 `20260904T195832Z-32916`。这次使用显式本地 Foundation checkout，不能视为已发布制品验证；正式依赖仍锁定已发布版本，0.1.4 发布及消费者制品更新移交 Workspace 发布加固 Work。

## Next

None. Foundation 候选发布与正式制品引用更新由 Workspace 发布加固 Work 承接，不属于本轮本地迁移完成门禁。

## References

- [BFF Location 修复与验证](references/bff-location.md)
- [稳定上下文](context.md)
- [Foundation HTTP Result Contract](../../../../foundation/flightdeck/knowledge/errors/http-result-contract.md)
- [Foundation Error Catalog](../../../../foundation/flightdeck/knowledge/errors/error-catalog.md)
- [Foundation Failure Feedback](../../../../foundation/flightdeck/knowledge/errors/failure-feedback.md)
