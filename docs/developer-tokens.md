# 开发者令牌 API

Nav 复用 Identity 账户中心的个人访问令牌（PAT）。令牌只通过 `Authorization: Bearer pat_…` 使用，不换取浏览器 Cookie；每次请求都检查“令牌勾选权限 ∩ 当前账号在 Nav 的实时权限”。撤销令牌或撤销 Nav 角色后，下次请求立即失效。

## 权限

| 权限键 | 用途 |
| --- | --- |
| `nav.link.submit` | 创建链接草稿；发布或推荐还需治理链接权限 |
| `nav.link.update` | 读取可管理链接并编辑当前账号授权范围内的链接 |
| `nav.link.moderate` | 发布、归档、删除或批量治理链接 |
| `nav.structure.manage` | 读取和管理分类、分组与标签 |
| `nav.health_check.run` | 读取健康状态、运行检查与维护豁免 |
| `nav.settings.manage` | 读取和修改 Nav 公开站点设置 |

Account 只展示当前账号实际拥有的权限。内容维护者的链接编辑继续受授权范围和“仅自己的链接”关系约束；令牌 scope 不会把管理员权限授予普通账号。成员管理、角色申请/邀请/授予、授权控制台、当前浏览器会话和其他未声明路由不向 PAT 开放。

## 提交链接

API 完整定义见 [OpenAPI](../contracts/openapi/nav.json)。下面示例中的令牌由调用者放在当前终端，不能写入仓库。

```powershell
$nav = 'http://127.0.0.1:8090'
$headers = @{ Authorization = "Bearer $env:NAV_TOKEN" }

# 管理员可用结构权限读取有效 categoryId/groupId；内容维护者应使用已获授权的分组。
$structure = Invoke-RestMethod "$nav/api/v1/admin/nav/structure" -Headers $headers
$category = $structure.categories | Where-Object { $_.groups.Count -gt 0 } | Select-Object -First 1
$group = $category.groups[0]

$body = @{
  categoryId = $category.id
  groupId = $group.id
  title = '示例链接'
  url = 'https://example.com/'
  description = '由开发者令牌创建'
  tags = @('示例')
  keywords = @()
  kind = 'tool'
  featured = $false
  status = 'draft'
  sortOrder = 0
} | ConvertTo-Json

$created = Invoke-RestMethod "$nav/api/v1/admin/nav/links" -Method Post `
  -Headers $headers -ContentType 'application/json' -Body $body
```

常用入口：

- `GET /api/v1/nav/catalog`、`GET /api/v1/nav/groups/{groupId}`：公开导航目录，有效 PAT 无需额外 scope。
- `GET /api/v1/admin/nav/links`：读取当前账号可编辑的链接，需要 `nav.link.update`。
- `POST /api/v1/admin/nav/links`：创建链接，需要 `nav.link.submit`；非草稿状态还需 `nav.link.moderate`。
- `PATCH /api/v1/admin/nav/links/{id}`：编辑链接，需要 `nav.link.update`；移动分组或改变审核状态会继续检查额外权限。
- `DELETE /api/v1/admin/nav/links/{id}`、`POST /api/v1/admin/nav/links/bulk`：治理链接，需要 `nav.link.moderate`。
- 结构、健康检查和站点设置端点沿用 OpenAPI，并分别受同名权限约束。

成功返回直接 DTO 或空的 204；失败返回 `application/problem+json`。权限不足为 403，无效或已撤销令牌为 401，资源状态和版本冲突沿用现有 404/409/428 合同。携带失败 PAT 的请求不会回退浏览器登录状态。

## 本地接入

从 Workspace 运行：

```powershell
./environments/nav-local/run.ps1 -Mode Isolated
```

默认入口：Nav `http://nav.dev.yuelili.test:3006`、Account `http://account-nav.dev.yuelili.test:3100/developer-tokens`、Nav API `http://127.0.0.1:8090`。本地环境同时登记 Identity 权限目录与 Nav 在线 PAT 验证。

独立部署需同时配置：

- Identity `pat.applications`：登记 `nav-yueli-web`、权限目录 URL 和 audience。
- Nav `nav.personalToken.*`：site ID、Identity `/api/v1/pat/verify` 和可信 HTTP 策略。

本地实现不代表生产已部署；使用生产 Origin 前必须先核对正式部署记录。
