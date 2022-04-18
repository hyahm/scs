package probe

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/pkg/config/alert"
	"github.com/hyahm/scs/pkg/message"
)

var monitors Scan

type Scan map[string]*Monitor

type Monitor struct {
	AI       *alert.AlertInfo
	Interval time.Duration
}

func NewMonitor() Scan {
	monitors = make(map[string]*Monitor)

	for _, v := range healthDetector.Config.Monitor {
		monitors[v] = &Monitor{
			AI: &alert.AlertInfo{
				AM:                 &message.Message{},
				ContinuityInterval: healthDetector.Config.ContinuityInterval,
			},
			Interval: healthDetector.Config.Interval,
		}
	}
	return monitors
}

func (m Scan) Update() {
	temp := make(map[string]struct{})
	for k := range m {
		temp[k] = struct{}{}
	}
	for _, v := range healthDetector.Config.Monitor {
		if _, ok := m[v]; ok {
			m[v].Interval = healthDetector.Config.Interval
			m[v].AI.ContinuityInterval = healthDetector.Config.ContinuityInterval
			delete(temp, v)
		} else {
			m[v] = &Monitor{
				AI: &alert.AlertInfo{
					AM:                 &message.Message{},
					ContinuityInterval: healthDetector.Config.ContinuityInterval,
				},
				Interval: healthDetector.Config.Interval,
			}
		}
	}
	for k := range temp {
		delete(m, k)
	}
}

func requests(domain string) bool {
	req, err := http.NewRequest(http.MethodPost, domain+"/probe", nil)
	if err != nil {
		golog.Error(err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client(5 * time.Second).Do(req)
	if err != nil {
		golog.Error(err)
		return false
	}
	defer resp.Body.Close()
	golog.Info(resp.StatusCode)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		golog.Error(err)
		return false
	}
	golog.Info(string(b))
	return resp.StatusCode == 200
}

func client(timeout time.Duration) *http.Client {
	if timeout == 0 {
		timeout = 3 * time.Second
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout,
	}

}

func (m Scan) Check() {
	for server, mm := range m {
		//http cookie接口
		failed := requests(server)
		if failed {
			mm.AI.AM.HostName = server
			mm.AI.BreakDown(fmt.Sprintf("服务器或scs服务出现问题: %s", server))
			continue
		}
		mm.AI.Recover(fmt.Sprintf("服务器或scs服务恢复: %s", server))

	}

}
