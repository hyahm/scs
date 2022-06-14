package controller

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
)

// 只有在删除的时候才会需要   svc.StopSignal 信号
// 只有在 Server.Removed 为 true的时候才会发送  svc.StopSignal 信号
// RemovePname 的时候才会用到
func RemoveScript(pname string) error {
	_, sok := store.Store.GetScriptByName(pname)
	if !sok {
		return errors.New("not found this pname:" + pname)
	}
	wg := &sync.WaitGroup{}
	for index := range store.Store.GetScriptIndex(pname) {
		wg.Add(1)
		subname := fmt.Sprintf("%s_%d", pname, index)
		svc, ok := store.Store.GetServerByName(subname)
		if !ok {
			golog.Error(pkg.ErrBugMsg)
			continue
		}
		go remove(svc, true, wg)
	}
	wg.Wait()
	atomic.AddInt64(&global.CanReload, -1)
	return nil
}

// update: 是否需要重新修改配置文件， 有锁
func Remove(svc *server.Server, update bool) {
	// 如果是always 为 true，那么直接修改为false
	golog.Info("svc start removed")
	svc.Remove()
	golog.Info("svc removed")
	<-svc.StopSignal
	store.Store.DeleteScriptIndex(svc.Name, svc.Index)
	store.Store.DeleteServerByName(svc.SubName)
	removeServer(svc.Name, svc.SubName, update)

	// 如果全部删光了， 那么scripts的name也要删除

	if store.Store.GetScriptLength(svc.Name) == 0 {
		store.Store.DeleteScriptByName(svc.Name)
	}
	atomic.AddInt64(&global.CanReload, -1)
}

func remove(svc *server.Server, update bool, wg *sync.WaitGroup) {
	// 如果是always 为 true，那么直接修改为false
	golog.Info("svc start removed")
	svc.Remove()
	golog.Info("svc removed")
	<-svc.StopSignal
	store.Store.DeleteScriptIndex(svc.Name, svc.Index)
	store.Store.DeleteServerByName(svc.SubName)
	removeServer(svc.Name, svc.SubName, update)

	// 如果全部删光了， 那么scripts的name也要删除

	if store.Store.GetScriptLength(svc.Name) == 0 {
		store.Store.DeleteScriptByName(svc.Name)
	}
	wg.Done()
}
