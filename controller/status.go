package controller

import (
	"fmt"
	"runtime"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/probe"
)

func getStatus(svc *server.Server) pkg.ServiceStatus {
	status := pkg.ServiceStatus{
		PName:        svc.Name,
		Name:         svc.SubName,
		IsCron:       svc.IsCron,
		Command:      svc.Status.Command,
		Version:      svc.Status.Version,
		CanNotStop:   svc.Status.CanNotStop,
		Path:         svc.Dir,
		Status:       svc.Status.Status,
		RestartCount: svc.Status.RestartCount,
		Pid:          svc.Status.Pid,
		Disable:      svc.Disable,
		OS:           runtime.GOOS,
		Start:        svc.Status.Start,
	}

	status.Cpu, status.Mem, _ = probe.GetProcessInfo(int32(status.Pid))
	return status
}

func ScriptName(pname, subname string) ([]pkg.ServiceStatus, error) {
	statuss := make([]pkg.ServiceStatus, 0)
	_, ok := store.Store.GetScriptByName(pname)
	if !ok {
		return nil, pkg.ErrNotFound
	}
	svc, ok := store.Store.GetServerByName(subname)
	if !ok {
		return nil, pkg.ErrNotFound
	}
	statuss = append(statuss, getStatus(svc))
	return statuss, nil

}

func ScriptPname(pname string) ([]pkg.ServiceStatus, error) {
	statuss := make([]pkg.ServiceStatus, 0)
	_, ok := store.Store.GetScriptByName(pname)
	if !ok {
		return nil, pkg.ErrNotFound
	}
	for i := range store.Store.GetScriptIndex(pname) {
		subname := fmt.Sprintf("%s_%d", pname, i)
		svc, ok := store.Store.GetServerByName(subname)
		if !ok {
			golog.Error(pkg.ErrBugMsg)
		}
		statuss = append(statuss, getStatus(svc))
	}

	return statuss, nil
}

// 获取所有服务的状态
func AllStatus() []pkg.ServiceStatus {
	statuss := make([]pkg.ServiceStatus, 0)
	for _, svc := range store.Store.GetAllServer() {
		statuss = append(statuss, getStatus(svc))
	}
	return statuss
}

// 获取所有服务的状态
func AllStatusFromScript(names map[string]struct{}) []pkg.ServiceStatus {
	statuss := make([]pkg.ServiceStatus, 0)
	for _, svc := range store.Store.GetAllServer() {
		if _, sok := names[svc.Name]; sok {
			statuss = append(statuss, getStatus(svc))
		}
	}
	return statuss
}
