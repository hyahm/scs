package controller

import (
	"errors"
	"sync/atomic"

	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg/config/scripts/subname"
)

// 只有在删除的时候才会需要   svc.StopSignal 信号
// 只有在 Server.Removed 为 true的时候才会发送  svc.StopSignal 信号
// RemovePname 的时候才会用到
func RemoveScript(pname string) error {

	if _, ok := store.ss[pname]; ok {
		replicate := store.ss[pname].Replicate
		if replicate == 0 {
			replicate = 1
		}

		for i := 0; i < replicate; i++ {
			subname := subname.NewSubname(pname, i)
			atomic.AddInt64(&global.CanReload, 1)
			go Remove(store.servers[subname.String()], true)
		}

	} else {
		return errors.New("not found this pname:" + pname)
	}
	return nil
}

// update: 是否需要重新修改配置文件， 有锁
func Remove(svc *server.Server, update bool) {
	// 如果是always 为 true，那么直接修改为false
	svc.Remove()
	<-svc.StopSignal
	store.mu.Lock()
	delete(store.serverIndex[svc.Name], svc.Index)
	delete(store.servers, svc.SubName)
	removeServer(svc.Name, svc.SubName, update)

	// 如果全部删光了， 那么scripts的name也要删除
	if len(store.serverIndex[svc.Name]) == 0 {
		delete(store.ss, svc.Name)
	}
	store.mu.Unlock()
	atomic.AddInt64(&global.CanReload, -1)
}

// remove没有锁
// func remove(svc *server.Server, update bool) {
// 	// 如果是always 为 true，那么直接修改为false
// 	svc.Remove()
// 	<-svc.StopSignal
// 	delete(store.serverIndex[svc.Name], svc.Index)
// 	delete(store.servers, svc.SubName)
// 	removeServer(svc.Name, svc.SubName, update)
// 	// 如果全部删光了， 那么scripts的name也要删除
// 	if len(store.serverIndex[svc.Name]) == 0 {
// 		delete(store.ss, svc.Name)
// 	}
// 	atomic.AddInt64(&global.CanReload, -1)
// }
