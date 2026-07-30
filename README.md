# 导航站产品

- 生命周期：活跃的可复用产品类型
- 权威来源：Catalog 产品类型 `nav`、`api/` 迁移/OpenAPI、`web/` 界面
- 消费者：`nav-yueli` 等导航站点实例
- 验证：`pnpm platformctl verify product --file catalog/overlays/local.yaml --root . nav`

Nav 负责精选链接组、主题、推荐链接、搜索和点击指标。`api/` 是持久领域服务，`web/` 提供公开发现与管理。`research/` 下的调研仅提供设计输入，不是运行时真值。

`api/internal/navauthz` 是 Authorization 公共 Interface 的第二个完整运行时消费者：它声明
Site → Category → Group → Link 四级 Scope、submitter 查询约束、内容维护者、自定义角色、范围委派、申请审批和策略发布，
并使用 Nav 实例自己的 authorization schema 与策略数据。Identity 只提供认证后的 Subject，不提供 Nav 角色。

管理 API 以 effective capability 和资源 Scope 作为唯一门禁；Catalog renderer 只负责将受保护管理员写入
`authorization.bootstrapAdministratorSubs`，不再保留 `nav.operatorSubs` bypass。链接跨 Group 移动时会通过 Foundation
`ResourceScopeRelocator` 同步 Scope 父级，避免沿用旧分组授权。

注册自动授权规则已经进入 Nav 本地策略控制面。产品授权入口先 preview、再按需 reconcile：启用规则后，注册用户第一次使用
Nav 时会在本实例幂等获得内容维护者 Grant；关闭规则不会撤销已有 Grant。该流程不要求 Identity 保存 Nav 角色或维护站点订阅。
