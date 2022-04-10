package controller

import (
	"fmt"
	"runtime"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
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
	status := &pkg.StatusList{
		Data:    make([]pkg.ServiceStatus, 0),
		Version: global.VERSION,
		Code:    200,
	}
	_, ok := store.Store.GetScriptByName(pname)
	if !ok {
		return pkg.NotFoundScript()
	}
	svc, ok := store.Store.GetServerByName(subname)
	if !ok {
		return pkg.NotFoundScript()
	}
	status.Data = append(status.Data, getStatus(svc))
	return status.Marshal()

}

func ScriptPname(pname string) []byte {
	statuss := &pkg.StatusList{
		Data:    make([]pkg.ServiceStatus, 0),
		Version: global.VERSION,
		Code:    200,
	}
	_, ok := store.Store.GetScriptByName(pname)
	if !ok {
		return pkg.NotFoundScript()
	}
	for i := range store.Store.GetScriptIndex(pname) {
		subname := fmt.Sprintf("%s_%d", pname, i)
		svc, ok := store.Store.GetServerByName(subname)
		if !ok {
			golog.Error(pkg.ErrBugMsg)
		}
		statuss.Data = append(statuss.Data, getStatus(svc))
	}

	return statuss.Marshal()
}

// 获取所有服务的状态
func AllStatus() []byte {
	statuss := &pkg.StatusList{
		Data:    make([]pkg.ServiceStatus, 0),
		Version: global.VERSION,
		Code:    200,
	}

	for _, svc := range store.Store.GetAllServer() {
		statuss.Data = append(statuss.Data, getStatus(svc))
	}
	return statuss.Marshal()
}

// 获取所有服务的状态
func AllStatusFromScript(names map[string]struct{}) []byte {
	statuss := &pkg.StatusList{
		Data:    make([]pkg.ServiceStatus, 0),
		Version: global.VERSION,
		Code:    200,
	}
	for _, svc := range store.Store.GetAllServer() {
		if _, sok := names[svc.Name]; sok {
			statuss.Data = append(statuss.Data, getStatus(svc))
		}

	}
	return statuss.Marshal()
}
