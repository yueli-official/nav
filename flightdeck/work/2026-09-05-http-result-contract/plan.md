# Plan

## P0 — 合同

- [x] 建立 v1 error catalog、operation-errors、Project v1 与生成物。
- [ ] Nav API items/201/202/204/Location 与 favicon 304 已实现并验证；BFF 丢失 Location 待 Foundation 修复。

## P1 — 实现与消费者

- [x] 建立 typed cause → application catalog mapping → HTTP projection seam。
- [x] 更新 Foundation 依赖、Nuxt 结构化 feedback 与迁移期 CI freshness/compatibility 门禁（Project 单命令等待 Foundation 发布）。

## P2 — 验证

- [x] Go test/race/vet/govuln 与 Web tests/typecheck/build 全绿。
- [ ] Workspace local checkout 与 CLI Playwright 完成公开导航、投稿、管理和失败反馈验收。
- [ ] 双轴 review Standards / Spec 无阻断发现。
