package controller

import (
	"runtime"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/pkg/config/probe"
	"github.com/hyahm/scs/status"
)

func getStatus(name, subname string) *status.ServiceStatus {
	// subname := svc.SubName
	golog.Info(servers[subname].Status.Pid)
	status := &status.ServiceStatus{
		PName:        name,
		Name:         subname,
		IsCron:       servers[subname].IsCron,
		Command:      servers[subname].Status.Command,
		Always:       ss[name].Always,
		Version:      servers[subname].Status.Version,
		CanNotStop:   servers[subname].Status.CanNotStop,
		Path:         servers[subname].Status.Path,
		Status:       servers[subname].Status.Status,
		RestartCount: servers[subname].Status.RestartCount,
		Pid:          servers[subname].Status.Pid,
		Disable:      ss[name].Disable,
		OS:           runtime.GOOS,
		Start:        servers[subname].Status.Start,
		SCSVerion:    global.VERSION,
	}

	status.Cpu, status.Mem, _ = probe.GetProcessInfo(int32(servers[subname].Cmd.Process.Pid))
	return status
}
