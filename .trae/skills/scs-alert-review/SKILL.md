---
name: "scs-alert-review"
description: "SCS 报警链路审查。当修改 pkg/config 下的报警器（email/telegram/dingding/weixin/rocket/callback）、dispatcher 去抖、probe 探测器、monitor 心跳时调用，用于识别半成品代码与重构陷阱。"
---

# SCS 报警链路审查

SCS 的报警系统当前处于**半重构状态**，很多关键循环被注释，改动前必须识别"就位"与"未通"的部分。

## 当前状态地图

### 已就位（可正常工作）
- **配置结构**：`pkg/config/alert.go Alert{Email, Rocket, Telegram, WeiXin, Callback, DingDing}`
- **接口**：`pkg/message/sendalerter.go SendAlerter.Send(body, to...)`
- **消息体**：`pkg/message/message.go Message`，带 `FormatBody(format)` 模板渲染
- **各通知器实现**：`email.go / telegram.go / weixin.go / dingding.go / rocket.go / callback.go`
- **dispatcher 去抖**：`pkg/config/dispatcher.go alertInfoMap`，`RespAlert.SendAlert()` 实现"首次立即发，后续需超过 ContinuityInterval"
- **探测器框架**：`pkg/config/probe.go Probe{Mem, Cpu, Disk, Interval, ContinuityInterval}`
- **HTTP 入口**：`api/handle/alert.go Alert()`（POST `/send/alert`）、`GetAlert()`
- **监控心跳**：`pkg/config/monitor.go` 对 Monitor 列表做 `/probe` 心跳（3 次重试 + telnet）

### 未通（关键循环被注释，报警主链路当前不工作）
- **`pkg/config/alert.go InitAlert()`**：按 alert 类型注册到 `alerter.Alerts` map 的逻辑**被整段注释**
- **`pkg/config/send.go AlertMessage()`**：按 `switch alert.(type)` 分发到各通知器的循环**被注释**，只剩 `golog.Errorf`
- **`pkg/config/probe.go initConfig() / CheckHardWare()`**：按 Cpu/Mem/Disk 分别 NewCheckPointer 并启动 ticker 的逻辑**被注释**
- **`internal/cache/server.go successAlert()`**：恢复通知的旧 AlertMessage 调用**被注释**，几乎是空壳
- **全局 `msgCache` channel**：`internal/cache.go init()` 创建了容量 1000 的缓冲，**没有消费者 goroutine**

## 审查任务

### 任务 A：确认改动是否触及"未通"代码
如果用户的修改涉及：
- `InitAlert` / `AlertMessage` / `successAlert` / `msgCache` / `initConfig` / `CheckHardWare`

→ **先与用户确认预期行为**：是要"接通"报警主链路，还是仅调整已就位的部分？不要假设被注释的代码应该启用。

### 任务 B：检查 dispatcher 锁（★★★ 已知 bug）
`pkg/config/dispatcher.go`：
- `SetDispatcher` 用了 `RLock()` 应为 `Lock()`
- `CleanAlert` 持有 `Lock()` 时调用 `GetDispatcherList()`（内部再 `RLock()`）会死锁

修复方向：
```go
func (a *alertInfoMap) SetDispatcher(key string, info AlertInfo) {
    a.Lock()          // 改这里
    defer a.Unlock()
    a.dispatcher[key] = info
}

// CleanAlert 不要在持有锁时调用会再次加锁的方法
func CleanAlert() {
    dispatcher.Lock()
    list := make([]string, 0, len(dispatcher.dispatcher))
    for k := range dispatcher.dispatcher {  // 内联，避免再调 GetDispatcherList
        list = append(list, k)
    }
    dispatcher.Unlock()
    // ... 后续清理
}
```

### 任务 C：新增报警后端
1. 在 `pkg/config/` 下新建 `<name>.go`，定义结构体（参考 `email.go AlertEmail`）
2. 实现 `SendAlerter` 接口（`Send(body, to...) error`）
3. 在 `Alert` 结构体（`alert.go`）加字段
4. 在 `InitAlert`（当前被注释，需要先与 owner 确认是否启用）注册
5. 在 `AlertMessage` 的分发循环（当前被注释）加 case
6. 在 `Message.FormatBody` 里加格式模板

### 任务 D：新增探测器
1. `pkg/config/probe.go Probe` 加字段
2. 实现 `CheckPointer` 接口（`Check() / Update()`，参考 `cpu.go / mem.go / disk.go`）
3. 在 `initConfig`（当前被注释）注册 ticker

## 输出约定
审查时明确区分：
- ✅ **就位可用**：明确指出文件与函数
- ⚠️ **半成品**：指出被注释的位置，提示需要 owner 确认
- 🐛 **已知 bug**：dispatcher 锁问题等
