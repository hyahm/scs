package controller

import "github.com/hyahm/scs/internal/store"

func WaitKillAllServer() {
	for _, svc := range store.Store.GetAllServer() {
		svc.Kill()
	}
}
