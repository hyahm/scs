---
name: "scs-concurrency-review"
description: "SCS 并发与 channel 安全审查。当修改 Server 状态机、全局 store（serverStore/Script/dispatcher）、restartOnExit 标志、Logger 或新增 goroutine 时调用，用于排查死锁、channel panic、goroutine 泄漏。"
---

# SCS 并发与 channel 安全审查

本 skill 用于审查 `d:\cander\scs` 项目里涉及并发的改动，避免以下高频问题：

## 审查清单

### 1. golog 的 Sync 陷阱（★★★ 已踩过坑）
**问题**：`golog.(*Log).Sync()` 内部 `close(l.task.cache)` 是无条件的，对同一 Logger 实例二次调用会 `panic: close of closed channel`。restart 路径下两次 `Start()` goroutine 重叠会覆盖 `svc.Logger`，旧 goroutine 的 defer 若读字段会误刷新进程日志。

**检查点**：
- `Start()` 内用**局部变量捕获 logger** + `defer func(){ defer func(){ _ = recover() }(); logger.Sync() }()`，确保只同步本次实例
- 新增 defer 时确认绑定的是局部变量，而非 `svc.Logger` 字段
- cron / always restart / 重启路径下 `Start()` 会被多次进入，Logger 会被重新赋值，旧 Logger 的 Sync 不能被重复触发

**修复模板**（`start.go` 当前实现）：
```go
logger := golog.NewLog(...)
svc.Logger = logger
defer func() {
    defer func() { _ = recover() }()
    logger.Sync()
}()
```

### 2. 全局 store 锁（★★★ 已知有 bug）
全局单例 store 都用 `sync.RWMutex` 包装，但有几处已知问题：

**`pkg/config/dispatcher.go`**：
- `SetDispatcher` 用了 `dispatcher.RLock()` —— **应为 `Lock()`**（写操作）
- `CleanAlert` 持有 `Lock()` 时调用 `GetDispatcherList()`（内部 `RLock()`）—— 同 goroutine 重入会死锁

**审查点**：
- 写 map 必须用 `Lock()`，读用 `RLock()`
- 持有锁时不要回调会再次获取同一把锁的方法
- store 方法要成对提供 `GetXxx` / `SetXxx`，外部不要绕过它们直接访问 map

### 3. channel 创建与关闭
| channel | 创建位置 | 关闭策略 |
|---------|----------|----------|
| `Server.Ready chan bool` | `factory.go:32` | 从未收发（CheckReady 未完成），改动时确认 |
| `Status.ChStop chan struct{}` | `script.go:52` | 用于 cannotstop 原子停止，未 close |
| 全局 `msgCache chan message.Message` | `internal/cache.go init()` | **无消费者**，发送方要注意非阻塞写 |

> 注：`Server.Exit chan int`（旧的 9/10/11/12 信号路由）已删除，改为 `restartOnExit bool`（`svc.mu` 保护）+ `Cancel()` 协作。

**审查点**：
- 新增 channel 必须明确：谁创建、谁 close、谁接收
- 发送端用 `select { case ch <- v: default: }` 防止阻塞（项目里 `golog.s()` 就是这个模式）
- 不要对已关闭 channel 发送（panic）/ 关闭已关闭 channel（panic）/ 重复 close（panic）

### 4. goroutine 泄漏
项目里常见的 goroutine：
- `StartAsync` → `go Start()`
- `read()` → 2 个 `appendRead` goroutine（监听 `Ctx.Done()`）
- `shell` 类命令 → `read` 包级函数启动的 goroutine（**不监听 Ctx**，靠 EOF 退出）
- `successAlert()` → `wait()` 里 `go svc.successAlert()`
- client 层多节点并发 → `sync.WaitGroup` + `go func()`

**审查点**：
- 每个无穷循环 goroutine 必须有退出条件（`Ctx.Done()` / EOF / channel close）
- HTTP handler 启动的 goroutine 要注意请求结束后是否还在跑
- `successAlert` 当前几乎是空壳，改动时确认是否有定时器未停

### 5. context 取消传播
- `Server.Ctx` 被 `Cancel()` 取消后：
  - `exec.CommandContext` 会杀进程
  - `appendRead` 的 `select` 会退出
  - `cron` 的 `select` 会退出
- 审查点：新增的 goroutine 如果持有 `Server.Ctx`，必须监听 `Done()`

## 审查输出格式
对每个发现的潜在问题，按以下格式报告：
```
[严重程度] 文件:行号
问题：简述
触发场景：什么情况下会出问题
修复建议：具体代码改动方向
```
严重程度分级：★★★ panic / 死锁 ｜ ★★ 数据竞争 ｜ ★ 潜在风险
