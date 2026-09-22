# Nav 开发者令牌

## Goal

让用户可以在 Account 创建绑定 `nav-yueli-web` 的细粒度开发者令牌，通过 Nav 现有 API 读取公开目录，并按账号当前 Nav 权限提交、编辑和治理链接，以及管理结构、健康检查和站点设置。每次请求都重新验证 PAT，并取令牌 scope 与账号当前 Nav 权限的交集。

完成时应满足：普通内容维护者只能使用自己当前拥有的链接提交/编辑能力；管理员能力不会授予普通用户；成员、授权控制台、角色与其他未声明路由不向 PAT 开放；显式 PAT 失败不回退浏览器 Cookie；开发者指南、OpenAPI、运行配置与真实本地浏览器/API 验收一致。

## Status

待查收

## Current

2026-09-22 本地实现与真实组合验收已完成。Nav 已升级到 Foundation Go `v0.5.0`，接入在线 PAT verifier、可信权限目录、显式路由 allowlist 和当前账号权限交集。Account 当前管理员目录精确展示六项 Nav 权限；链接管理集合查询也同时受 `nav.link.update` 的令牌 scope 约束。成员管理、授权控制台、角色工作流、`/api/v1/me`、内部权限目录和未声明路由保持拒绝。

Nav 现有 v1 授权目录已经完整表达这些能力，因此没有升级授权目录或增加数据库授权迁移。验证通过：`GOWORK=off go test -timeout 60s ./...`、全量 `go vet ./...` 和 43 项 OpenAPI/HTTP Result freshness。隔离 Session `20260922T064758Z-21196` 四进程 ready；CLI Playwright 1/1 通过，覆盖 Account UI 创建六权限令牌、桌面/390px 手机目录、公开目录、草稿链接创建/查询/编辑/删除、健康检查与设置读取、敏感端点 403，以及撤销后 Nav 401。截图已目检且无横向溢出，临时链接和令牌已清理，组合验收后已停止。

本 Work 只负责 Nav，本地实现与验收，不部署生产。Gallery 与 Paste 已在各自独立 Work 完成本地交付并等待用户确认；Distribution 后续建立独立单站 Work。

## Next

等待用户查收 Nav 开发者令牌交付；收到反馈后处理对应问题。2026-09-22 用户已授权本地提交，尚未推送或生产部署。

## References

- [上下文与权限边界](context.md)
- [实施与验收计划](plan.md)
- [Foundation 个人令牌合同](../../../../foundation/flightdeck/knowledge/authorization/personal-access-tokens.md)
