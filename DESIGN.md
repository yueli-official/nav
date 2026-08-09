---
name: 月离导航
description: 以中性操作界面呈现身份、本站关系与权限三条可追溯链。
colors:
  primary: "var(--ui-primary)"
  success: "var(--ui-success)"
  warning: "var(--ui-warning)"
  page: "oklch(98.4% 0.006 78)"
  region: "oklch(96.6% 0.008 78)"
  card: "oklch(99.4% 0.003 78)"
  inset: "oklch(94.4% 0.01 78)"
  overlay: "oklch(99.8% 0.002 78)"
typography:
  display:
    fontFamily: "Space Grotesk, DM Sans, system-ui, sans-serif"
    fontSize: "1.25rem"
    fontWeight: 600
    lineHeight: 1.4
  title:
    fontFamily: "DM Sans, system-ui, sans-serif"
    fontSize: "0.875rem"
    fontWeight: 600
    lineHeight: 1.5
  body:
    fontFamily: "DM Sans, system-ui, sans-serif"
    fontSize: "0.875rem"
    fontWeight: 400
    lineHeight: 1.5
  label:
    fontFamily: "DM Sans, system-ui, sans-serif"
    fontSize: "0.75rem"
    fontWeight: 500
    lineHeight: 1.25
rounded:
  control: "0.5rem"
  surface: "0.5rem"
  feature: "0.75rem"
  pill: "9999px"
spacing:
  control-gap: "0.5rem"
  field-gap: "0.75rem"
  region-padding: "1rem"
  page-padding: "1.5rem"
components:
  collection-frame:
    backgroundColor: "{colors.card}"
    rounded: "{rounded.feature}"
  relationship-row:
    backgroundColor: "{colors.card}"
    padding: "0.75rem 1rem"
  member-detail-panel:
    backgroundColor: "{colors.overlay}"
    width: "min(100%, 32rem)"
---

# Design System: 月离导航

## Overview

**Creative North Star: “三线关系账本”**

月离导航的管理界面是一套安静、紧凑、可追溯的操作表面。它以共享的中性管理壳和 CollectionPanel 为视觉底座，用有限的主色表达当前项与可操作状态，让管理员把注意力留给关系、状态与后果，而不是装饰。

成员页是这套世界已经建成的签名表达：管理员按“身份 → 本站关系 → 权限与活动”的顺序阅读一行记录；详情层沿用同一顺序，并把暂停、恢复等高影响操作放到充分说明上下文之后。这个顺序对应真实领域所有权，不是一张账号副本表。

**Key Characteristics:**

- 中性纸面底色、细边界和极轻层次，优先保证长时间操作的清晰度。
- 紧凑行密度，但身份、关系和权限之间始终留有明确分组。
- 主色用于选择、角色与行动线索；成功和警告色只表达有文字说明的状态。
- 详情和确认层负责解释跨系统边界与高影响操作，列表保持可扫描。

## Colors

界面以暖中性表面为主体，天蓝主色只承担导航、选择与权限线索；绿色和琥珀色分别标记正常与暂停，但都必须与文字或图标共同出现。

### Primary

- **导航天蓝：**用于当前导航项、角色标记、可交互焦点和少量身份占位色，不铺满大面积管理表面。

### Secondary

- **正常绿：**用于“正常”等正向成员状态，并配套状态名称与人物检查图标。
- **暂停琥珀：**用于“已暂停”、暂停说明和相关操作，强调影响范围而不把成员描述成系统错误。

### Neutral

- **纸面页底：**承载整个管理壳，略带暖意以降低纯白后台的眩光。
- **区域底与卡片白：**区分侧栏、集合框架和内容行；层级主要由色调与边界形成。
- **内嵌灰与浮层白：**分别用于详情统计块和 Slideover、Modal 等临时操作层。

**The Sparse Accent Rule.** 主色是定位和操作信号，不是装饰底色；同一视口的大部分面积保持中性。

**The Named State Rule.** 成员状态不得只靠颜色表达，必须同时提供中文状态、语义图标或清楚的后果说明。

## Typography

**Display Font:** Space Grotesk（以 DM Sans、system-ui 回退）  
**Body Font:** DM Sans（以 system-ui 回退）

中文由系统字体自然回退；Space Grotesk 为标题和品牌提供克制的几何骨架，DM Sans 让数字、用户键与密集控制保持清楚。界面不引入展示型大标题，信息层级来自字重、字号和分组。

### Hierarchy

- **Display：**仅用于管理页标题和品牌名称，保持短而稳定。
- **Title：**用于详情区块、关系名称与行内主要身份，采用半粗体。
- **Body：**用于名称、状态解释和操作后果，默认使用紧凑正文尺度。
- **Label：**用于列头、移动端语义标签、时间与辅助信息；不牺牲可读性追求更小密度。
- **Identifiers：**公共用户键、策略键等技术标识使用等宽呈现、允许截断或换行，并与可读名称分层。

**The Human Before Key Rule.** 有可读身份时先显示姓名与句柄，再显示稳定用户键；用户键不能取代主标题。

