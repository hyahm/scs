package probe

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"scs/alert"
	"time"

	"github.com/hyahm/golog"
)

type Monitor struct {
	Monitor  []string
	AI       *alert.AlertInfo
	Interval time.Duration
}

func NewMonitor(monitor []string, interval, continuityInterval time.Duration) *Monitor {
	return &Monitor{
		Monitor: monitor,
		AI: &alert.AlertInfo{
			AM:                 &alert.Message{},
			ContinuityInterval: continuityInterval,
		},
		Interval: interval,
	}
}

func (m *Monitor) Update(probe *Probe) {
	m.Monitor = probe.Monitor
	m.Interval = probe.Interval
	m.AI.ContinuityInterval = probe.ContinuityInterval
}

func (m *Monitor) Check() {
	for _, server := range m.Monitor {
		var failed bool
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		//http cookie接口
		cookieJar, _ := cookiejar.New(nil)
		c := &http.Client{
			Jar:       cookieJar,
			Transport: tr,
			Timeout:   time.Second * 5,
		}

		resp, err := c.Get(server + "/probe")
		if err != nil {
			golog.Error(err)
			failed = true
		} else {
			if resp.StatusCode != 200 {
				golog.Error(resp.StatusCode)
				failed = true
			}
		}

		if failed {
			m.AI.AM.HostName = server
			m.AI.BreakDown(fmt.Sprintf("服务器或scs服务出现问题: %s", server))
			continue
		}
		m.AI.Recover(fmt.Sprintf("服务器或scs服务恢复: %s", server))

	}

}
