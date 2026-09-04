# 统一错误与 HTTP Result 合同

## Goal

让 Nav API、application mapping 与 Nuxt 消费者采用 Foundation Project v1、声明式错误目录、标准成功 DTO 和完整结构化失败反馈，并通过真实公开导航、投稿与管理流程验收。

## Status

Open

## Current

Nav checkout 干净，当前 main 比 origin/main 多 2 个既有提交；HTTP Result 产品实现尚未修改。仓库已有 canonical OpenAPI，但前端仍有多处 `data.message` 展示，membership/authorization seam 仍从 `err.Error()` 推导公开语义，需要先清点 operation、错误目录和现有返回形状。

## Next

依据 [执行计划](plan.md) 清点 canonical OpenAPI、公开错误构造、operation 成功形状与前端失败载体，建立 Project v1 候选并运行严格生成器；不得发布版本。

## References

- [稳定上下文](context.md)
- [Foundation HTTP Result Contract](../../../../foundation/flightdeck/knowledge/errors/http-result-contract.md)
- [Foundation Error Catalog](../../../../foundation/flightdeck/knowledge/errors/error-catalog.md)
- [Foundation Failure Feedback](../../../../foundation/flightdeck/knowledge/errors/failure-feedback.md)
