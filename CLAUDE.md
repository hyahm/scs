# CLAUDE.md

本文件给 Claude Code 提供项目全局上下文。详细的可执行审查清单见 `.trae/skills/`（报警/并发/进程管理）。

---

## 项目定位

**SCS（Service Control Service）** —— Go 1.25 编写的跨平台进程编排与监控守护进程。
- `scsd`（`cmd/scsd/`）：常驻守护进程，通过 YAML 描述一组脚本/服务并管理其生命周期（启动/停止/重启/kill/更新/副本/定时任务/监控报警）。定位类似 `supervisor`，但内置 TLS、多通道报警、原子性停止（cannotstop）、WebSocket 实时日志。
- `scsctl`（`cmd/scsctl/`）：配套 cobra CLI，可并发管理多个 scsd 节点。
- module：`github.com/hyahm/scs`，当前版本 `v3.8.5`（`global/var.go`）。

## 技术栈

- Web 框架：自研 `github.com/hyahm/xmux`（分组、token 权限 keys、中间件 module、WebSocket）
- CLI：`github.com/spf13/cobra`
- 日志：自研 `github.com/hyahm/golog`（⚠️ `(*Log).Sync()` 会无条件 close channel，**不可重复调用**；`Start()` 内局部捕获 logger 实例 + `defer` 同步，避免 restart 重入时误刷新进程日志）
- 配置：`gopkg.in/yaml.v3`
- 监控：`github.com/shirou/gopsutil`（CPU/内存/磁盘）
- 报警后端：gomail（邮件）/ Rocket.Chat / Telegram / 企业微信 / 钉钉 / HTTP 回调

## 常用命令

```bash
# 运行守护进程（默认读 scs.yaml，可用 -f 指定）
go run ./cmd/scsd/scsd.go -f scs.yaml
go run ./cmd/scsd/scsd.go -v          # 查看版本

# 运行 CLI
go run ./cmd/scsctl/scsctl.go

# 测试（目前仅 internal/cache/word_test.go）
go test ./...

# 构建（见 build.bat / update.sh）：交叉编译 windows/linux/darwin 的 scsd + scsctl
#   Windows 下用 build.bat（内部是 PowerShell 设置 $env:GOOS）
#   Linux 下用 update.sh（git pull + GOPROXY=https://goproxy.cn）
GOOS=linux   go build -o bin/scsd   ./cmd/scsd/scsd.go
GOOS=linux   go build -o bin/scsctl ./cmd/scsctl/scsctl.go
```

配置文件：根目录 `scs.yaml`（已在 `.gitignore`，含本地 token），结构参考 `default.yaml` 与 `README.md`。

## 架构分层

```
cmd/scsd/scsd.go        守护进程入口：ReadConfig → FirstStartAllScript → WriteConfig → HttpServer
cmd/scsctl/scsctl.go    CLI 入口（读 ~/.scsctl.yaml）

api/server.go           HttpServer() 入口（监听 + TLS）
api/handle.go           ★ 三级 token 权限路由（AdminHandle / ScriptHandle / simpleHandle）
api/handle/*.go         各 HTTP handler
api/module/             中间件（CheckToken、UpdateConfig、写回配置）

internal/cache/         ★★★ 进程管理核心，所有新功能放这里
  ├─ server.go            Server 结构体 + Start/Stop/Kill/Restart/Remove/UpdateServer 状态机
  ├─ server_store.go      serverStore 全局单例（map[string]*Server + RWMutex）
  ├─ script_store.go      Script 配置态仓库（map[string]config.Script + RWMutex）
  ├─ start.go             StartAsync / Start（局部捕获 ctx/logger，避免 restart 重入误杀）
  ├─ stop.go / wait.go    停止 / 等待退出
  ├─ cron.go              定时任务循环
  ├─ factory.go           Server/Script 工厂
  ├─ log.go               stdout/stderr 捕获 → Logger
  ├─ prestart.go          启动前准备
  ├─ script_windows.go    Windows: powershell -c + taskkill /F /T /PID
  └─ script_unix.go       Unix:   /bin/sh -c + syscall.Kill(-pgid, SIGKILL)

internal/config.go       ReadConfig / WriteConfig（YAML 读写）
internal/cache.go        全局 msgCache channel（⚠️ 当前无消费者 goroutine，见陷阱）
pkg/config/              配置结构体 + 探测器(probe) + 报警器实现 + dispatcher 去抖
pkg/message/             报警消息体 + SendAlerter 接口
pkg/status.go            状态常量 STOP/RUNNING/WAITSTOP/WAITRESTART
client/command/          scsctl 的 cobra 子命令
controller/              ⚠️ 旧三层架构，几乎全部被注释 —— 死代码，禁止在此新增功能
```

