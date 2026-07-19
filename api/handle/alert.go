package handle

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/xmux"
)

func Alert(w http.ResponseWriter, r *http.Request) {

	ra := &config.RespAlert{}
	err := json.NewDecoder(r.Body).Decode(ra)
	if err != nil {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
		return
	}
	ra.SendAlert()
}

func GetAlert(w http.ResponseWriter, r *http.Request) {
	xmux.GetInstance(r).Response.(*pkg.Response).Data = config.GetDispatcherList()
}

func Probe(w http.ResponseWriter, r *http.Request) {
	var addr string
	if config.Cfg.ProxyHeader == "" {
		addr = strings.Split(r.RemoteAddr, ":")[0]
	} else {
		addr = r.Header.Get(config.Cfg.ProxyHeader)
	}

	needToken := true
	// 检查是否是被监控的
	for _, v := range config.Cfg.Probe.Monitored {
		if strings.Contains(addr, v) {
			needToken = false
			break
		}
	}
	if needToken {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 203
		return
	}

}
