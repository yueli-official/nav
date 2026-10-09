# 2026-10-09 正式发布准备

用户明确纳入 旧版迁移 范围，Distribution 排除。本文件记录准备证据，不表示已经上线。

域名候选：nav.yuelili.com（Nav/YTV 保留原入口；其余新域名等待用户确认）。服务器正式产品目录拟为 `/projects/yuelili.com/nav`。当前 checkout `697947f89a0e64b40259e25229c1f04b92369cf4`；应用制品来源以下方独立 SHA 为准，后续文档或迁移工具提交不会自动改变既有应用制品。

后端源 `b1fe6ebc4494a91d8a2fa0068fb9cd6a0aaf88a8`：GOWORK=off 的 go test ./...、Linux/amd64/CGO=0 固定构建通过；命令 bindingcheck, bootstrap, healthcheck, nav。迁移文件 28 项 SHA 记录在 manifest.json，数据库集成覆盖以测试配置/日志为准，不能把未设置数据库的单测当生产验证。

前端：固定源 d3ea85c8ba05d71cbbe389e4096cb06c62ca3474；冻结安装、前端测试、类型、生产构建、已安装内容 verified 均通过。Web tar SHA-256 2ba129f791e2b984b1633b3dd6a2ec4d330ecfe8dab15ab70ac7c76f243d4672。UI 为 .6 固定 SHA 06407c95454b215c61b5cf3ef432476de6677d3a78542b07fe706a7c52106c0a，sourceOverrides=[]。

迁移入口：api/manifest/sql/migrations + cmd/bootstrap。只在新建/明确备份的目标数据库预检并应用，不运行 devseed，不复制本地账号/业务数据，不修改旧 checksum。现有生产 Provider 按精确能力绑定，禁止直接套用 standalone compose 重建中央服务。

所需注册：Identity Web client、成员查询与权限目录绑定。

管理员/所有权：现有启动入口仍要求 Bootstrap subject，必须补齐共享首位管理员认领。

网络与密钥：Web/API 仅发布宿主机回环端口。所有 browser issuer/callback 使用正式 HTTPS Origin；Seal Secret、服务客户端 secret  由私有部署配置独立生成并备份，不使用开发默认值。证书需覆盖所选域名，现有证书不能作为新域名证明。

旧版文件完整备份 `/projects/yuelili.com/backups/ui-components-20261009-3/nav-legacy`，逐文件 SHA 核对。tools/prepare-legacy-import.py 生成只读报告与 SQL：8 栏目、58 分组、357 链接，10 组重复 URL 保留。独立 PostgreSQL 验证两次导入保持数量一致、破损关联为 0；验证库已经删除，生产数据没有修改。SQL 为新增 legacy-* 记录，不删除默认 seed；最终切换前仍需核对 preset 重叠及排序展示。

旧数据准备工具 4 项测试通过，独立 PostgreSQL 两次导入数量均为 8/58/357，破损关联 0，临时验证库已删除。

完成门禁：固定源码/依赖与制品可追溯 → 域名/TLS → Provider 注册及绑定检查 → 数据库迁移预检/初始化 → 候选健康/API → 精确切换 → 正式域名 CLI Playwright 手机/桌面、明暗主题、实际交互与登录门禁。初始化认领由用户显式完成；线上管理写入不以匿名测试代替。

现有部署资料：deploy\compose.nav.yaml, deploy\contracts\requirements.json, deploy\deployment.lock.json, deploy\includes\identity\compose.yaml, deploy\topologies\nav.standalone.yaml。

本机证据 `C:/Users/yl/.codex/task-artifacts/sites-ui-upgrade/20261009-3/nav` 与 `nav-api`；原始依赖对照在相邻 20261009-2/first-release-original-dependencies.json。源码快照/依赖为可再生过程产物，收尾删除；固定归档、二进制、Web 包、哈希、日志和迁移准备保留用于后续发布。