## Server 状态机（核心，务必理解）

每个 `Server` 副本（`internal/cache/server.go`）通过 3 个并发原语协作：

1. **`Ctx + Cancel`**（`context.WithCancel`）：取消日志 goroutine、cron ticker、进程本身（`exec.CommandContext(ctx, ...)`，ctx 取消即杀进程）。`Start()` 内用**局部变量**捕获 `ctx, cancel` 并 `defer cancel()`，避免 restart 重入时旧 goroutine 的 defer 误杀新进程。
2. **`restartOnExit bool`**（`svc.mu` 保护）：Restart 的"退出后重启"意图。`Restart()` 在 RUNNING 时置 true 并 `Cancel()`，`wait()` 进程退出后消费它并 `StartAsync()`；`Stop/Kill/Remove` 入口清零，避免被自动拉起。
3. **`Status.Status`**（`pkg/status.go`）：`STOP / RUNNING / WAITSTOP / WAITRESTART`。

**cannotstop（原子停止）**：`/cannotstop/name` 置 `CanNotStop=true`，`/canstop/name` 置 false 并向 `Status.ChStop` 发信号。`Stop/Kill/Restart/Remove` 入口都 `if GetCanNotStop() { <-ChStop }` 阻塞等待，防止脚本数据处理中途被杀。

操作流转：
- `StartAsync()` → 仅当 `Status==STOP` 时 `go Start()`
- `Start()` → 局部捕获 Logger + ctx/cancel → cron 路径或普通 `start()` → `wait()`
- `Stop()` → 清 restartOnExit → 等 cannotstop 解除 → `Cancel()`
- `Kill()` → 清 restartOnExit → 等 cannotstop 解除 → `Cancel()` + `kill()` 强杀进程组
- `Restart()` → 等 cannotstop 解除 → STOP 直接 `StartAsync()`；RUNNING 置 restartOnExit + `Cancel()`（由 wait 重启）
- `Remove()` → 清 restartOnExit → 等 cannotstop 解除 → `Cancel()` + `kill()` → `RemoveServer()` 删 server 缓存（script 配置删除 + 配置文件回写由 handler 层 + UpdateConfig module 负责，remove handler 当前仍为空壳）

## HTTP 权限模型（`api/handle.go`）

三级嵌套 RouteGroup，越内层权限越低：
```
AdminHandle   (key=admin,   + CheckToken module)   改配置/重载/报警管理
 └─ ScriptHandle (key=script)                      stop/kill/remove/cannotstop
     └─ simpleHandle (key=simple)                  status/start/update（调试用）
```
token 由 `config.Token` 派生出 admin/script/simple 三种角色 key，经 `AddPageKeys` + `module.CheckToken` 校验。

## 协作约定（重要）

1. **中文优先**：注释、错误信息、commit message、文档全部用中文（与现有代码一致）。
2. **跨平台**：任何进程操作都要同时考虑 `script_windows.go` 和 `script_unix.go`，用 build tag 区分，不要写死某一平台。
3. **死代码禁区**：不要在 `controller/` 下新增代码（整层被注释）。新功能一律放 `internal/cache` + `api/handle`。
4. **全局 store 访问**：`serverStore` / `script_store` / `dispatcher` 都用 `sync.RWMutex` 包装，写 map 用 `Lock()`、读用 `RLock()`，**持有锁时不要回调会再次获取同一把锁的方法**。

## 已知陷阱（改动前必看）

