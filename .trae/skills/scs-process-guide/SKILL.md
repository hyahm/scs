---
name: "scs-process-guide"
description: "SCS 进程管理操作指南。当需要新增/修改 scsd 对脚本的生命周期管理（启动、停止、重启、kill、副本、定时任务）或调整 Server 状态机时调用。"
---

# SCS 进程管理操作指南

本指南用于在 `d:\cander\scs` 项目里安全地修改进程生命周期相关逻辑。

## 核心文件
- 状态机入口：`internal/cache/server.go`（`Start/Stop/Kill/Restart/Remove/UpdateServer`）
- 启动流程：`internal/cache/start.go`（`StartAsync` / `Start`，局部捕获 ctx/logger）
- 等待退出：`internal/cache/wait.go`
- 定时任务：`internal/cache/cron.go`（`cron` / `doTicker`）
- 跨平台实现：
  - Windows：`internal/cache/script_windows.go`（`powershell -c` + `taskkill /F /T /PID`）
  - Unix：`internal/cache/script_unix.go`（`/bin/sh -c` + `syscall.Kill(-pgid, SIGKILL)`）
- 日志捕获：`internal/cache/log.go`（`read` 成员方法 + `read` 包级函数）

## Server 三大并发原语（务必理解）
1. **`Ctx + Cancel`**（`context.WithCancel`）：用于取消
   - 日志 goroutine（`log.go appendRead` 监听 `Ctx.Done()`）
   - cron ticker（`cron.go` 的 `select`）
   - 进程本身（`exec.CommandContext(ctx, ...)`，ctx 取消 → 进程被杀）
2. **`restartOnExit bool`**（`svc.mu` 保护）：Restart 的"退出后重启"意图
   - RUNNING 时 `Restart()` 置 true + `Cancel()`，`wait()` 进程退出后消费并 `StartAsync()`
   - `Stop/Kill/Remove` 入口必须 `setRestartOnExit(false)`，否则 Restart 后 Stop 仍会被拉起
3. **`Status.Status`**（`pkg/status.go` 常量）：`STOP / RUNNING / WAITSTOP / WAITRESTART`

## 修改 Check List
改任何进程操作前，依次确认：
- [ ] 是否同时处理了 Windows 和 Unix？（用 build tag 拆分到 `script_*.go`）
- [ ] 是否会触发 `Logger.Sync()` 重复调用？→ `Start()` 内已局部捕获 logger + defer 同步，新增同步照此模式
- [ ] goroutine 退出路径：是否有 `Ctx.Done()` 分支？是否会泄漏？
- [ ] Stop/Kill/Remove 是否清零了 `restartOnExit`？Restart 是否按 Status 分支（STOP 直接 StartAsync，RUNNING 走标志）？
- [ ] 访问 `serverStore` / `Script` 是否通过它们的 RWMutex 包装方法？
- [ ] cron 路径：`IsCron=true` 时 `wait()` 不应改 `Status` 为 STOP（保持 RUNNING 等下次）

## 常见任务模板

### 新增一个进程操作 API
1. `api/handle/<verb>.go`：写 handler，从 `xmux.GetInstance(r)` 拿请求/参数
2. `api/handle.go`：在对应 RouteGroup（AdminHandle/ScriptHandle/simpleHandle）注册路由
3. `internal/cache/server.go`：实现 `Server.<Verb>()` 方法，遵循状态机约定
4. `client/command/<verb>.go`：给 scsctl 加 cobra 子命令（通常并发请求多节点，用 `sync.WaitGroup`）

### 修改启动命令拼接
- 入口在 `script_windows.go` / `script_unix.go` 的 `start()` 里 `exec.CommandContext(...)`
- 环境变量注入：`svc.Cmd.Env = append(..., k+"="+v)`，注意 Windows 下 `\x00` 已被过滤
- 工作目录：`svc.Cmd.Dir = svc.Dir`

### 调整定时任务
- 配置结构：`pkg/config/cron.go Cron{Start, Loop, IsMonth, Times, LoopTime, First}`
- 循环体：`cron.go` 的 `for { select { case <-Ctx.Done(): ...; case <-ticker.C: doTicker() } }`
- `doTicker` 会 `Times--`，到 0 自动 `stopStatus()`

## 不要做的事
- ❌ 直接 `svc.Logger.Sync()` —— golog 的 Sync 会 close channel，重复调用 panic。`Start()` 内已用局部捕获 logger + `defer func(){ recover(); logger.Sync() }()`，新增同步逻辑照此模式
- ❌ 在 Stop/Kill/Remove 里忘记 `setRestartOnExit(false)` —— 否则 Restart 后 Stop 仍会被 wait() 自动拉起
- ❌ 在 `controller/` 下加新代码 —— 整层是死代码
- ❌ 给 `msgCache` 加发送方但不加消费者 —— 会阻塞或丢弃

## 验证
```powershell
go build ./...
go test ./internal/cache/... ./pkg/...
```
跨平台编译验证（可选）：
```powershell
$env:GOOS="linux"; go build ./...
$env:GOOS="windows"; go build ./...
```
