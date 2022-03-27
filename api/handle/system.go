package handle

import (
	"encoding/json"
	"net/http"

	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
	"github.com/shirou/gopsutil/host"
)

func GetOS(w http.ResponseWriter, r *http.Request) {
	role := xmux.GetInstance(r).Get("role").(string)
	res := &pkg.Response{
		Role: role,
	}
	hi, err := host.Info()

	if err != nil {
		w.Write(res.ErrorE(err))
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
	w.Write(b)
}
