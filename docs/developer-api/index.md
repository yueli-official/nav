# 开发者 API 生成入口

- [能力目录](developer-api.json)：可申请权限、允许路由与 scope 组合，供脚本和助手读取。
- [开发者 OpenAPI](developer-openapi.json)：只包含 PAT 白名单操作及其引用 Schema。
- [接口说明](developer-api.md)：自动生成的人类可读接口目录。
- [使用指南](../developer-tokens.md)：接入流程、部署状态和使用限制。

声明来自 Nav 实际使用的 PAT 权限目录和路由 allowlist，字段来自产品 OpenAPI。公开目录路由会在生成合同中明确标为 public-or-pat；Nav 没有 PAT 媒体 Profile。

~~~powershell
./scripts/generate-developer-api.ps1 -GeneratorPath ../foundation/tools/developer-api/generate.py
./scripts/generate-developer-api.ps1 -GeneratorPath ../foundation/tools/developer-api/generate.py -Check
~~~

-Check 会重新导出运行时声明和 OpenAPI 并比对全部制品。生成成功只表示当前源码与文档一致，不表示生产已部署或真实令牌已验收。
