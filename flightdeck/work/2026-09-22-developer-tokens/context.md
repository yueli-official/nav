# 上下文与权限边界

用户要求按 BVideo、Paste、Gallery、Nav、Distribution 的方向让更多站点支持开发者令牌。仓库规则禁止建立笼统的多站推广 Work，因此本 Work 只交付 Nav。

Nav 使用 Identity 作为凭证权威，使用自身 Authorization Runtime 判断当前业务权限。Nav 没有直接媒体上传接口；站点图标由服务端依据链接 URL 抓取并缓存，因此本产品 PAT 不接入 Asset。

可按当前账号权限委派的能力：

- `nav.link.submit`：创建链接草稿；非草稿状态仍需审核能力。
- `nav.link.update`：读取可管理链接并编辑当前授权范围内的链接；内容维护者继续受提交者关系约束。
- `nav.link.moderate`：发布、归档、删除或批量治理链接。
- `nav.structure.manage`：读取和管理分类、分组与标签。
- `nav.health_check.run`：读取健康检查状态、运行检查与维护豁免。
- `nav.settings.manage`：读取和修改公开站点设置。

公开目录、分组、点击记录和 favicon 路由允许有效 PAT 调用，但不需要额外 scope。成员目录、成员状态、角色申请/邀请/授予、授权控制台、`/api/v1/me`、可信权限目录本身和其他未声明路由默认拒绝 PAT。权限目录只接受 `identity-svc` 且带 `personal-token:permissions` scope 的服务身份。

本地 site/application ID 为 `nav-yueli-web`，资源 audience 沿用 `nav-yueli-web`。Workspace Environment 负责连接 Identity 权限目录与 Nav 在线验证；产品不复制令牌存储或 Cookie 行为。
