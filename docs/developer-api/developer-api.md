# Nav 开发者 API

> 自动生成，请勿手改；从运行时 PAT 声明和产品 OpenAPI 派生。

[使用指南](../developer-tokens.md) · [机器能力目录](developer-api.json) · [PAT OpenAPI](developer-openapi.json)

本文件描述当前源码允许的范围。生产部署：**未知**；真实调用验证：**未知**。实际生效情况查产品部署记录。

请求携带 `Authorization: Bearer <PAT>` 与 `X-Yueli-Token-Site: <实际站点 ID>`。
同一分支内 scope 全部需要，分支之间任选其一；仍须通过当前账号与资源权限检查。目录不是授权凭证。

## 可申请权限

| Scope | 名称 | 当前业务能力 | 说明 |
| --- | --- | --- | --- |
| nav.health_check.run | 管理链接健康检查 | nav.health_check.run | 读取健康状态、运行检查并维护豁免；仍需当前管理员权限。 |
| nav.link.moderate | 治理链接 | nav.link.moderate | 发布、归档、删除或批量治理链接；仍需当前账号拥有对应权限。 |
| nav.link.submit | 提交链接 | nav.link.submit | 创建链接草稿；发布或推荐仍需审核权限。 |
| nav.link.update | 读取和编辑链接 | nav.link.update | 读取可管理链接并编辑当前账号授权范围内的链接。 |
| nav.settings.manage | 管理站点设置 | nav.settings.manage | 读取和修改 Nav 公开站点设置；仍需当前管理员权限。 |
| nav.structure.manage | 管理导航结构 | nav.structure.manage | 读取和管理分类、分组与标签；仍需当前管理员权限。 |

## 媒体 Profile

| Namespace / Profile | 所需 Scope |
| --- | --- |

## 允许调用的操作

| 方法 / 路由 | 操作 | 所需 Scope |
| --- | --- | --- |
| `POST /api/v1/admin/nav/categories` | Create navigation category | `nav.structure.manage` |
| `DELETE /api/v1/admin/nav/categories/{id}` | Delete navigation category | `nav.structure.manage` |
| `PATCH /api/v1/admin/nav/categories/{id}` | Update navigation category | `nav.structure.manage` |
| `GET /api/v1/admin/nav/checks` | List navigation link health checks | `nav.health_check.run` |
| `GET /api/v1/admin/nav/checks/jobs/{jobId}` | Get a navigation link health check job | `nav.health_check.run` |
| `POST /api/v1/admin/nav/checks/run` | Start a navigation link health check job | `nav.health_check.run` |
| `PUT /api/v1/admin/nav/checks/{id}/exemption` | Set a navigation link health check exemption | `nav.health_check.run` |
| `POST /api/v1/admin/nav/groups` | Create navigation group | `nav.structure.manage` |
| `DELETE /api/v1/admin/nav/groups/{id}` | Delete navigation group | `nav.structure.manage` |
| `PATCH /api/v1/admin/nav/groups/{id}` | Update navigation group | `nav.structure.manage` |
| `GET /api/v1/admin/nav/links` | List navigation links | `nav.link.update` |
| `POST /api/v1/admin/nav/links` | Create a navigation link | `nav.link.submit` |
| `POST /api/v1/admin/nav/links/bulk` | Bulk update navigation links | `nav.link.moderate` |
| `DELETE /api/v1/admin/nav/links/{id}` | Delete a navigation link | `nav.link.moderate` |
| `PATCH /api/v1/admin/nav/links/{id}` | Update a navigation link | `nav.link.update` |
| `GET /api/v1/admin/nav/settings` | Get navigation site settings | `nav.settings.manage` |
| `PUT /api/v1/admin/nav/settings` | Update navigation site settings | `nav.settings.manage` |
| `GET /api/v1/admin/nav/structure` | List categories and groups | `nav.structure.manage` |
| `GET /api/v1/admin/nav/tags` | List navigation tags | `nav.structure.manage` |
| `POST /api/v1/admin/nav/tags/delete` | Delete a navigation tag | `nav.structure.manage` |
| `POST /api/v1/admin/nav/tags/rename` | Rename or merge a navigation tag | `nav.structure.manage` |
| `GET /api/v1/nav/catalog` | Get the published navigation catalog | 无额外 scope |
| `GET /api/v1/nav/groups/{groupId}` | Get a paginated published navigation topic | 无额外 scope |
| `POST /api/v1/nav/links/{id}/click` | Record a navigation link visit | 无额外 scope |
| `GET /api/v1/nav/links/{id}/favicon` | Get a proxied navigation link favicon | 无额外 scope |

字段、请求体、响应与引用 Schema 见生成的 OpenAPI。未列入的接口不代表允许 PAT 调用。
