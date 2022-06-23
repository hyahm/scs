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
func AddSignalRequest(name string, sr *pkg.SignalRequest) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := signalHandle[name]; ok {
		signalHandle[name].Parameter = sr.Parameter
		signalHandle[name].Notice = sr.Notice
		signalHandle[name].Restart = sr.Restart
		return false
	}
	signalHandle[name] = sr
	return true
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