- **golog `Sync()` panic**：`(*Log).Sync()` 内部无条件 `close(channel)`，对同一 Logger 二次调用会 `panic: close of closed channel`。`Start()` 内用**局部变量捕获 logger** + `defer func(){ recover(); logger.Sync() }()`，确保 restart 重入时只同步本次实例。
- **Restart 重入的字段覆盖**：`Ctx/Cancel/Logger/Cmd` 是共享字段，restart 路径下两次 `Start()` goroutine 重叠会互相覆盖。已通过"`Start()` 局部捕获 ctx/cancel/logger"根治 defer 误杀新进程；但 `appendRead`/`successAlert`/`cron` 仍读取 `svc.Ctx` 字段（靠 EOF 自愈，非致命），彻底修需把它们改为参数传入 ctx（后续）。
- **restartOnExit 必须由 Stop/Kill/Remove 清零**：否则 Restart 后立即 Stop，进程退出仍会被 `wait()` 自动拉起，违背用户意图。


## 已知问题与半成品（勿误判为新引入的 bug）

### `go test ./...` 已知 FAIL（均非本次重构引入，勿误修）
- **vet 拦截 build failed**：[client.go:446](client/client.go#L446) `fmt.Sprintf`、[rocket.go:88](pkg/config/rocket.go#L88) `fmt.Errorf` 用了非常量格式串 → `go test` 默认 `go vet` 拦截（`go build` 正常通过）。修法：固定格式串。
- **`pkg.TestRand`**（rand_test.go）：测 `random.go`。
- 验证自己的改动用 `go test ./internal/cache/...`（本重构核心包，已通过）。

### 未完成 / 半成品（改动前先与 owner 确认预期）
- **Remove 仅删 server 缓存**：`cache.Server.Remove()` 只停进程 + `RemoveServer()`；script 配置删除 + 配置文件回写属 handler 层 + `UpdateConfig` module，[remove.go](api/handle/remove.go) 当前空壳未激活。
- **Always 自动重启未生效**：`wait()` 里基于信号的重启 switch 被整段注释，进程意外崩溃不自愈。
- **报警主链路未通**：`InitAlert` / `AlertMessage` / `probe.CheckHardWare` 关键循环被注释，全局 `msgCache` 无消费者。
- **dispatcher 锁 bug**（`pkg/config/dispatcher.go`）：`SetDispatcher` 用 `RLock` 应为 `Lock`；`CleanAlert` 持 `Lock` 时调 `GetDispatcherList()`（内部 `RLock`）会死锁。

### 后续技术债（已评估为非致命，未在本次修）
- **`ChStop` 无缓冲 chan**：同 svc 多个 goroutine 阻塞 `<-ChStop` 时，`canstop` 只发一个 token → 潜在死锁。`SetCanNotOperation` 已串行化同 svc 操作，实际触发罕见；彻底修改为 `close` 广播。
- （restart 重入的 ctx 字段覆盖：`appendRead`/`successAlert`/`cron` 读 `svc.Ctx` 字段靠 EOF 自愈——见上方【已知陷阱】）

## 新增一个进程操作 API 的步骤

1. `api/handle/<verb>.go`：写 handler，从 `xmux.GetInstance(r)` 取请求/参数。
2. `api/handle.go`：在对应 RouteGroup 注册路由（按权限级别选 Admin/Script/simple）。
3. `internal/cache/server.go`：实现 `Server.<Verb>()`，遵循状态机约定（注意跨平台 + Sync 陷阱）。
4. `client/command/<verb>.go`：给 scsctl 加 cobra 子命令（多节点并发用 `sync.WaitGroup`）。

## 深度审查清单（按需查阅）

改动涉及以下领域时，先读对应 skill 再动手：
- 报警链路（`pkg/config` 报警器/probe/dispatcher/monitor）→ `.trae/skills/scs-alert-review/SKILL.md`
- 并发与 channel（Server 状态机/全局 store/Logger/goroutine）→ `.trae/skills/scs-concurrency-review/SKILL.md`
- 进程生命周期（start/stop/restart/kill/cron/副本）→ `.trae/skills/scs-process-guide/SKILL.md`
