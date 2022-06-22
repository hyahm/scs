package controller

import (
	"fmt"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/scripts"
)

func KillScript(s *scripts.Script) {
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := fmt.Sprintf("%s_%d", s.Name, i)
		svc, ok := store.Store.GetServerByName(subname)
		if ok {
			svc.Kill()
		}

	}
}

func KillAndStartServer(param string, svc *server.Server) {
	go func() {
		svc.Kill()
		svc.Start(param)
	}()

}

func NeedStop(s *scripts.Script) bool {
	// 更新server
	// 判断值是否相等
	script, ok := store.Store.GetScriptByName(s.Name)
	if !ok {
		golog.Error(pkg.ErrBugMsg)
		return false
	}
	return !scripts.EqualScript(s, script)
}
