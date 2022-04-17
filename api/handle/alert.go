package handle

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/alert"
	"github.com/hyahm/xmux"
)

func Alert(w http.ResponseWriter, r *http.Request) {
	if global.CanReload != 0 {
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 201)
		return
	}
	ra := &alert.RespAlert{}
	err := json.NewDecoder(r.Body).Decode(ra)
	if err != nil {
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 500)
		return
	}
	ra.SendAlert()
}

func GetAlert(w http.ResponseWriter, r *http.Request) {
	xmux.GetInstance(r).Response.(*pkg.Response).Data = alert.GetDispatcher()
}

func Probe(w http.ResponseWriter, r *http.Request) {
	var addr string
	if global.ProxyHeader == "" {
		addr = strings.Split(r.RemoteAddr, ":")[0]
	} else {
		addr = r.Header.Get(global.ProxyHeader)
	}

	needToken := true
	// 检查是否是被监控的
	for _, v := range global.GetMonitored() {
		if strings.Contains(addr, v) {
			needToken = false
			break
		}
	}
	if needToken {
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 203)
		return
	}

}
