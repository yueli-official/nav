# 导航站产品

Nav 是独立的导航消费者：管理精选链接、分类、主题、推荐、搜索、点击指标和站点展示配置。
`api/` 持有领域、权限策略和数据库事实，`web/` 提供公开发现与运营管理，`research/` 只保存设计输入。

## 边界

- Nav 只依赖 `auth.oidc` 能力，不引用 Identity 的数据库、角色模型或源码包。
- Identity 只签发 Subject；Site → Category → Group → Link 的 Scope、Grant、委派和策略发布均由
  `api/internal/navauthz` 持有。
- `api/internal/runtime` 适配日志、认证、健康检查、限流、Telemetry 和配置；领域代码不感知运行平台。
- `api/internal/naverr` 是 Nav 自己的不可变公开错误契约。
- Nav 的 PostgreSQL、migration、初始站点记录、OpenAPI 和前端依赖锁均由本仓库自持。

仓库不使用 `doctor.yaml`、Platform 私有 Go 包、`@platform/*`、pnpm catalog 或 workspace link。

## 本地一键开发

本地多仓组合由相邻 Workspace 仓库统一管理，Nav 只声明自己的进程、配置和 OIDC 能力需求：

```powershell
cd ..\workspace
.\environments\nav-local\run.ps1
```

默认 Nav API 为 `8090`，Web 为 `3006`。所有端口都可通过 `LOCAL_*_PORT` 环境变量覆盖，不会要求
停止其他项目。Workspace 负责准备独立数据库、注册 OIDC client、执行 migration/bootstrap，
再按 Identity → Account → Nav API → Nav Web 启动。

## Docker Compose 部署

仓库提供两个标准入口，要求 Docker Compose `2.20.0+`：

| 入口 | Identity | Nav 数据库/API/Web |
|---|---|---|
| `compose.yaml` | 锁定版本、自动部署和注册 Nav client | 自动部署 |
| `compose.attach.yaml` | 使用已存在且已注册 client 的服务 | 自动部署 |

完整独立部署：

```powershell
Copy-Item .env.example .env
# 填写所有空 secret、管理员和数据库值
docker compose config --quiet
docker compose up -d --wait
docker compose down
```

复用已有 Identity：

```powershell
Copy-Item deploy/env/attach.env.example deploy/env/attach.env
# 填写外部 HTTPS endpoint、Nav client 和 secret
docker compose --env-file deploy/env/attach.env -f compose.attach.yaml config --quiet
docker compose --env-file deploy/env/attach.env -f compose.attach.yaml up -d --wait
docker compose --env-file deploy/env/attach.env -f compose.attach.yaml down
```

普通 `down` 保留数据库和 Identity publisher 数据；只有明确销毁实例时才执行 `down --volumes`。
该操作不可恢复。能力需求记录在 `deploy/contracts/requirements.json`，完整部署使用的基础服务版本记录在
`deploy/deployment.lock.json`。`nav-binding` 会在 API 启动前验证 issuer、discovery 和 JWKS，
attach 模式不会创建、重启或关闭外部 Identity。

生产必须使用 HTTPS origin、`NAV_COOKIE_SECURE=true` 和至少 32 字节随机的
`NAV_WEB_SEAL_SECRET`。

## 工具链与验证

- Go `1.25.12`
- Node.js `24.18.x`
- pnpm `10.28.x`

前端安装必须忽略上级 Platform workspace：

```powershell
cd web
corepack pnpm install --ignore-workspace --frozen-lockfile
```

代码测试、类型检查、生产构建、容器和 Compose 验证由本仓库 CI 自持。当前迁移阶段按统一安排暂不执行；
等所有消费者迁完后再一次性进行端到端验收。
