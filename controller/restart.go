package controller

import (
	"errors"
	"fmt"

	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/scs/pkg/config/scripts/subname"
)

func RestartServer(svc *server.Server) {
	mu.RLock()
	defer mu.RUnlock()
	// 禁用 script 所在的所有server
	if _, ok := ss[svc.Name]; ok {
		restartServer(ss[svc.Name], svc)
	}

}

func restartServer(s *scripts.Script, svc *server.Server) {
	svc.Always = s.Always
	svc.Command = s.Command
	svc.ContinuityInterval = s.ContinuityInterval
	svc.Cron = s.Cron
	svc.Dir = s.Dir
	svc.Disable = s.Disable
	svc.Replicate = s.Replicate
	svc.Status.RestartCount = 0
	svc.Port = s.Port
	if s.Token != "" {
		svc.Token = s.Token
	}
	svc.Update = s.Update
	svc.Version = s.Version

	svc.Restart()
}

// 异步重启
func RestartScript(s *scripts.Script) error {
	mu.RLock()
	defer mu.RUnlock()
	// 禁用 script 所在的所有server
	if _, ok := ss[s.Name]; !ok {
		return errors.New("not found " + s.Name)
	}
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		name := s.Name + fmt.Sprintf("_%d", i)

		if _, ok := servers[name]; ok {
			restartServer(s, servers[name])
		}
	}
	return nil
}

func RestartAllServer() {
	mu.RLock()
	defer mu.RUnlock()
	for _, svc := range servers {
		RestartServer(svc)
	}
}

func RestartPermAllServer(token string) {
	mu.RLock()
	defer mu.RUnlock()
	for _, svc := range servers {
		if svc.Token == token {
			RestartServer(svc)
		}

	}
}

// 返回成功还是失败
func UpdateAndRestartScript(s *scripts.Script) bool {
	mu.RLock()
	defer mu.RUnlock()
	if _, ok := ss[s.Name]; !ok {
		return false
	}
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		go servers[subname.String()].UpdateAndRestart()
	}
	return true
}

func UpdateAndRestartAllServer() {
	mu.RLock()
	defer mu.RUnlock()
	for _, s := range ss {
		go UpdateAndRestartScript(s)
	}
}

func UpdatePermAndRestartAllServer(token string) {
	mu.RLock()
	defer mu.RUnlock()
	for _, s := range ss {
		if s.Token == token {
			go UpdateAndRestartScript(s)
		}
	}
}
