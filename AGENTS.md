# AGENTS.md

SCS（Service Control Service）—— Go 编写的跨平台进程编排与监控守护进程，定位类似 supervisor，内置 TLS、多通道报警、原子性停止、WebSocket 实时日志。

## 项目概览

- **Module**: `github.com/hyahm/scs` v3.8.5
- **Go 版本**: 1.25+
- **入口**:
  - `cmd/scsd/scsd.go` — 守护进程，`go run ./cmd/scsd/scsd.go -f scs.yaml`
  - `cmd/scsctl/scsctl.go` — cobra CLI，多节点并发管理
- **配置文件**: `scs.yaml`（默认），结构见 `default.yaml`

## 技术栈

| 组件 | 库 |
|------|-----|
| Web 框架 | `github.com/hyahm/xmux`（自研，分组路由 + token 权限 + 中间件 + WebSocket） |
| CLI | `github.com/spf13/cobra` |
| 日志 | `github.com/hyahm/golog`（自研） |
| 配置 | `gopkg.in/yaml.v2` |
| 监控 | `github.com/shirou/gopsutil`（CPU/内存/磁盘） |
| 报警 | gomail / Rocket.Chat / Telegram / 企业微信 / 钉钉 / HTTP 回调 |

## 目录结构

```
cmd/scsd/scsd.go          守护进程入口
cmd/scsctl/scsctl.go      CLI 入口
api/
  server.go               HttpServer（监听 + TLS）
  handle.go               三级路由注册（AdminHandle > ScriptHandle > simpleHandle）
  handle/*.go             21 个 HTTP handler
  module/*.go             中间件（CheckToken / UpdateConfig）
internal/
  config.go               配置读写（ReadConfig / WriteConfig / AddScriptToConfigFile）
  cache.go                报警消息缓存通道 + 消费者（PushAlert / StartAlertConsumer）
  cache/                  ★ 进程管理核心
    server.go             Server 结构体 + 状态机（Start/Stop/Kill/Restart/Remove/UpdateServer）
    server_store.go       全局 Server 单例（map[string]*Server + RWMutex）
    script_store.go       全局 Script 配置仓库（map[string]config.Script + RWMutex）
    script.go             AddScript 工厂 + newCommand 跨平台封装
    start.go              StartAsync / Start / FirstStartAllScript
    wait.go               wait() 进程退出等待 + Always 自动重启
    cron.go               定时任务循环
    factory.go            工厂函数 ServerOption
    log.go                stdout/stderr 实时捕获
    prestart.go           启动前 PreStart 检查
    script_unix.go        Unix: /bin/sh -c + syscall.Kill(-pgid, SIGKILL)
    script_windows.go     Windows: powershell -c + taskkill /F /T /PID
pkg/
  status.go               状态常量 STOP/RUNNING/WAITSTOP/WAITRESTART/INSTALL/STARTING/REMOVING
  config/                 配置结构体 + 报警器 + 探测器 + dispatcher 去抖
  message/                报警消息体 + SendAlerter 接口
client/                   scsctl HTTP 客户端 + 20 个 cobra 子命令
controller/               ⚠️ 仅 auth.go 和 signal.go 仍在使用，其余为死代码，禁止新增
global/var.go             VERSION 常量
```

## 启动流程

```
main()
  → golog.Sync defer
  → go message.GetIp()            异步获取公网 IP
  → internal.Setrlimit()          设置文件描述符上限
  → flag.Parse()                  解析 -f / -v
  → signal.Notify                 SIGTERM/SIGINT/SIGPIPE → os.Exit(1)
  → internal.ReadConfig()         读取 YAML → InConfig 单例
  → config.InitAlert()            初始化报警通道
  → internal.StartAlertConsumer() 启动报警消费协程
  → go config.CleanAlert()        启动 dispatcher 清理
  → config.InitDetector()         初始化硬件检测器
  → go config.CheckHardWare()     启动硬件检测（可选）
  → cache.FirstStartAllScript()   遍历配置启动所有脚本
  → internal.WriteConfig()        回写配置
  → api.HttpServer()              启动 HTTP 服务（阻塞）
```

## Server 状态机（核心）

### 状态常量

```
STOP → RUNNING → WAITSTOP / WAITRESTART → STOP
       INSTALL → RUNNING
       STARTING → RUNNING
       REMOVING → STOP
```

### 三个并发原语

1. **Ctx + Cancel**（`context.WithCancel`）：取消进程、日志 goroutine、cron ticker。`Start()` 内用局部变量捕获，defer cancel，避免 restart 重入时旧 goroutine 误杀新进程。
2. **restartOnExit**：Restart 的"退出后重启"意图。`Restart()` 在 RUNNING 时置 true + Cancel()；`wait()` 进程退出后消费并 StartAsync()；Stop/Kill/Remove 入口清零。
3. **Status.Status**：当前运行状态。

### 操作流转

