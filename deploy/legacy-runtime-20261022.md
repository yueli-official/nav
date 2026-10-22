# 旧版Nav运行时安全重建

本记录对应仍在线的旧Vue实现，不是当前本地新版产品部署。2026-10-22从固定生产commit `7e208877c3be00c6e7a332b18e5c504818d245d4`重建Go1.26.10，只应用相邻`legacy-security-20261022.patch`去除后台口令明文日志；游客收藏新增/编辑权限保持原样。权限收紧尚待用户选择，不能把候选鉴权补丁应用到本次重建。

## 固定输入

- 原生产image ID：`sha256:8ca8f7b44ad77d9ef28736eb7252f8ab131b56c8b32ec870d9a6806d0d7f8e16`。
- 原commit的git archive SHA256：`c9efbeeccd06f3c93ef763cd59189ed09b802fa7fd97efa23b329f410f3be53a`。
- 当前主程序SHA256：`d4444908a27cbc5e30236413388baa3788bebd3efccd1064e2e440c2c5ee6141`。
- 当前镜像：`yueli/nav:runtime-go1.26.10-20261022-2`，image ID `sha256:217479fb6a390ce810ddccd0a751b89bc897146293d12a3fb168ed8abc69e09c`。

在正式checkout之外提取该commit的来源，应用补丁，使用Go1.26.10、`GOWORK=off`、`GOTOOLCHAIN=local`检查`go test -mod=readonly ./...`；原包无测试用例，此步骤只构成包检查。设置`GOOS=linux`、`GOARCH=amd64`、`CGO_ENABLED=0`后执行：

```text
go build -mod=readonly -trimpath -ldflags="-s -w" -o server .
```

以原生产image ID为父镜像，只复制新`server`到`/app/server`，将已修复的宿主机CA bundle复制到`/etc/ssl/certs/ca-certificates.crt`。复用父镜像Entrypoint、原六个Vue静态制品、Compose环境及数据挂载。不得从当前前端源码构建后混入该旧版镜像。

## 线上入口和恢复

服务器`/projects/yuelili.com/nav/docker-compose.yml`，精确服务`nav`，保留9004映射和原restart=no策略。旧版读取`data/config.yaml`与`data/apikey.txt`，本次没有SQLite迁移、密钥轮换或权限变更。

切换前用复制数据和随机测试口令验证：公开导航JSON一致，游客POST/PUT缺失口令时的参数校验状态与原版相同，原六个静态文件摘要相同，新stdout/stderr均不包含测试口令。切换后核对运行二进制、环境变量映射、重启策略和静态制品；CLI Playwright桌面/手机通过并直接查看截图。新进程日志不再输出后台口令，历史日志未做全盘清除。

正式Nav443当前没有TLS虚拟主机；浏览器通过SSH隧道读取原HTTP9004，合成HTTPS Origin只用于路由测试，不能作为公网Nav HTTPS通过的证据。

精确回执位于服务器`/projects/yuelili.com/.deploy/runtime-security-20261022-02/nav-runtime/`及`deployment-release.json`；私有备份在`/projects/yuelili.com/backups/runtime-security-20261022-02/`，包含两次切换前的Compose/YAML/API key，并已有离机副本。恢复前逐项比较当前Compose及image ID，仅回退精确服务；不要覆盖后续部署或直接覆盖运行中的数据。
