# SCS 项目协作提示词（System Prompt）

> 把这一段放进对话的 system / 首条 user 消息里，AI 即可获得项目全局上下文。

---

你正在协作的项目是 **SCS（Service Control Service）**，位于 `d:\cander\scs`。

## 项目定位
SCS 是一个用 Go 1.25 编写的**跨平台进程编排与监控守护进程**（`scsd`）+ 配套 CLI（`scsctl`）。
通过 YAML 描述一组脚本/服务，管理其生命周期：启动/停止/重启/kill/更新/副本/定时任务/监控报警。
定位类似 `supervisor`，但内置 TLS、多通道报警、原子性停止（cannotstop）、WebSocket 实时日志。

## 技术栈
- 语言：Go 1.25（module: `github.com/hyahm/scs`），当前版本 v3.8.5（`global/var.go`）
- Web 框架：自研 `github.com/hyahm/xmux`（支持分组、token 权限 keys、中间件 module、WebSocket）
- CLI：`github.com/spf13/cobra`
- 日志：自研 `github.com/hyahm/golog`（注意：`(*Log).Sync()` 会无条件 close channel，不可重复调用）
- 配置：`gopkg.in/yaml.v3`
- 监控：`github.com/shirou/gopsutil`（CPU/内存/磁盘）
- 报警后端：gomail（邮件）/ Rocket.Chat / Telegram / 企业微信 / 钉钉 / HTTP 回调

## 架构分层
```
api/handle.go         三级 token 权限路由（AdminHandle / ScriptHandle / simpleHandle）
api/server.go         HttpServer() 入口（监听 + TLS）
internal/cache/       ★ 进程管理核心
  ├─ server.go          Server 结构体（单个运行副本）
  ├─ server_store.go    serverStore 全局单例（map[string]*Server + RWMutex）
  ├─ script_store.go    Script 配置态仓库（map[string]config.Script + RWMutex）
  ├─ start.go/wait.go   启动 / 等待退出
  ├─ cron.go            定时任务循环
  ├─ log.go             stdout/stderr 捕获 → Logger
  ├─ script_windows.go  Windows: powershell -c + taskkill
  └─ script_unix.go     Unix: /bin/sh -c + syscall.Kill(-pgid, SIGKILL)
internal/config.go     ReadConfig / WriteConfig（YAML 读写）
pkg/config/            配置结构体 + 探测器 + 报警器实现
pkg/message/           报警消息体 + SendAlerter 接口
client/command/        scsctl 的 cobra 子命令
controller/            旧三层架构，几乎全部被注释（死代码）
```

## Server 状态机（重要）
每个 `Server` 副本通过 3 个并发原语协作：
- `Ctx context.Context + Cancel context.CancelFunc`：取消日志 goroutine、cron ticker、CommandContext 进程
- `Exit chan int`（buffer 2）：信号路由，9=kill / 10=restart / 11=stop / 12=remove
- `Status.Status`：`STOP / RUNNING / WAITSTOP / WAITRESTART`（常量在 `pkg/status.go`）

操作流转：
- `StartAsync()` → 仅当 `Status==STOP` 时 `go Start()`
- `Start()` → 新建 Logger + Exit chan + Ctx → 走 cron 或普通 `start()` → `wait()`
- `Stop()` → cron 直接 Cancel；RUNNING 走 `Cancel()`（context 取消会杀进程）
- `Kill()` → RUNNING 时 `Exit <- 10` + `kill()`；WAITRESTART/WAITSTOP 先 `<-Exit` 排空
- `Restart()` → 通过 `Exit chan` 传递 10
- `Remove()` → `Exit <- 12`

## 协作约定
1. **中文优先**：注释、错误信息、commit message、文档全部用中文。
2. **跨平台**：任何进程操作都要同时考虑 `script_windows.go` 和 `script_unix.go`，用 build tag 区分。
3. **不要碰 controller/**：整层是被注释的死代码，新功能一律放 `internal/cache` + `api/handle`。
4. **golog 坑**：调用 `Logger.Sync()` 必须用 `syncLogger()`（已加 recover），不能直接裸调。
5. **并发安全**：访问全局 store（serverStore / Script / InConfig / dispatcher）必须走它们的 RWMutex 包装方法，不要直接读 map。
6. **重构痕迹**：大量旧代码被注释保留（wait.go / send.go / probe.go / global/var.go），改动前先确认目标代码是否已被注释替代。
7. **报警链路当前是半成品**：`InitAlert` 和 `AlertMessage` 的关键分发循环被整段注释，只有 dispatcher（去抖）和探测器框架就位——改报警相关功能前先与 owner 确认预期行为。

## 常用命令
```powershell
# 编译（Windows）
.\build.bat
# 编译（Linux）
bash update.sh
# 单独构建
go build ./...
go test ./...
# 运行单元测试（项目里测试很少，主要在 internal/cache 和 pkg）
go test ./internal/... ./pkg/...
```

## 已知风险点（改这些地方要格外小心）
- `pkg/config/dispatcher.go SetDispatcher` 用了 `RLock` 应为 `Lock`；`CleanAlert` 持有 `Lock` 时又调 `GetDispatcherList()`（内部 RLock），同 goroutine 重入会死锁。
- 全局 `msgCache` channel（容量 1000）没有消费者 goroutine，是预留管道。
- `Server.Ready chan bool` 创建了但从未收发（CheckReady 未完成）。
- `appendRead`（log.go 方法版）依赖 `Ctx.Done()` 退出，但 `read`（包级函数版，用于 shell 命令）不监听 Ctx，仅靠 EOF 退出——短命令可接受，长命令需注意。
