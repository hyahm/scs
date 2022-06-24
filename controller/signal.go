package controller

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/alert"
)

var signalHandle map[string]*pkg.SignalRequest
var mu sync.RWMutex

func init() {
	signalHandle = make(map[string]*pkg.SignalRequest)
	mu = sync.RWMutex{}
}

// 添加信号请求，如果添加成功返回true， 修改就是false
func AddSignalRequest(name string, sr *pkg.SignalRequest) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := signalHandle[name]; !ok {
		signalHandle[name] = sr
	}
}

// 添加信号请求，如果添加成功返回true， 修改就是false
func UpdateSignalRequest(name string, sr *pkg.SignalRequest) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := signalHandle[name]; ok {
		signalHandle[name].Parameter = sr.Parameter
		signalHandle[name].Notice = sr.Notice
		signalHandle[name].Restart = sr.Restart
		return true
	}
	return false
}

func GetSignalRequest(name string) *pkg.SignalRequest {
	mu.RLock()
	defer mu.RUnlock()
	if sr, ok := signalHandle[name]; ok {
		return sr
	}
	return nil
}

func DeleteSignalRequest(name string) {
	mu.Lock()
	defer mu.Unlock()
	delete(signalHandle, name)
}

func UnStop(ctx context.Context, name string, timeout time.Duration) {
	select {
	case <-time.After(time.Second * timeout):
		pkg.DeleteAtomSignal(name)
		sr := GetSignalRequest(name)
		if sr == nil {
			return
		}
		// 清除超时信号
		defer DeleteSignalRequest(name)
		// 报警
		if sr.Notice {
			ra := &alert.RespAlert{
				Name:   name,
				Title:  "原子操作超时",
				Reason: fmt.Sprintf("原子操作超时超过 %d 秒没有执行完成", sr.Timeout),
			}
			if sr.ContinuityInterval > 0 {
				ra.ContinuityInterval = sr.ContinuityInterval
			}
			ra.SendAlert()
		}
		if sr.Restart {
			if server, ok := store.Store.GetServerByName(name); ok {
				golog.Info("restart atom")
				KillAndStartServer(sr.Parameter, server)
			}
		}
	case <-ctx.Done():

	}
}

// 保存是否可停止的状态

// type SNS struct {
// 	timeout int64           // 保存超时时间
// 	ctx     context.Context // 上下文管理
// 	cancel  context.CancelFunc
// }

// var StopSignal map[string]SNS
// var signalMu sync.RWMutex
// var signal chan string

// func init() {
// 	StopSignal = make(map[string]SNS)
// 	signalMu = sync.RWMutex{}
// 	signal = make(chan string, 1000)
// }

// // 启动主线程的时候就启动这个
// func StartSignalTime() {
// 	for {
// 		select {
// 		case name := <-signal:
// 			// 把这个key删掉，并且结束对应的goroutine
// 			removeSignal(name)
// 		}
// 	}
// }

// func removeSignal(name string) {
// 	signalMu.Lock()
// 	defer signalMu.Unlock()
// 	// 为了避免短暂的操作， 还是先判断下是否存在这个key
// 	if v, ok := StopSignal[name]; ok {
// 		// 先停止goroutine
// 		v.cancel()
// 	}

// }

// // 启动一个不能停止的信号
// func StartCanNotStop(name string, timeout int64) {
// 	// name: name
// 	// timeout: 超时时间， 超过这个时间没收到停止信号就发送报警
// 	if timeout == 0 {
// 		return
// 	}
// 	ctx, cancel := context.WithCancel(context.Background())
// 	sns := SNS{timeout: timeout, ctx: ctx, cancel: cancel}
// 	signalMu.Lock()
// 	defer signalMu.Unlock()

// 	if _, ok := StopSignal[name]; ok {
// 		// 如果存在这个key就是重复提交
// 		return
// 	}

// 	StopSignal[name] = sns
// 	go func(sns SNS) {
// 		// 开一个协程
// 		select {
// 		case <-time.After(time.Second * time.Duration(sns.timeout)):
// 			//  超时了发报警
// 			// 还是要判断一下，这个是不是运行状态
// 			fmt.Println("超时了")
// 		case <-sns.ctx.Done():
// 			// 然后删掉这个key
// 			delete(StopSignal, name)
// 		}
// 	}(sns)

// }

// // 启动一个可用停止的信号
// func StartCanStop(name string) {
// 	signalMu.RLock()
// 	defer signalMu.RUnlock()
// 	if _, ok := StopSignal[name]; ok {
// 		signal <- name
// 	}

// }
