package controller

import (
	"fmt"
	"runtime"

	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/probe"
)

func getStatus(svc *server.Server) pkg.ServiceStatus {
	status := pkg.ServiceStatus{
		PName:        svc.Name,
		Name:         svc.SubName,
		IsCron:       svc.IsCron,
		Command:      svc.Command,
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

func ScriptName(pname, subname string) []byte {
	store.mu.RLock()
	defer store.mu.RUnlock()
	status := &pkg.StatusList{
		Data:    make([]pkg.ServiceStatus, 0),
		Version: global.VERSION,
		Code:    200,
	}
	if _, ok := store.ss[pname]; !ok {
		return pkg.NotFoundScript()
	}
	if _, ok := store.servers[subname]; !ok {
		return pkg.NotFoundScript()
	}
	status.Data = append(status.Data, getStatus(store.servers[subname]))
	return status.Marshal()

}

func ScriptPname(pname string) []byte {
	store.mu.RLock()
	defer store.mu.RUnlock()
	statuss := &pkg.StatusList{
		Data:    make([]pkg.ServiceStatus, 0),
		Version: global.VERSION,
		Code:    200,
	}
	if _, ok := store.ss[pname]; !ok {
		return pkg.NotFoundScript()
	}
	for i := range store.serverIndex[pname] {
		subname := fmt.Sprintf("%s_%d", pname, i)
		statuss.Data = append(statuss.Data, getStatus(store.servers[subname]))
	}

	return statuss.Marshal()
}

// 获取所有服务的状态
func AllStatus() []byte {
	store.mu.RLock()
	defer store.mu.RUnlock()
	statuss := &pkg.StatusList{
		Data:    make([]pkg.ServiceStatus, 0),
		Version: global.VERSION,
		Code:    200,
	}
	for _, svc := range store.servers {
		if _, ok := store.ss[svc.Name]; ok {
			statuss.Data = append(statuss.Data, getStatus(svc))
		}
	}
	return statuss.Marshal()
}

// 获取所有服务的状态
func AllStatusFromScript(names map[string]struct{}) []byte {
	store.mu.RLock()
	defer store.mu.RUnlock()
	statuss := &pkg.StatusList{
		Data:    make([]pkg.ServiceStatus, 0),
		Version: global.VERSION,
		Code:    200,
	}
	for _, svc := range store.servers {
		if _, sok := names[svc.Name]; sok {
			if _, ok := store.ss[svc.Name]; ok {
				statuss.Data = append(statuss.Data, getStatus(svc))
			}
		}

	}
	return statuss.Marshal()
}
