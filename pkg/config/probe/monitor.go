package probe

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	golog.Info(monitors)
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

// func NewClient(timeout ...time.Duration) *SCSClient {
// 	var rto time.Duration
// 	if len(timeout) > 0 {
// 		rto = timeout[0]
// 	}

// 	return &SCSClient{
// 		Domain:  "https://127.0.0.1:11111",
// 		Token:   os.Getenv("TOKEN"),
// 		Pname:   os.Getenv("PNAME"),
// 		Name:    os.Getenv("NAME"),
// 		Timeout: rto,
// 	}
// }

func requests(domain string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, domain+"/probe", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client(5 * time.Second).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
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
		var failed bool
		//http cookie接口
		resp, err := requests(server)
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
