package handle

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shirou/gopsutil/host"
)

func GetOS(w http.ResponseWriter, r *http.Request) {
	hi, err := host.Info()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 1, "msg": "%v"}`, err)))
		return
	}
	type Info struct {
		Hostname        string `json:"hostname"`
		Uptime          uint64 `json:"uptime"`
		OS              string `json:"os"`
		Platform        string `json:"platform"`
		PlatformFamily  string `json:"platformFamily"`
		PlatformVersion string `json:"platformVersion"`
	}
	info := &Info{
		Hostname:        hi.Hostname,
		Uptime:          hi.Uptime,
		OS:              hi.OS,
		Platform:        hi.Platform,
		PlatformFamily:  hi.PlatformFamily,
		PlatformVersion: hi.PlatformVersion,
	}
	type Resp struct {
		Code int   `json:"code"`
		Data *Info `json:"data"`
	}
	res := &Resp{
		Data: info,
	}

	b, _ := json.Marshal(res)
	w.Write(b)
}
