package controller

import (
	"encoding/json"
	"runtime"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/config/probe"
	"github.com/hyahm/scs/status"
)

type StatusList struct {
	Data    []*status.ServiceStatus `json:"data"`
	Code    int                     `json:"code"`
	Msg     string                  `json:"msg"`
	Version string                  `json:"version"`
	Role    string                  `json:"role"`
}

func (sl *StatusList) Marshal() []byte {
	sl.Code = 200
	b, err := json.Marshal(sl)
	if err != nil {
		golog.Error(err)
	}
	return b
}

func getStatus(name, subname string) *status.ServiceStatus {
	// subname := svc.SubName
	status := &status.ServiceStatus{
		PName:        name,
		Name:         subname,
		IsCron:       servers[subname].IsLoop,
		Command:      servers[subname].Status.Command,
		Always:       ss[name].Always,
		Version:      servers[subname].Status.Version,
		CanNotStop:   servers[subname].Status.CanNotStop,
		Path:         servers[subname].Status.Path,
		Status:       servers[subname].Status.Status,
		RestartCount: servers[subname].Status.RestartCount,
		// Up:           servers[subname].Status.Up,
		Disable:   ss[name].Disable,
		OS:        runtime.GOOS,
		Start:     servers[subname].Status.Start,
		SCSVerion: global.VERSION,
	}
	if servers[subname].Cmd != nil && servers[subname].Cmd.Process != nil {
		status.Pid = servers[subname].Cmd.Process.Pid
		status.Cpu, status.Mem, _ = probe.GetProcessInfo(int32(servers[subname].Cmd.Process.Pid))

	}
	return status
}
