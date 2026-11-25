# 月离导航首次发布

2026-11-24 21:12 CST，`navigation-20261124-1` 切换到 [月离导航](https://nav.yuelili.com/)。首页使用 Vue SPA，提供 10 张站点卡、29 个快捷入口、统一搜索和共享九宫格。后台启用浏览统计和站点设置，API 包含公开品牌资料和开发者令牌权限目录。

## 运行与验收

- Compose 项目为 `nav-next`，入口为服务器 `/projects/yuelili.com/.deploy/nav-navigation-20261124-1/compose.production.yaml`。Web 监听宿主机回环端口 `13104`，API 监听 `18104`，PostgreSQL 在项目私网运行。迁移 `15` 已完成，状态为 clean。
- API 和 Web 健康，镜像 ID 与候选一致。公网 HTTPS 浏览器通过 390 和 1440 像素首页、搜索、九宫格及匿名后台跳转。36 个入口静态资源摘要匹配。候选的五次浏览均实际入库，验收记录在开通公网前清理。
- Nav 专用 OIDC 客户端和开发者令牌应用均为 `nav-yueli-web`。已核对中央登录跳转、生产数据库中的一个管理员授权、无效令牌返回 `401`、权限目录仅在内部网络可访问。真实管理员登录回调、后台写入和有效生产 PAT 调用仍待实测。
- DNSPod A 记录 ID 为 `2458635338`。HTTPS 在 OpenResty 的 `nav.yuelili.com.conf` 接入。证书自动续期由 `nav-yuelili-tls-renew.timer` 管理，续期演练通过。

## 制品和恢复

本次用工作区快照发布。源码快照 SHA256 为 `d21f14acfe466937c12176ad25afad13a295c4c09c375ae24a6f15bac7441aeb`，批准的前端依赖基线为 `5085bc0bcd4078c6095869b015ba3d98178a89a9c217cdc06c87ed1cdeca461d`。后端使用 Foundation Go v0.5.0 及此前软件站发布归档中已核对的四个品牌文件。双端制品、镜像 ID 和验证回执见同目录 `candidate-release.json`、`production-runtime.json`。

旧版 `nav` 容器与目录 `/projects/yuelili.com/nav` 保留，HTTP `9004` 仍可访问。旧数据与切换前新版 PostgreSQL 备份在 `/projects/yuelili.com/backups/nav-navigation-20261124-1/`，归档和备份均已校验可读。旧目录未导入新版生产数据库。

关闭新版公网入口可恢复本次归档的 `nav-vhost.acme-only.conf`，再检查并重载 OpenResty。旧服务可直接通过其原端口核对。Identity 的 Nav 专用配置变更前备份为 `identity-compose.before.yaml`，恢复前须核对后续共享配置变化，仅撤销 Nav 条目。数据库恢复须单独判断。

## fun 入口更新

同日发布 Web `fun-domain-20261124-1`，视频分享站卡片、站点图标和快捷入口指向 [fun](https://fun.yuelili.com/)。正式浏览器通过手机和桌面卡片及视频入口验证。API 复用此次首发镜像，36 个入口静态资源摘要匹配。后续 Web 发布与恢复回执在 `/projects/yuelili.com/.deploy/nav-fun-20261124-1/`。
