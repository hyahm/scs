package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scs/alert"
	"scs/global"
	"scs/probe"
	"strings"
)

func Alert(w http.ResponseWriter, r *http.Request) {

	ra := &alert.RespAlert{}
	err := json.NewDecoder(r.Body).Decode(ra)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code":1, "msg": "%s"}`, err.Error())))
		return
	}
	ra.SendAlert()
	w.WriteHeader(http.StatusOK)
	return
}

func GetAlert(w http.ResponseWriter, r *http.Request) {
	w.Write(alert.GetDispatcher())
	return
}

func Probe(w http.ResponseWriter, r *http.Request) {
	addr := strings.Split(r.RemoteAddr, ":")[0]
	needToken := true
	// 检查是否是被监控的
	for _, v := range probe.GlobalProbe.Probe.Monitored {
		if v == addr {
			needToken = false
			break
		}
	}
	if !needToken {
		w.WriteHeader(http.StatusOK)
		return
	}
	// 检查是否可以被忽略token
	for _, v := range global.IgnoreToken {
		if v == addr {
			needToken = false
			break
		}
	}
	if !needToken {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNetworkAuthenticationRequired)
	return
}
