package controller

import (
	"fmt"

	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg/config"
)

func WaitKillAllServer() {
	for _, svc := range store.GetStore().GetAllServer() {
		svc.Kill()
	}
}

func KillScript(s config.Script) {
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := fmt.Sprintf("%s_%d", s.Name, i)
		svc, ok := store.GetStore().GetServerByName(subname)
		if ok {
			svc.Kill()
		}

	}
}

func KillAndStartServer(param string, svc *server.Server) {
	go func() {
		svc.Kill()
		svc.Start()
	}()

}
