package handle

import (
	"encoding/json"
	"net/http"

	"github.com/hyahm/scs/api/module"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
	"github.com/shirou/gopsutil/host"
)

func GetOS(w http.ResponseWriter, r *http.Request) {
	res := &pkg.Response{}
	hi, err := host.Info()

	if err != nil {
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 500)
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
	res.Data = info
	b, _ := json.Marshal(res)
	module.Write(w, r, b)
}
