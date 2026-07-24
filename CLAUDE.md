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
- 日志：自研 `github.com/hyahm/golog`（⚠️ `(*Log).Sync()` 会无条件 close channel，**不可重复调用**，必须走 `Server.syncLogger()`）
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
  ├─ start.go             StartAsync / Start / syncLogger（已加 recover）
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

1. **`Ctx + Cancel`**（`context.WithCancel`）：取消日志 goroutine、cron ticker、进程本身（`exec.CommandContext(ctx, ...)`，ctx 取消即杀进程）。
2. **`Exit chan int`**（buffer 2）：信号路由。`9=kill / 10=restart / 11=stop / 12=remove`。从不显式 close，靠 `Ctx.Done()` + `cmd.Wait()` 退出。
3. **`Status.Status`**（`pkg/status.go`）：`STOP / RUNNING / WAITSTOP / WAITRESTART`。

操作流转：
- `StartAsync()` → 仅当 `Status==STOP` 时 `go Start()`
- `Start()` → 服务启动的时候调用新建 Logger + Exit chan + Ctx → cron 路径或普通 `start()` → `wait()`
- `Stop()` → 如果收到cannotstop信号，等待收到canstop才能 直接 Cancel
- `Kill()` → 如果收到cannotstop信号，等待收到canstop才能 直接 Cancel
- `Restart()` →  如果收到cannotstop信号，等待收到canstop才能 直接 Cancel 然后再StartAsync 启动
- `Remove()` → `Exit <- 12`

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

- **golog `Sync()` panic**：`(*Log).Sync()` 内部无条件 `close(channel)`，对同一 Logger 二次调用会 `panic: close of closed channel`。任何刷新日志的地方一律走 `svc.syncLogger()`（`start.go`，已加 `defer recover()`）。新增 defer 时确认不会对同一 Logger 重复触发。
- **dispatcher 锁 bug**（`pkg/config/dispatcher.go`）：`SetDispatcher` 用了 `RLock()` 应为 `Lock()`；`CleanAlert` 持有 `Lock()` 时调用 `GetDispatcherList()`（内部再 `RLock()`）会死锁。
- **报警主链路当前未通**：`InitAlert()` / `AlertMessage()` / `probe.CheckHardWare()` / `server.successAlert()` 中按类型分发/启动 ticker 的循环被注释；全局 `msgCache`（`internal/cache.go`）无消费者。若要"接通"报警，先与用户确认预期，不要假设被注释代码应当启用。
- **`Exit chan` 排空**：状态切换前若有积压信号会导致状态错乱，`Kill/Restart` 会先 `<-Exit` 排空。
- **cron 路径**：`IsCron=true` 时 `wait()` 不应把 `Status` 改成 STOP（保持 RUNNING 等下一次）。

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
