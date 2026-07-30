# 导航站首页信息架构与功能对标

> 调研日期：2026-07-14  
> 范围：中文网址导航、个人起始页、书签工作台  
> 方法：仅采用产品官网、官方帮助中心、官方产品页面等一手来源。文中的优先级判断是基于入口位置、官方功能排序与跨产品共性做出的产品推断。

## 结论先行

1. **顶部应保留品牌、一个真正可用的搜索入口和账户/主题等全局操作。** 360 导航把搜索作为核心任务，Raindrop.io、start.me、Toby 和 Linkwarden 则把搜索用于找回已收集内容；它们的共同点不是“顶部链接多”，而是核心动作明确。当前 Nav 顶部的“搜索 / 分类目录”只是页内锚点，缺少状态变化和任务反馈，建议合并为一个可聚焦或打开搜索面板的“搜索”按钮。参考：[360 导航首页](https://hao.360.com/)、[Raindrop.io 搜索](https://help.raindrop.io/using-search)、[start.me 搜索](https://support.start.me/en/articles/9182859-search-for-bookmarks)、[Toby 产品页](https://www.gettoby.com/)、[Linkwarden 产品页](https://linkwarden.app/)。
2. **分类目录仍然有信息架构价值，但不值得占用顶部黄金位置。** 中文导航把分类放在搜索和高频入口之后；书签产品则用集合、侧栏、标签和过滤器承载结构。Nav 当前有 7 个分类、18 个主题、72 个链接，已经超过纯搜索可以轻松浏览的规模，因此应保留正文分类栏和移动端分类选择器，只移除冗余的顶部“分类目录”锚点。参考：[hao123 首页](https://www.hao123.com/)、[hao123 网站地图](https://www.hao123.com/sitemap)、[360 导航首页](https://hao.360.com/)、[Raindrop.io Collections](https://help.raindrop.io/collections)、[Symbaloo Webmix](https://en.help.symbaloo.com/portal/en/kb/articles/what-is-a-webmix)。
3. **热门榜有用，但名称必须与信号匹配。** Pinboard 的 Popular 会显示每个链接的收藏人数，360/hao123 的“热搜”来自搜索趋势；两者都让用户理解“为什么热门”。Nav 目前只有累计 `clickCount` 与最后点击时间，没有按天保存点击事件；实测仅 Behance、Dribbble 各 1 次，其余为 0，直接命名“今日热门”或“近 7 天热门”会失真。第一阶段应显示“编辑精选”或“推荐上手”，并用 `featured` 作冷启动回退；建立时间窗事件数据后，再升级为可信的“近 7 天热门”。参考：[Pinboard Popular](https://pinboard.in/popular/)、[360 导航首页](https://hao.360.com/)、[hao123 首页](https://www.hao123.com/)。
4. **个性化先做低成本、高频动作。** 成熟产品普遍先解决保存、分组、找回和固定入口，再扩展复杂布局：start.me 有书签 Inbox 与页面/Widget/组，Raindrop.io 有收藏、集合与 Favorites，Linkwarden 支持固定链接和集合，Toby 支持 Starred Collections。Nav 适合先做“收藏 / 固定”和“最近访问”，暂不做自由拖拽主页或大量小组件。参考：[start.me 功能页](https://start.me/en/features)、[start.me Inbox](https://support.start.me/en/articles/11424918-bookmarks-inbox)、[Raindrop.io Filters](https://help.raindrop.io/filters)、[Linkwarden 产品页](https://linkwarden.app/)、[Toby Web App](https://web.gettoby.com/)。
5. **提交收录和失效反馈属于治理能力，不应成为首页主任务。** 中文导航通常把“反馈”放在辅助区；书签产品更倾向自动检查失效链接，再让用户修正或删除。Nav 已有后台健康状态字段和检查页面，应优先复用自动检查，卡片菜单提供“链接失效”补充反馈，页尾再放“提交工具”。参考：[360 导航反馈入口](https://hao.360.com/)、[Raindrop.io Broken links](https://help.raindrop.io/broken-links)、[start.me 显示 Broken Links 设置](https://support.start.me/en/articles/9182899-change-user-preferences-display-settings)。

## 产品对比

| 产品 | 首页 / 顶部导航 | 搜索模型 | 分类与发现 | 个性化、提交与治理 | 对 Nav 的启示 |
| --- | --- | --- | --- | --- | --- |
| [hao123](https://www.hao123.com/) | 顶部以天气、日期、登录等轻量工具为辅；页面核心是百度搜索、搜索热点和常用网址。 | 外部网页搜索，可切地图、图片、视频、新闻等搜索域。 | 正文用高密度频道与站点目录承载浏览；另有独立[网站地图](https://www.hao123.com/sitemap)。 | 登录后可承接个人能力；首页发现重点是热点和常用站点。 | 搜索与高频入口优先；目录是正文能力，不是顶部按钮。 |
| [360 导航](https://hao.360.com/) | 顶部辅助入口包括反馈、调研、换肤等；核心区是带明确按钮的综合搜索。 | 外部搜索，可切网页、资讯、视频、图片、地图、AI 等垂直域。 | 搜索下方是高频站点、热搜和按场景组织的目录。 | 有“我的网址 / 添加我的网址”；反馈放辅助区。 | 顶部搜索必须执行真实动作；热门和个性化可提高回访，但不应挤掉搜索。 |
| [2345 导航](https://www.2345.com/) | 顶部提供天气、邮箱、产品、反馈、换肤、登录和收藏夹；搜索仍是首屏核心。 | 百度外部搜索，并可开启历史搜索。 | 搜索下方先放实时热点，再放常用站点和酷站分类。 | 账户收藏夹是个人入口；反馈属于辅助操作。 | 中文大众导航普遍采用“搜索 → 热点/常用 → 分类”的序列。 |
| [start.me](https://start.me/en/features) | 页面由列和 Widget 组成，书签 Widget 是核心；顶部搜索可被 Pro 用户隐藏，说明它是工具而不是品牌装饰。 | Header 搜索自己的书签、笔记和任务；另有 Search Widget 可调用 Google、Bing 或自定义搜索引擎。 | 用页面、Widget、组组织内容，也可浏览社区公开页面。 | 支持一键保存、Inbox、主题/布局、共享、点击分析、重复与失效链接检查。 | 站内搜索和外部搜索是两个明确模式；首页应优先可执行工具，而不是口号。来源：[Search widget](https://support.start.me/en/articles/9182850-search-widget)、[Bookmarks widget](https://support.start.me/en/articles/9324154-bookmarks-widget)、[显示设置](https://support.start.me/en/articles/9182899-change-user-preferences-display-settings)。 |
| [Raindrop.io](https://help.raindrop.io/quickstart) | 应用以顶部搜索、侧栏集合和内容列表构成工作台。 | 一个输入框搜索书签、集合、标签、高亮、批注和正文，并支持当前集合范围与高级过滤；不承担全网搜索。 | Collections 可嵌套、排序，并为不同内容切换列表、网格、Headlines、Moodboard。 | Favorites 是全局过滤器；集合默认私有，可发布为[公共页面](https://help.raindrop.io/public-page)；Pro 自动检测失效链接。 | 对公共 Nav，搜索应先覆盖名称、域名、描述、标签与分类；分类和过滤互补，不是二选一。来源：[搜索](https://help.raindrop.io/using-search)、[集合](https://help.raindrop.io/collections)、[过滤器](https://help.raindrop.io/filters)、[失效链接](https://help.raindrop.io/broken-links)。 |
| [Toby](https://www.gettoby.com/) | 新标签页是视觉工作台，Spaces 和 Collections 是主要结构。 | 应用提供组织级搜索；扩展菜单也允许快速搜索集合。 | Collection 按项目/研究/日常资源分组，支持拖放；可保存整个浏览会话。 | 支持 Starred Collections、公开/团队共享、导入导出。 | 对工作台型导航，“保存当前上下文、固定常用项、快速找回”比资讯流更重要。来源：[Collections](https://help.gettoby.com/support/solutions/articles/66000526357-collections)、[Extension Menu](https://help.gettoby.com/support/solutions/articles/66000497270-how-to-use-the-extension-menu)、[Web App](https://web.gettoby.com/)。 |
| [Symbaloo](https://en.help.symbaloo.com/portal/en/kb/articles/what-is-symbaloo) | Webmix 是可自定义的磁贴网格，不同布局把搜索置于顶部、顶部居中或更大的 Search Layout。 | 每个 Webmix 的搜索框可连接用户选择的搜索引擎；Tile Search 用于发现可加入的资源。 | 多个 Webmix 按主题组织，磁贴可分组、标色和共享。 | 支持浏览器扩展保存、导入书签、公开/私有共享与协作。 | 搜索入口可以很突出，但只有与真实搜索目标绑定才有价值；复杂布局是成熟后的个性化层。来源：[六种布局](https://en.help.symbaloo.com/portal/en/kb/articles/exploring-the-six-symbaloo-layout-views)、[功能页](https://blog.symbaloo.com/pages/symbaloo-features)。 |
| [Pinboard](https://pinboard.in/popular/) | 公共顶部仅有 Recent、Popular、Tour、Howto、Login，极其克制。 | 登录后的核心是个人书签和标签检索。 | Popular 直接按公共收藏信号列出链接，并显示收藏人数；页面末尾展示热门标签。 | 付费存档用于对抗 link rot；产品定位强调速度而非社交。 | 热门榜应展示排序依据；一个小而可信的榜单胜过没有信号解释的“排行榜”。来源：[Popular](https://pinboard.in/popular/)、[官方博客对产品定位的说明](https://blog.pinboard.in/2017/06/pinboard_acquires_delicious/)。 |
| [Linkwarden](https://linkwarden.app/) | 营销页顶部为 Features、Pricing、Blog、Docs、Login、Get Started；登录后的产品以 Dashboard、Collections 与搜索为主。 | 在全部集合中使用搜索与高级运算符查找已保存内容。 | 可固定重要链接与集合；最近内容承担回访发现。 | 支持公开分享、团队协作、浏览器扩展、导入导出，并自动以多种格式保存网页副本。 | 工作台的“热门”不一定是公共趋势；固定项、最近内容和长期可用性同样是回访机制。来源：[产品页](https://linkwarden.app/)、[官方文档](https://docs.linkwarden.app/)。 |

## 功能优先级的跨产品判断

### P0：搜索、访问与可浏览结构

- **一个真实可执行的搜索入口。** 当前 Nav 已能即时过滤名称、描述、域名、标签和分类，顶部按钮应聚焦首页搜索框，或打开全局搜索面板；`/` 与 `Ctrl/Cmd + K` 继续作为快捷方式。[当前搜索组件](../web/app/components/NavigationSearch.vue)
- **搜索范围要写清楚。** 默认只搜本站已收录工具，避免用户误以为是全网搜索。若增加外部搜索，应作为结果面板底部的显式动作，例如“在百度搜索「关键词」”，不要让同一个 Enter 行为在站内/站外之间含糊切换。start.me 将个人内容搜索和外部 Search Widget 分开，是更清晰的范式：[站内书签搜索](https://support.start.me/en/articles/9182859-search-for-bookmarks)、[外部 Search Widget](https://support.start.me/en/articles/9182850-search-widget)。
- **分类目录保留在正文。** 7 分类 / 18 主题 / 72 链接需要稳定的信息气味和扫读路径；桌面侧栏、移动端选择器与搜索结果中的分类筛选应继续存在。[当前首页](../web/app/pages/index.vue)
- **卡片访问必须快且可预测。** 名称、用途、域名/类别和打开方式应一眼可见；这比先增加新闻、天气、小组件更接近导航站核心任务。

### P1：可信的热门 / 新增发现、收藏与最近访问

- **首屏标题组改为任务区，而不是再加一层品牌口号。** 推荐结构：左侧紧凑搜索，右侧或下方是 4–6 条“编辑精选”；每条显示排名、名称、类别与一句用途。数据不足时必须标为“精选”，不能标为“今日热门”。
- **建立点击事件时间窗后再上真实榜单。** 现有模型只有累计 `clickCount` 和 `lastClickedAt`（见[当前数据模型](../api/internal/model/navigation.go)），无法计算“最近 7 天访问最多”或趋势变化。需要按链接和日期保存事件/聚合计数，过滤机器人和管理员测试流量，榜单至少显示“近 7 天 · N 次访问”。Pinboard 的榜单直接显示收藏人数，提供了可信度最低要求：[Pinboard Popular](https://pinboard.in/popular/)。
- **冷启动顺序：** `近 7 天有效访问` → 样本不足时混入 `featured` → 完全无数据时显示 `编辑精选`；不要把人工顺序包装成算法趋势。
- **增加“最近新增”。** 它只依赖已有发布时间，解释成本低，也能给长期用户稳定的新鲜感。
- **个性化从“收藏 / 固定”和“最近访问”开始。** 登录用户跨设备同步；游客可先本地保存。参考：[Raindrop.io Favorites](https://help.raindrop.io/filters)、[Linkwarden Pin Links and Collections](https://linkwarden.app/)、[Toby Starred Collections](https://web.gettoby.com/)。

### P2：提交收录、失效反馈与更深个性化

- **提交收录：** 放页尾或搜索零结果状态，不占顶部主导航。表单至少包含 URL、名称、用途、建议分类和提交理由；进入审核队列，不直接公开。
- **失效反馈：** 卡片的更多菜单提供“链接失效 / 内容不符”，但后台自动健康检查应是主路径。自动检查需容忍防爬、内网和临时故障，Raindrop.io 将失效标记定义为“提示而非裁决”，并提供不同严格度，可作为治理参考：[Raindrop.io Broken links](https://help.raindrop.io/broken-links)。
- **更深个性化：** 拖拽排序、自由布局、主题背景、共享工作台等，应等收藏与回访数据证明需求后再做。start.me 和 Symbaloo 已证明这些能力适合成熟的个人起始页，但不是当前 Nav 的首要缺口：[start.me 功能页](https://start.me/en/features)、[Symbaloo 功能页](https://blog.symbaloo.com/pages/symbaloo-features)。

## 建议的首页信息架构

1. **全局顶部栏**
   - 品牌标识；
   - 一个“搜索”按钮，点击聚焦搜索或打开全局搜索面板；
   - 主题切换、账户菜单；
   - 移除“分类目录”锚点；管理入口继续只对管理员出现。
2. **首屏任务区**
   - 搜索框及清晰的本站搜索说明；
   - 数据冷启动阶段展示“编辑精选 / 推荐上手”，4–6 项；
   - 有足够时间窗数据后替换为“近 7 天热门”，并展示访问次数。
3. **目录主体**
   - 保留桌面分类侧栏、移动端分类选择器；
   - 当前分类下继续按主题分组；
   - 可增加“全部 / 收藏 / 最近访问 / 最近新增”等真实状态筛选。
4. **治理入口**
   - 页尾：“提交工具”“意见反馈”；
   - 卡片菜单：“链接失效”“内容不符”；
   - 后台自动健康检查与人工复核。

## 暂不建议

- 不照搬 hao123 / 360 的新闻资讯流、天气、换肤和广告密度；Nav 的定位是精选工具工作台，不是大众门户。
- 不做没有时间窗、没有样本量说明的“今日热门”或趋势箭头。
- 不把“分类目录”删除为纯搜索页；当前内容规模和探索场景仍需要目录。
- 不在顶部同时放搜索按钮、大搜索框锚点和分类锚点，避免三个入口指向同一首页状态。
- 不优先做 AI 推荐、多人协作、自由拖拽 Widget；这些功能在 start.me、Toby、Symbaloo、Linkwarden 中建立在成熟的收藏与组织模型之上。

## 来源清单

- 中文导航：[hao123 首页](https://www.hao123.com/)、[hao123 网站地图](https://www.hao123.com/sitemap)、[360 导航](https://hao.360.com/)、[2345 导航](https://www.2345.com/)
- 起始页：[start.me Features](https://start.me/en/features)、[start.me Help Center](https://support.start.me/en/)、[Symbaloo Help Center](https://en.help.symbaloo.com/portal/en/kb)
- 书签工作台：[Raindrop.io Help](https://help.raindrop.io/)、[Toby 官网](https://www.gettoby.com/)、[Toby Help Centre](https://help.gettoby.com/support/home)、[Pinboard Popular](https://pinboard.in/popular/)、[Linkwarden 官网](https://linkwarden.app/)、[Linkwarden Docs](https://docs.linkwarden.app/)

