# 导航站产品

- 生命周期：活跃的可复用产品类型
- 权威来源：Catalog 产品类型 `nav`、`api/` 迁移/OpenAPI、`web/` 界面
- 消费者：`nav-yueli` 等导航站点实例
- 验证：`pnpm platformctl verify product --file catalog/overlays/local.yaml --root . nav`

Nav 负责精选链接组、主题、推荐链接、搜索和点击指标。`api/` 是持久领域服务，`web/` 提供公开发现与管理。`research/` 下的调研仅提供设计输入，不是运行时真值。

`api/internal/navauthz` 是 Authorization 公共 Interface 的第二消费者验证：它声明 Site → Category → Group → Link 四级 Scope、
submitter 查询约束、内容维护者、自定义角色和范围委派，并拥有独立生成的 authorization schema migration。当前 Nav 运行时迁移
仍应作为单独产品变更执行，不能把该验证包或 Identity 角色当作现有管理端的隐式授权。
