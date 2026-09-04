# Context

- Nav 拥有导航分组、链接、分类、成员、健康检查、站点设置和授权业务语义；Foundation 只拥有合同生成与通用 HTTP Runtime。
- 普通 JSON API 使用直接 DTO 与 Foundation Problem；重定向、探测或外部协议保留专用 Adapter。
- Domain/Provider cause 与公开 Problem 分离，在 Nav application seam 映射一次；不得公开 `err.Error()` 或上游响应正文。
- 字段 violations、未映射摘要及 trace 技术详情必须保留，反馈优先落在字段或受影响区域。
- 本阶段允许本地提交和 Workspace 真实组合验收，不创建 tag、Release 或镜像。
