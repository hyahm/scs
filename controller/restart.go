package controller

import (
	"errors"
	"fmt"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/scripts"
)

// 没有锁，只是为了外部访问
func RestartServer(svc *server.Server) {
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
	svc.MakeServer(script, svc.Port)
	svc.Start()
}

// 重启第一步
func RestartScript(s *scripts.Script) error {
	// 禁用 script 所在的所有server
	script, ok := store.Store.GetScriptByName(s.Name)
	if !ok {
		return errors.New("not found script: " + s.Name)
	}
	replicate := script.Replicate
	if replicate == 0 {
		replicate = 1
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
		svc.Restart()
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

// 返回成功还是失败
func UpdateAndRestartScript(s *scripts.Script) bool {
	_, ok := store.Store.GetScriptByName(s.Name)
	if !ok {
		return false
	}
	return updateAndRestartScript(s)
}

func updateAndRestartScript(s *scripts.Script) bool {
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := fmt.Sprintf("%s_%d", s.Name, i)
		svc, ok := store.Store.GetServerByName(subname)
		if ok {
			go svc.UpdateAndRestart()
		}
	}
	return true
}

func UpdateAllServer() {
	for _, svc := range store.Store.GetAllServer() {
		go svc.UpdateAndRestart()
	}
}

func UpdateAllServerFromScript(names map[string]struct{}) {
	for pname := range names {
		script, ok := store.Store.GetScriptByName(pname)
		if ok {
			go updateAndRestartScript(script)
		}
	}
}
