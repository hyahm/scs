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

func RestartServer(svc *server.Server, script *scripts.Script) {
	// 禁用 script 所在的所有server
	// restartServer(store.ss[svc.Name], svc)

	// 先修改值
	svc.MakeServer(script, svc.Port)
	golog.Info(svc.Command)
	atomic.AddInt64(&global.CanReload, 1)
	go restartServer(svc)

}

func restartServer(svc *server.Server) {

	// restartServer(store.ss[svc.Name], svc)
	if _, ok := store.ss[svc.Name]; !ok {
		golog.Error("not found script: ", svc.Name)
		return
	}
	// 先修改值
	svc.MakeServer(store.ss[svc.Name], store.servers[svc.SubName].Port)
	svc.Restart()
	atomic.AddInt64(&global.CanReload, -1)
}

// 异步重启
func RestartScript(s *scripts.Script) error {
	store.mu.RLock()
	defer store.mu.RUnlock()
	// 禁用 script 所在的所有server

	for i := range store.serverIndex[s.Name] {
		subname := fmt.Sprintf("%s_%d", s.Name, i)
		if _, ok := store.servers[subname]; ok {
			atomic.AddInt64(&global.CanReload, 1)
			go restartServer(store.servers[subname])
		}
	}
	return nil
}

func RestartAllServer(token string) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, svc := range store.servers {
		if token != "" && token != svc.Token {
			continue
		}
		if _, ok := store.ss[svc.Name]; !ok {
			golog.Error("not found script: ", svc.Name)
			return
		}
		RestartServer(svc, store.ss[svc.Name])
	}
}

// 返回成功还是失败
func UpdateAndRestartScript(s *scripts.Script) bool {
	store.mu.RLock()
	defer store.mu.RUnlock()
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

func UpdateAndRestartAllServer() {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, s := range store.ss {
		go UpdateAndRestartScript(s)
	}
}

func UpdatePermAndRestartAllServer(token string) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, s := range store.ss {
		if s.Token == token {
			go UpdateAndRestartScript(s)
		}
	}
}
