package controller

import (
	"runtime"

	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/pkg/config/probe"
	"github.com/hyahm/scs/status"
)

func getStatus(name, subname string) status.ServiceStatus {
	// subname := svc.SubName
	status := status.ServiceStatus{
		PName:        name,
		Name:         subname,
		IsCron:       store.servers[subname].IsCron,
		Command:      store.servers[subname].Status.Command,
		Always:       store.ss[name].Always,
		Version:      store.servers[subname].Status.Version,
		CanNotStop:   store.servers[subname].Status.CanNotStop,
		Path:         store.servers[subname].Status.Path,
		Status:       store.servers[subname].Status.Status,
		RestartCount: store.servers[subname].Status.RestartCount,
		Pid:          store.servers[subname].Status.Pid,
		Disable:      store.ss[name].Disable,
		OS:           runtime.GOOS,
		Start:        store.servers[subname].Status.Start,
		SCSVerion:    global.VERSION,
	}

	status.Cpu, status.Mem, _ = probe.GetProcessInfo(int32(store.servers[subname].Cmd.Process.Pid))
	return status
}
