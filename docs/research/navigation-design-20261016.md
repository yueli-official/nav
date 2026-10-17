# Nav 导航站视觉与 Vue 动效调研

调研日期：2026-10-16。范围是案例与技术选型，不改产品实现、不安装依赖、不部署。以下推荐是针对本站“做精不做大”、首页搜索＋精选＋分类卡片的设计判断，不是市场份额排名。

## 结论

优先参考 TOOOLS.design 的资源导航表达、Raycast 的品牌色与光影层次、Linear 的克制排版；以现有 Nuxt UI/Foundation 为基础，少量采用 Inspira UI 效果，或用 Motion for Vue 实现统一的交互动效。Vue Bits 适合继续浏览创意，但它的附加许可限制需要按具体采用和分发方式核对。

这条路线保留现有基础组件与版本门禁，重点改善搜索聚焦、精选卡片、分类图标与鼠标/键盘反馈。当前没有做候选组件的安装、打包体积、Nuxt SSR 集成或手机帧率测试，因此“适合”是选型判断，不是兼容性通过结论。

## 导航与发现类案例

| 案例／官方来源 | 页面中可观察的结构 | 对本站的判断 |
| --- | --- | --- |
| [Product Hunt](https://www.producthunt.com/) | 首页含当日产品、上周／上月榜单与产品类别。 | 借鉴有限数量的精选榜单、图标＋一句话价值说明和时间区间；本站规模不需要引入社区、投票和发布平台。 |
| [Uneed](https://www.uneed.best/) | 每日产品发布列表、搜索入口、产品标签与类别。 | 借鉴卡片信息顺序和搜索入口；首页应围绕使用者挑选工具。 |
| [TOOOLS.design](https://www.toools.design/) | 分类入口、设计资源、精选工具、按工作类型组织的资源入口。 | 最贴近本站“首页分类卡片→分类内容”结构，可重点参考图标统一、轻背景层次和卡片留白；本站仍保持紧凑，不直接复制其大字号。 |
| [优设导航](https://hao.uisdc.com/) | 分类侧栏、搜索、热门推荐和多组资源卡片。 | 参考中文命名和扫读密度；它的大目录布局适合大量资源，本站已选择更精简的首页与分类页。 |

以上案例均依据其自有页面，不引用第三方流量或“最佳”排名。

## 视觉与酷炫效果案例

| 案例／官方来源 | 值得观察 | 对本站的判断 |
| --- | --- | --- |
| [Linear](https://linear.app/) | 本次首屏截图呈现深色底、低对比边框、集中排版与产品界面分层。 | 借鉴层次、边框和排版秩序；首屏标题与留白需按本站紧凑目标缩小。 |
| [Raycast](https://www.raycast.com/) | 本次首屏截图呈现深色底、红色光带、清晰主标题和收敛的主操作。 | 借鉴品牌光晕、搜索焦点和局部色彩，不必复制整屏宣传页。 |
| [Lusion](https://lusion.co/) | 官网介绍自身的3D视觉叙事、动效和交互体验，并列出相应作品。 | 作为高强度视觉创意库，适合看氛围和局部互动；当前仅查阅官方页面文字，未实际验证其3D交互。 |
| [Awwwards](https://www.awwwards.com/) | 提供 Interaction Design、CSS animations、WebGL、3D 等作品分类。 | 用于继续寻找具体参考作品，不等于本站应采用其所有效果。 |

## 本站可选视觉方向

1. **清透精选**：浅色或深浅双主题，背景局部柔和渐变、统一分类图标、轻边框卡片。以 TOOOLS 的资源表达为主要参考，适合长期日常使用。
2. **深色科技**：保留内容布局，用 Raycast 的品牌光影和 Linear 的克制层次；主要效果集中在搜索和首页精选。
3. **创意实验**：首页增加一处可交互光场或动态背景，目录内容维持稳定。这是探索方向，需先证明手机性能与阅读体验。

建议先做1或2，效果控制在搜索、精选、卡片悬停/焦点与内容切换四处。卡片进入可轻微错峰，单卡轻抬升与边框渐亮；动效时间、幅度与装饰强度需在原型中试调。输入和列表首屏保持SSR可读，系统要求减少动态效果时应提供静态版本。

## 调研证据与限制

- CLI Playwright 实际打开 TOOOLS、优设、Linear、Raycast、Inspira UI、Vue Bits 六个官网，HTTP均200；检查了1440×1000首屏截图。
- 本机截图与 observations.json 位于 D:/coding/yueli-official/task-artifacts/nav/public-navigation-20261016/research/。使用 reducedMotion=reduce，故静态截图不能证明动效流畅、完整播放或移动端性能。
- Inspira 在该首屏截图中仅页头和框架明显可见，不把它当完整组件展示或兼容验收。
- 本报告没有改变当前导航页、组件版本、统计实现或运行服务。当前Nav累计点击与周榜不是同一统计口径；若后续呈现“本周热门”，应先补齐时间区间统计。


## Vue 动效组件调研

调研日期：2026-10-16。范围是官方文档、许可证与部分源码声明；本轮没有安装依赖、修改页面，也没有做 Nuxt 运行兼容性、SSR hydration 或性能实测。以下推荐是选型判断，不代表兼容性验收通过。

本地 Nav `web/package.json` 当前声明 Nuxt 4.5.1、Vue 3.5.39、Tailwind CSS 4.3.1、Nuxt UI 4.9.0，并接入 Foundation 依赖基线。因此应在现有设计系统上加少量视觉增强，不整套替换，也不绕过跨站依赖门禁。

### 候选比较

| 方案 | 定位与采用方式 | Vue / Nuxt / SSR 证据 | 许可证与适用判断 |
| --- | --- | --- | --- |
| **Inspira UI** | 复制所需 Vue 源码；覆盖卡片、背景、文本、交互、布局。适合挑局部效果，并保留产品自己的样式控制。 | 官方安装页面向 Vue / Nuxt，当前指南面向 Tailwind v4，还直接给出 Nuxt UI 颜色变量映射。不同效果的客户端要求不同，不能把“支持 Nuxt”理解成所有组件都能完整 SSR。 | 开源仓库为 MIT，需保留版权/许可声明；Pro 单独收费，不等于开源版授权。是本轮优先考察对象。 |
| **Vue Bits** | React Bits 的官方 Vue 移植；按组件复制或 CLI 安装源码，并非整包 UI 依赖。偏视觉表现、文本与背景实验。 | 官方称面向 Vue / Nuxt；当前介绍声明 Vue 3 + TypeScript + Tailwind。没有在已读材料中找到所有组件统一 SSR 保证，DOM / Canvas / WebGL 效果仍需逐件检查。 | **MIT + Commons Clause，不是纯 MIT**。网站内使用与把组件本身再分发是不同情形；放进 Foundation 再发布前需单独核对条款。适合参考和谨慎挑选。 |
| **Motion for Vue (`motion-v`)** | 动效引擎，提供弹簧、入退场、布局、手势、滚动等，不提供一整套站点视觉。可局部 import 或用 Nuxt module。 | 官方支持 `motion-v/nuxt`；`motion` 组件支持 SSR 初始状态，部分需 DOM 测量的 SVG 属性除外；另有按需特性加载和减少动态效果 API。 | 核心 MIT；Motion+ 示例和高级产品单独区分。适合统一短促、克制的交互语言。 |
| **TresJS** | Vue 声明式 Three.js / 3D 场景；需 Three.js，可选 Cientos、后期处理。不是常规 UI 组件库。 | 有官方 Nuxt module，自动将 `TresCanvas` 作为客户端渲染，并处理模板编译器；3D 场景本身不是 SSR 内容。 | 核心 MIT；模型、贴图和附加资源仍各有授权。当前轻量导航首页不优先采用。 |

官方依据：[Inspira 安装与 Nuxt UI 映射](https://docs.inspira-ui.com/docs/en/getting-started/installation)、[Inspira 开源许可证](https://github.com/unovue/inspira-ui/blob/main/LICENSE)、[Inspira Pro](https://pro.inspira-ui.com/activate-license)；[Vue Bits 介绍](https://vue-bits.dev/get-started/introduction)、[官方仓库](https://github.com/DavidHDev/vue-bits)、[完整许可证](https://github.com/DavidHDev/vue-bits/blob/main/LICENSE.md)；[Motion Vue 接入](https://motion.dev/docs/vue)、[SSR 说明](https://motion.dev/docs/vue-motion-component)、[核心许可证](https://github.com/motiondivision/motion-vue/blob/master/LICENSE)；[TresJS 安装与 Nuxt 模块](https://docs.tresjs.org/getting-started/installation/)、[TresJS 许可证](https://github.com/Tresjs/tres/blob/main/LICENSE)。

Vue Bits 的 Vue 支持另以实际源码交叉确认：[SpotlightCard.vue](https://github.com/DavidHDev/vue-bits/blob/main/src/content/Components/SpotlightCard/SpotlightCard.vue) 使用 Vue 单文件模板、`script setup lang="ts"`、`ref` 与 `useTemplateRef`；不能依据官网某段营销示例的 JSX 外观判断其组件属于 React。这个源码检查同样不代表所有组件已经过本站 SSR 验证。

### 适合 Nav 的具体效果

1. **分类卡片轻微聚光**：可看 [Inspira Card Spotlight](https://docs.inspira-ui.com/docs/en/components/cards/card-spotlight)。仅用于分类卡或本周主推卡，鼠标经过时小面积提亮；触屏保持静态可读，键盘焦点用清晰描边。不要让所有站点卡持续发光。
2. **搜索框的焦点质感**：可参考 [Inspira Input](https://docs.inspira-ui.com/docs/en/components/input-and-forms/input) 的径向光效或 [Neon Border](https://docs.inspira-ui.com/docs/en/components/special-effects/neon-border) 的边框思路。实际输入控件继续采用现有 Foundation / Nuxt UI 行为，按现有主题色做轻量装饰，避免复制第二套输入组件和全局样式。
3. **精选卡片入场、筛选结果切换**：用 Motion 的 opacity / transform 与布局过渡建立节奏，减少突变。仅 hover 颜色、阴影和小位移时直接 CSS transition 即可，Motion 官方也推荐简单效果使用 CSS。[Motion 使用与 CSS 边界](https://motion.dev/docs/vue)
4. **首页一处弱背景**：渐变、静态噪点或低密度网格就可形成品牌感；Vue Bits 可作为效果灵感库。其官方建议一页不宜叠太多效果，并注意移动端关闭部分效果；这更符合“精选导航”而非特效展厅。[Vue Bits 性能建议](https://vue-bits.dev/get-started/introduction)

不建议当前采用：持续粒子 / 全屏 WebGL、鼠标拖尾、自定义光标、首屏大段逐字动画、强制横向滚动、每张卡片 3D 翻转。它们增加导航任务的操作和性能成本；这属于本项目适配判断，不是对库质量的否定。TresJS 可留给将来确有 3D 展示需求的专题，而非作为本轮首页基础依赖。

### 与现有技术栈一起用的边界

- **保留一套基础 UI 与主题变量**：Inspira 安装页已经有 Nuxt UI 映射，优先复用 `--ui-*`；不要再引入一套 shadcn / Reka 全局主题覆盖。按具体复制组件检查附加依赖，不把安装指南中的整套配套包等同于每个效果的必要依赖。[Inspira 安装页](https://docs.inspira-ui.com/docs/en/getting-started/installation)
- **动效引擎只选一种实际需要的方案**：Inspira 当前依赖指导包含 `motion-v`；若采用其相关组件，应复用同一版本，避免混入旧教程的 `@vueuse/motion` 并行维护。Motion 的 Vue API 与 Nuxt module 有官方入口，仍需通过 Nav 的现有基线、lockfile、typecheck 与最终构建验证。[Inspira 安装](https://docs.inspira-ui.com/docs/en/getting-started/installation)、[Motion Nuxt](https://motion.dev/docs/vue)
- **SSR 内容不依赖动效才能看见**：标题、搜索、分类和站点链接应正常输出。依赖浏览器 API 的装饰在客户端挂载，并预留尺寸；先检查随机数、DOM 测量、观察器、清理行为。Vue 官方明确浏览器 API 及生命周期在 SSR 的差异。[Vue SSR 编写原则](https://vuejs.org/guide/scaling-up/ssr)
- **减少动态效果与按需加载**：接入后需尊重系统减少动态效果设置，Motion 提供 `useReducedMotion`；确实需要引擎时可评估 `LazyMotion`，实际包体以本站构建结果为准，不照搬文档数字。[减少动态效果](https://motion.dev/docs/vue-use-reduced-motion)、[按需特性加载](https://motion.dev/docs/vue-lazymotion)
- **Vue Bits 许可证要保留区别**：官方条款允许作为应用/网站/产品的一部分使用，但限制组件自身以独立、打包、模板或移植形式销售、再许可或再分发。因此不把它作为纯 MIT 源码自动推广进共享 UI 包；真正采用具体组件前确认用途与源文件许可。这里是对已读取条款的工程使用边界记录，不扩大解释为“所有商业网站禁用”。[Vue Bits LICENSE](https://github.com/DavidHDev/vue-bits/blob/main/LICENSE.md)

### Aceternity / Magic UI 与 Vue 的区别

[Aceternity UI](https://ui.aceternity.com/) 官方是 React + Tailwind + Motion；[Magic UI](https://magicui.design/docs) 官方也是 React 组件与营销页面素材路线。它们很适合找视觉方向，但不能把 React 代码直接装进 Nuxt。Inspira 是 Vue / Nuxt 方向的独立项目，吸收这两者的灵感及社区贡献，并非这两个产品官方提供的 Vue 包；应按 Inspira 自身源码、依赖和许可证评估。[Inspira 作者介绍](https://igrislabs.com/products/inspira-ui/)

### 本项目选型建议

先做“排版、留白、图标与颜色统一 + 搜索框一个亮点 + 分类卡片轻交互”，以现有 Foundation / Nuxt UI 为底。额外组件优先从 **Inspira 挑 1–2 个局部效果**；需要复杂入退场再评估 **Motion for Vue**。Vue Bits 用来浏览灵感并按许可逐件决定，TresJS 暂不引入。后续原型验收再看手机、键盘、减少动态效果、SSR 和实际 JS 成本，不能以官方示例好看替代本站验证。
