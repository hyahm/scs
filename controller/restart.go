package controller

import (
	"fmt"
	"sync/atomic"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/scs/pkg/config/scripts/subname"
)

// 没有锁，只是为了外部访问
func RestartServer(svc *server.Server, script *scripts.Script) {
	// 禁用 script 所在的所有server
	// 先修改值, 因为是restart， 所以端口在svc初始化的时候就固定了
	atomic.AddInt64(&global.CanReload, 1)
	go restartServer(svc, script)

}

func restartServer(svc *server.Server, script *scripts.Script) {
	// 先修改值
	svc.Restart()
	//已经停止了。
	<-svc.StopSignal
	// 更新server并启动
	svc.MakeServer(script, svc.Port)
	svc.Start()
	atomic.AddInt64(&global.CanReload, -1)
}

// 重启第一步
func RestartScript(s *scripts.Script) error {
	store.mu.RLock()
	defer store.mu.RUnlock()
	// 禁用 script 所在的所有server
	replicate := store.ss[s.Name].Replicate
	if replicate == 0 {
		replicate = 1
	}

	for i := 0; i < replicate; i++ {
		subname := fmt.Sprintf("%s_%d", s.Name, i)
		if _, ok := store.servers[subname]; ok {
			// 这里主要是发送restart信号
			// 只是停止 + 发送restart信号
			golog.Info("restart ", subname)
			RestartServer(store.servers[subname], s)
			// 等待停止信号，我们就要重新makeserver

		}
	}
	// for i := range store.serverIndex[s.Name] {
	// 	subname := fmt.Sprintf("%s_%d", s.Name, i)
	// 	if _, ok := store.servers[subname]; ok {
	// 		// 这里主要是发送restart信号
	// 		// 只是停止 + 发送restart信号
	// 		RestartServer(store.servers[subname], s)
	// 		// 等待停止信号，我们就要重新makeserver

	// 	}
	// }
	return nil
}

func RestartAllServer() {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, svc := range store.servers {
		if _, ok := store.ss[svc.Name]; !ok {
			golog.Error("not found script: ", svc.Name)
			return
		}
		RestartServer(svc, store.ss[svc.Name])
	}
}

func RestartAllServerFromScripts(names map[string]struct{}) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, svc := range store.servers {
		if _, sok := names[svc.Name]; sok {
			if _, ok := store.ss[svc.Name]; ok {
				RestartServer(svc, store.ss[svc.Name])
			}
		}

	}
}

func UpdateAndRestartScript(s *scripts.Script) bool {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return updateAndRestartScript(s)
}

// 返回成功还是失败
func updateAndRestartScript(s *scripts.Script) bool {
	if _, ok := store.ss[s.Name]; !ok {
		return false
	}
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		go store.servers[subname.String()].UpdateAndRestart()
	}
	return true
}

func UpdateAllServer() {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, s := range store.ss {
		go updateAndRestartScript(s)
	}
}

func UpdateAllServerFromScript(names map[string]struct{}) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, s := range store.ss {
		if _, ok := names[s.Name]; ok {
			go updateAndRestartScript(s)
		}

	}
}
