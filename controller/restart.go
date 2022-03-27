package controller

import (
	"errors"
	"fmt"

	"github.com/hyahm/scs/internal/config/scripts"
	"github.com/hyahm/scs/internal/config/scripts/subname"
	"github.com/hyahm/scs/internal/server"
)

func RestartServer(svc *server.Server) {
	mu.RLock()
	defer mu.RUnlock()
	// 禁用 script 所在的所有server
	svc.Status.RestartCount = 0
	makeReplicateServerAndStart(ss[svc.Name], svc.Replicate)
}

// 异步重启
func RestartScript(s *scripts.Script) error {
	mu.RLock()
	defer mu.RUnlock()
	// 禁用 script 所在的所有server
	if _, ok := ss[s.Name]; !ok {
		return errors.New("")
	}
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		if _, ok := servers[s.Name+fmt.Sprintf("_%d", i)]; !ok {
			servers[s.Name+fmt.Sprintf("_%d", i)] = &server.Server{Index: i}
		} else {
			servers[s.Name+fmt.Sprintf("_%d", i)].Status.RestartCount = 0
		}

		serverIndex[s.Name] = make(map[int]struct{})
		makeReplicateServerAndStart(s, replicate)
	}
	return nil
}

func RestartAllServer() {
	mu.RLock()
	defer mu.RUnlock()
	for _, svc := range servers {
		go svc.Restart()
	}
}

func RestartPermAllServer(token string) {
	mu.RLock()
	defer mu.RUnlock()
	for _, svc := range servers {
		if svc.Token == token {
			go svc.Restart()
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
