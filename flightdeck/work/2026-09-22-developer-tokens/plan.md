# 实施与验收计划

- [x] 产品认证与授权：升级 Foundation Go，接入 PAT verifier、可信权限目录、当前权限交集与显式路由 allowlist。
- [x] 集合查询边界：让链接管理查询计划同时受 `nav.link.update` 的令牌 scope 约束。
- [x] 产品合同与文档：更新 OpenAPI/HTTP Result 合同、开发者令牌指南和 README 入口。
- [x] Workspace 接线：登记 Identity PAT application 与 Nav verifier，不影响其他 Target。
- [x] 验证：相关 Go 测试、合同 freshness，以及真实 Workspace 组合中的桌面/手机 Account 权限目录、Nav API 闭环、敏感端点拒绝与撤销失效。