## Layout

管理台由可收起侧栏、固定页头和独立内容面板组成。集合页在宽屏使用受控的最大宽度；搜索与筛选固定在集合顶部，三段关系列头紧随其后，分页固定在集合底部，只有成员行区域滚动。这种拓扑保证管理员在长列表中始终知道筛选条件、字段含义与当前位置。

成员关系账本的宽屏列依次为“身份”“本站关系”“权限与活动”“操作”。每行只暴露进入详情的轻量操作；Slideover 从右侧打开并限制为约 32rem，使列表上下文仍可见。

窄屏从 320px 起保持可用：侧栏收起为菜单，桌面列头隐藏，每行改为纵向信息块，并显式重复“本站关系”“权限与活动”等语义标签；分页与每页数量纵向排列且常驻。详情层占满可用宽度，所需操作不得因断点消失。

**The Stable Collection Rule.** 搜索、控制区、字段语义和分页属于集合框架，数据量变化只能影响中间的行区。

## Elevation & Depth

系统默认平坦，通过相邻暖中性色、1px 语义边界和极轻卡片阴影建立层级。悬浮、选择与焦点使用色调变化和明确轮廓，不以强阴影或位移制造“漂浮卡片”。Slideover 和 Modal 是少数真正覆盖内容的层级。

**The Flat-by-Default Rule.** 常驻集合与行保持平坦；只有临时覆盖层和共享卡片允许可感知的提升。

## Shapes

控制与常规表面使用温和圆角（0.5rem），集合框架等较大边界使用稍大的圆角（0.75rem）。头像和状态计数使用圆形或胶囊形；这种轮廓只服务身份与短状态，不能把整行切成独立药丸。细边界承担主要分组，分隔线保持连续。

## Components

### Collection Panel

- 搜索始终位于首行；宽屏筛选在其下平铺，窄屏以可操作的“筛选”按钮展开，并显示当前总数。
- 列头、加载、空结果、请求失败和分页都属于同一有标签的集合区域；空集合与筛选无结果使用不同文案和恢复动作。
- 行区独立滚动，列头和分页不随成员记录离开视口。

### Relationship Ledger

- 身份列组合头像、显示名、句柄和稳定用户键；头像失败时回退为首字占位。
- 本站关系列只显示成员状态、加入时间和最近进入；不得混入 Identity 账号状态。
- 权限与活动列用角色徽章及投稿、待处理申请摘要；无角色明确写“普通成员”或“无维护角色”。
- 宽屏按列横向对齐；窄屏按同一语义顺序纵向展开，不能仅靠位置暗示字段。

### Member Detail Panel

- 详情固定按 Identity 用户、Nav 成员关系、本站权限与活动、成员操作四段展开。
- Identity 段提供稳定用户键和前往用户中心的外链，不在 Nav 内复制资料编辑能力。
- 暂停成员时展示时间、操作人、原因及“公开内容仍可浏览、已认证操作被阻止”的审计说明。

### Status and Role Badges

- 成员状态使用 subtle 徽章和人物语义图标；角色使用主色 soft 徽章；无角色使用中性 soft 徽章。
- 徽章是摘要，不替代详情中的完整含义、来源或影响说明。

### High-impact Actions

- 暂停入口只出现在成员详情，先说明本站范围，再打开带必填原因的确认 Modal。
- 当前管理员不能暂停自己的成员资格；禁用状态旁必须显示原因。
- 恢复操作保留原有角色、申请与审计历史，不暗示重新创建账号。

### Navigation and Authorization

- “成员”与“权限策略”是并列管理入口：前者管理产品关系，后者管理角色、能力、申请和自动授权规则。
- 权限页的策略变更先进入草稿，验证影响后整体发布；新成员自动授权必须明确它不改变用户中心身份或全局角色。

## Do's and Don'ts

### Do:

- **Do** 在涉及参与者治理的集合中明确“身份 → 产品关系 → 权限/活动”的所有权和阅读顺序。
- **Do** 把高影响操作放入有上下文的详情层，写清作用范围、保留内容和恢复方式。
- **Do** 保留键盘可见焦点、语义 section/table-like collection 标签、可读按钮名称以及不依赖颜色的状态文本。
- **Do** 让其他站点复用领域分离、稳定集合拓扑和审计式操作说明，同时由各产品定义自己的成员状态与活动字段。

### Don't:

- **Don't** 把 Nav 成员称为 Identity 用户副本，也不要把成员存在等同于角色、授权或申请。
- **Don't** 在 Nav 复制账号安全、资料编辑或全局用户治理；这些动作必须导向 Identity 用户中心。
- **Don't** 让整个页面随长列表滚动，或让筛选、列头、分页在主要操作期间消失。
- **Don't** 把关系账本的三列和“正常/暂停”语义逐字复制到无相同领域模型的站点；可复用的是分层方法，不是 Nav 的业务字段。
- **Don't** 用强阴影、大面积品牌色或装饰性卡片削弱管理集合的连续性和扫描效率。
