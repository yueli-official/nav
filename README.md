# 导航站产品

- 生命周期：活跃的可复用产品类型
- 权威来源：Catalog 产品类型 `nav`、`api/` 迁移/OpenAPI、`web/` 界面
- 消费者：`nav-yueli` 等导航站点实例
- 验证：`pnpm platformctl verify product --file catalog/overlays/local.yaml --root . nav`

Nav 负责精选链接组、主题、推荐链接、搜索和点击指标。`api/` 是持久领域服务，`web/` 提供公开发现与管理。`research/` 下的调研仅提供设计输入，不是运行时真值。
