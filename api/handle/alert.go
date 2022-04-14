package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hyahm/scs/api/module"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/alert"
)

func Alert(w http.ResponseWriter, r *http.Request) {
	if global.CanReload != 0 {
		module.Write(w, r, pkg.WaitingConfigChanged())

		return
	}
	res := pkg.Response{}
	ra := &alert.RespAlert{}
	err := json.NewDecoder(r.Body).Decode(ra)
	if err != nil {
		module.Write(w, r, res.ErrorE(err))
		return
	}
	ra.SendAlert()
	module.Write(w, r, res.Sucess("send alert message"))
}

func GetAlert(w http.ResponseWriter, r *http.Request) {
	res := pkg.Response{
		Data: alert.GetDispatcher(),
	}
	module.Write(w, r, res.Sucess(""))
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
	if !needToken {
		module.Write(w, r, []byte(`{"code": 200, "msg": "ok"}`))
		// w.WriteHeader(http.StatusOK)
		return
	}
	// 检查是否可以被忽略token
	for _, v := range global.GetIgnoreToken() {
		if v == addr {
			needToken = false
			break
		}
	}
	if !needToken {
		module.Write(w, r, []byte(`{"code": 200, "msg": "ok"}`))
		// w.WriteHeader(http.StatusOK)
		return
	}
	module.Write(w, r, []byte(fmt.Sprintf(`{"code": 500, "msg": "StatusNetworkAuthenticationRequired"}, "ips": "%s"`, addr)))
	// w.WriteHeader(http.StatusNetworkAuthenticationRequired)
}

// 报警相关配置
