package script

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyahm/golog"
)

var monitors Scan

type Scan map[string]*Monitor

type Monitor struct {
	AI       *AlertInfo
	Interval time.Duration
}

func NewMonitor() Scan {
	monitors = make(map[string]*Monitor)

	for _, v := range healthDetector.Config.Monitor {
		monitors[v] = &Monitor{
			AI: &AlertInfo{
				AM:                 &Message{},
				ContinuityInterval: healthDetector.Config.ContinuityInterval,
			},
			Interval: healthDetector.Config.Interval,
		}
	}
	golog.Info(monitors)
	return monitors
}

func (m Scan) Update() {
	temp := make(map[string]struct{})
	for k, _ := range m {
		temp[k] = struct{}{}
	}
	for _, v := range healthDetector.Config.Monitor {
		if _, ok := m[v]; ok {
			m[v].Interval = healthDetector.Config.Interval
			m[v].AI.ContinuityInterval = healthDetector.Config.ContinuityInterval
			delete(temp, v)
		} else {
			m[v] = &Monitor{
				AI: &AlertInfo{
					AM:                 &Message{},
					ContinuityInterval: healthDetector.Config.ContinuityInterval,
				},
				Interval: healthDetector.Config.Interval,
			}
		}
	}
	for k, _ := range temp {
		delete(m, k)
	}
}

func (m Scan) Check() {
	c := NewClient()
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
			rest := &struct {
				Code int    `json:"code"`
				Msg  string `json:"msg"`
			}{}
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
