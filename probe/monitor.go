package probe

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyahm/scs/alert"
	"github.com/hyahm/scs/client"
	"github.com/hyahm/scs/internal"

	"github.com/hyahm/golog"
)

// var monitors Scan

type Scan map[string]*Monitor

type Monitor struct {
	AI       *alert.AlertInfo
	Interval time.Duration
}

func NewMonitor(monitor []string, interval, continuityInterval time.Duration) Scan {
	monitors := make(map[string]*Monitor)
	golog.Info(monitors)
	for _, v := range monitor {
		monitors[v] = &Monitor{
			AI: &alert.AlertInfo{
				AM:                 &alert.Message{},
				ContinuityInterval: continuityInterval,
			},
			Interval: interval,
		}
	}
	golog.Info(monitors)
	return monitors
}

func (m Scan) Update(probe *Probe) {
	temp := make(map[string]struct{})
	for k, _ := range m {
		temp[k] = struct{}{}
	}
	for _, v := range probe.Monitor {
		if _, ok := m[v]; ok {
			m[v].Interval = probe.Interval
			m[v].AI.ContinuityInterval = probe.ContinuityInterval
			delete(temp, v)
		} else {
			m[v] = &Monitor{
				AI: &alert.AlertInfo{
					AM:                 &alert.Message{},
					ContinuityInterval: probe.ContinuityInterval,
				},
				Interval: probe.Interval,
			}
		}
	}
	for k, _ := range temp {
		delete(m, k)
	}
}

func (m Scan) Check() {
	c := client.NewClient()
	golog.Info(m)
	for server, mm := range m {
		c.Domain = server
		var failed bool

		//http cookie接口

		resp, err := c.Probe()
		if err != nil {
			golog.Error(err)
			failed = true
		} else {
			rest := &internal.Resp{}
			golog.Info(string(resp))
			err := json.Unmarshal(resp, rest)
			if err != nil {
				golog.Error(err)
				break
			}
			if rest.Code != 200 {
				golog.Error(rest.Msg)
				failed = true
			}
		}

		if failed {
			mm.AI.AM.HostName = server
			mm.AI.BreakDown(fmt.Sprintf("服务器或scs服务出现问题: %s", server))
			continue
		}
		mm.AI.Recover(fmt.Sprintf("服务器或scs服务恢复: %s", server))

	}

}
