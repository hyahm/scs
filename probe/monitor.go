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

// var monitors Scan

type Scan map[string]*Monitor

type Monitor struct {
	AI       *alert.AlertInfo
	Interval time.Duration
}

func NewMonitor(monitor []string, interval, continuityInterval time.Duration) Scan {
	monitors := make(map[string]*Monitor)
	for _, v := range monitor {
		monitors[v] = &Monitor{
			AI: &alert.AlertInfo{
				AM:                 &alert.Message{},
				ContinuityInterval: continuityInterval,
			},
			Interval: interval,
		}
	}
	return monitors
}

func (m Scan) Update(probe *Probe) {
	for _, v := range probe.Monitor {
		if _, ok := m[v]; ok {
			m[v].Interval = probe.Interval
			m[v].AI.ContinuityInterval = probe.ContinuityInterval
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
}

func (m Scan) Check() {
	for server, mm := range m {
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
			mm.AI.AM.HostName = server
			mm.AI.BreakDown(fmt.Sprintf("服务器或scs服务出现问题: %s", server))
			continue
		}
		mm.AI.Recover(fmt.Sprintf("服务器或scs服务恢复: %s", server))

	}

}
