package controller

import (
	"errors"
	"fmt"

	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/scs/pkg/config/scripts/subname"
)

func RestartServer(svc *server.Server) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	// 禁用 script 所在的所有server
	if _, ok := store.ss[svc.Name]; ok {
		restartServer(store.ss[svc.Name], svc)
	}

}

func restartServer(s *scripts.Script, svc *server.Server) {
	svc.Always = s.Always
	svc.Command = s.Command
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
	store.mu.RLock()
	defer store.mu.RUnlock()
	// 禁用 script 所在的所有server
	if _, ok := store.ss[s.Name]; !ok {
		return errors.New("not found " + s.Name)
	}
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		name := s.Name + fmt.Sprintf("_%d", i)

		if _, ok := store.servers[name]; ok {
			restartServer(s, store.servers[name])
		}
	}
	return nil
}

func RestartAllServer() {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, svc := range store.servers {
		RestartServer(svc)
	}
}

func RestartPermAllServer(token string) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, svc := range store.servers {
		if svc.Token == token {
			RestartServer(svc)
		}

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