| 操作 | 流程 |
|------|------|
| `StartAsync()` | 仅 Status==STOP 时 `go Start()` |
| `Start()` | 局部捕获 logger + ctx/cancel → cron 路径或普通 start() → wait() |
| `Stop()` | 清 restartOnExit → 等 cannotstop → Cancel() |
| `Kill()` | 清 restartOnExit → 等 cannotstop → Cancel() + kill() 强杀进程组 |
| `Restart()` | 等 cannotstop → STOP 直接 StartAsync()；RUNNING 置 restartOnExit + Cancel() |
| `Remove()` | 清 restartOnExit → 等 cannotstop → Cancel() + kill() + RemoveServer() |

### cannotstop（原子停止）

- `/cannotstop/name` → `CanNotStop = true`
- `/canstop/name` → `CanNotStop = false` + `close(ChStop)` 广播 + `make(chan struct{})` 重建
- Stop/Kill/Restart/Remove 入口 `if CanNotStop { <-ChStop }` 阻塞等待

### Always 自动重启

- 进程异常退出 → `wait()` 设置 `broken = true` → `successAlert()` 3 秒后发恢复通知
- `Restart()` 设置 `restartOnExit = true` → `wait()` 消费后 `StartAsync()`
- Stop/Kill/Remove 入口清零 `restartOnExit`，防止误拉起

## HTTP 路由与权限

```
AdminHandle (key=admin, +CheckToken)
  GET /-/reload        重载配置
  GET /-/fmt           格式化回写配置
  GET /-/config        查看完整配置
  GET /get/alert       获取报警列表
  GET /get/alarms      获取硬件报警
  GET /get/repo        获取仓库信息
  GET /script          添加脚本（BindJson + UpdateConfig）
  GET /enable/name     启用脚本
  GET /disable/name    禁用脚本
  └─ ScriptHandle (key=script)
       GET /stop/name      停止
       GET /stop           停止全部
       GET /kill/name      强杀
       GET /get/servers    获取所有 server
       GET /server/info/name 获取 server 详情
       GET /cannotstop/name  原子停止
       GET /canstop/name     解除原子停止
       GET /parameter/name   设置 SignalRequest
       GET /get/scripts     获取所有脚本
       GET /env/name        获取环境变量
       GET /remove/name     删除脚本（+UpdateConfig）
       GET /send/alert      发送报警
       └─ simpleHandle (key=simple)
            GET /status/name    状态
            GET /start/name     启动
            GET /update/name    更新（git pull）
            GET /update         更新全部
            GET /start          启动全部
            GET /status         全部状态
            GET /log/name       WebSocket 实时日志
            POST /restart/name  重启
            POST /restart       重启全部
```

## 报警系统

### 6 个通道

Email / Rocket.Chat / Telegram / 企业微信 / 钉钉 / HTTP 回调

### 链路

```
硬件检测器(CPU/Mem/Disk/Monitor) → AlertInfo.BreakDown()/Recover()
HTTP /send/alert → RespAlert.SendAlert()
进程崩溃 → wait() 设置 broken=true → successAlert() 恢复通知
```

### 去抖

`dispatcher` 按 `ContinuityInterval` 控制重复报警间隔，`CleanAlert()` 每 10 分钟清理超 10 小时未发送的条目。

### 消息缓存

`internal.PushAlert()` → `msgCache`（1000 缓冲）→ `StartAlertConsumer()` 异步消费。

## 协作约定

1. **中文优先**：注释、错误信息、commit message 全部用中文。
2. **跨平台**：进程操作同时考虑 `script_windows.go` 和 `script_unix.go`，不可写死平台。
3. **新功能位置**：`internal/cache` + `api/handle`，禁止在 `controller/` 新增代码。
4. **全局 store 访问**：`serverStore`/`script_store`/`dispatcher` 都用 `sync.RWMutex`，写 Lock()，读 RLock()，**持有锁时不可回调会再次获取同一把锁的方法**。
5. **函数不超过 50 行**，**导出函数必须有 godoc 注释**。

## 已知陷阱

- **golog Sync() panic**：`(*Log).Sync()` 内部无条件 close(channel)，二次调用会 panic。`Start()` 内用局部变量捕获 logger + defer 确保只同步本次实例。
- **Restart 重入字段覆盖**：`Start()` 局部捕获 ctx/cancel/logger 已根治，但 `appendRead`/`successAlert`/`cron` 仍读 `svc.Ctx` 字段（靠 EOF 自愈）。
- **restartOnExit 必须由 Stop/Kill/Remove 清零**：否则 Restart 后立即 Stop，进程退出仍会被 wait() 自动拉起。

## 新增 API 步骤

1. `api/handle/<verb>.go`：写 handler，从 `xmux.GetInstance(r)` 取请求参数。
2. `api/handle.go`：在对应 RouteGroup 注册路由（AdminHandle / ScriptHandle / simpleHandle）。
3. `internal/cache/server.go`：实现 `Server.<Verb>()`，遵守状态机约定。
4. `client/command/<verb>.go`：给 scsctl 加 cobra 子命令。

## 常用命令

```bash
go run ./cmd/scsd/scsd.go -f scs.yaml    # 启动守护进程
go run ./cmd/scsd/scsd.go -v             # 查看版本
go run ./cmd/scsctl/scsctl.go            # 运行 CLI
go vet ./...                             # 静态检查
go build ./...                           # 编译检查
go test ./...                            # 测试（TestRand 为已知预失败）
```