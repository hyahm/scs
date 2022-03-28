package controller

import (
	"errors"
	"sync/atomic"

	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg/config/scripts/subname"
	"github.com/hyahm/scs/status"
)

// 只有在删除的时候才会需要   svc.StopSigle 信号
// 只有在 Server.Removed 为 true的时候才会发送  svc.StopSigle 信号
// RemovePname 的时候才会用到
func RemoveScript(pname string) error {
	mu.RLock()
	defer mu.RUnlock()

	if _, ok := ss[pname]; ok {
		replicate := ss[pname].Replicate
		if replicate == 0 {
			replicate = 1
		}

		for i := 0; i < replicate; i++ {
			subname := subname.NewSubname(pname, i)
			atomic.AddInt64(&global.CanReload, 1)
			go Remove(servers[subname.String()], true)
		}

	} else {
		return errors.New("not found this pname:" + pname)
	}
	return nil
}

// update: 是否需要重新修改配置文件
func Remove(svc *server.Server, update bool) {
	// 如果是always 为 true，那么直接修改为false
	if svc == nil {
		return
	}
	if svc.Always {
		svc.Always = false
	}
	svc.Removed = true
	if svc.Status.Status != status.STOP && svc.IsCron {
		svc.Remove()
		<-svc.StopSigle
	}

	mu.Lock()
	delete(serverIndex[svc.Name], svc.Index)
	removeServer(svc.SubName, update)
	mu.Unlock()
	atomic.AddInt64(&global.CanReload, -1)
}

// func RemoveAllScripts() {
// 	// 删除所有脚本
// 	config.RemoveAllScriptToConfigFile()
// 	mu.RLock()
// 	defer mu.RUnlock()

// 	for _, svc := range servers {
// 		replicate := svc.Replicate
// 		for i := 0; i < replicate; i++ {
// 			atomic.AddInt64(&global.CanReload, 1)
// 			go Remove(svc)
// 		}
// 	}
// }
