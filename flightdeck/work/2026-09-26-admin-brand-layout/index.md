# nav 后台品牌布局

Status: Open

## 目标

将已确认的共享后台方案推广至本产品：品牌标题卡、紧凑集合、右上角账户和响应式布局；装饰图由本产品提供。保持原品牌色、业务权限和编辑流程。

## 当前

已接入固定 UI 候选 `yueli-ui-0.6.0-preview.20260926.admin-brand.1.tgz`，启用后台品牌主题并将账户移至顶栏。实现与验证进行中；尚未部署此轮改动。

## 验收

类型检查、构建及真实本地 CLI Playwright；有对应线上实例时完成备份、候选冒烟、切换及正式域名验收。

## 用户确认的发布边界（2026-09-26）

用户要求保留旧线上 Nav，等用户本地验收新版后再推进迁移发布。本轮新版已通过类型检查、构建及四个代表页面在 1440/1024/768/390 宽度的 Playwright 验收，账户菜单可用；验收独立组合已由 Workspace CLI 停止，现有站点预览未动。后续恢复本 Work，先供用户本地验收，不直接覆盖旧线上数据。

复验入口为 /manage；本轮隔离参数为 LOCAL_IDENTITY_PORT=8281、LOCAL_ACCOUNT_PORT=3200、LOCAL_IDENTITY_DB=brand_nav_identity，然后运行 Workspace environments/nav-local/run.ps1 -Action Up -Mode Isolated。重启前用 Workspace CLI 核对端口与现有 Session。

2026-09-26 用户授权整理提交；本轮保存为本地 Git 提交，仍保留用户本地验收后再推送、上线的边界。
