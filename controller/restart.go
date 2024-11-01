package controller

import (
	"fmt"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/scripts"
)

// 没有锁，只是为了外部访问
func RestartServer(svc *server.Server) {
	if svc.Disable {
		return
	}
	// 禁用 script 所在的所有server
	// 先修改值, 因为是restart， 所以端口在svc初始化的时候就固定了
	go restartServer(svc)

}

func restartServer(svc *server.Server) {
	// 先修改值
	svc.Restart()
	//已经停止了。
	<-svc.StopSignal
	// 更新server并启动
	script, ok := store.Store.GetScriptByName(svc.Name)
	if !ok {
		golog.Error(pkg.ErrBugMsg)
		return
	}
	svc.MakeServer(script)
	svc.Start()
}

// 重启第一步
func RestartScript(s *scripts.Script) error {
	// 禁用 script 所在的所有server
	if s.Disable {
		return nil
	}

	for index := range store.Store.GetScriptIndex(s.Name) {

		subname := fmt.Sprintf("%s_%d", s.Name, index)
		svc, ok := store.Store.GetServerByName(subname)
		if !ok {
			golog.Error(pkg.ErrBugMsg)
			continue
		}
		RestartServer(svc)
	}
	return nil
}

func RestartAllServer() {
	for _, svc := range store.Store.GetAllServer() {
		if svc.Disable {
			continue
		}
		RestartServer(svc)
	}
}

func RestartAllServerFromScripts(names map[string]struct{}) {
	for pname := range names {
		for index := range store.Store.GetScriptIndex(pname) {
			subname := fmt.Sprintf("%s_%d", pname, index)
			svc, ok := store.Store.GetServerByName(subname)
			if !ok {
				golog.Error(pkg.ErrBugMsg)
				continue
			}
			RestartServer(svc)
		}
	}
}
