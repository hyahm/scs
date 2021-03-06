package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/alert"
	"github.com/hyahm/scs/global"
)

func Alert(w http.ResponseWriter, r *http.Request) {

	ra := &alert.RespAlert{}
	err := json.NewDecoder(r.Body).Decode(ra)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())))
		return
	}
	ra.SendAlert()
	w.Write([]byte(`{"code":200, "msg": "send alert message"}`))
}

func GetAlert(w http.ResponseWriter, r *http.Request) {
	w.Write(alert.GetDispatcher())
}

func Probe(w http.ResponseWriter, r *http.Request) {
	addr := strings.Split(r.RemoteAddr, ":")[0]
	needToken := true
	// 检查是否是被监控的
	for _, v := range global.Monitored {
		if v == addr {
			needToken = false
			break
		}
	}
	if !needToken {
		w.Write([]byte(`{"code": 200, "msg": "ok"}`))
		// w.WriteHeader(http.StatusOK)
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
		w.Write([]byte(`{"code": 200, "msg": "ok"}`))
		// w.WriteHeader(http.StatusOK)
		return
	}
	golog.Info(global.Monitored)
	w.Write([]byte(`{"code": 511, "msg": "StatusNetworkAuthenticationRequired"}`))
	// w.WriteHeader(http.StatusNetworkAuthenticationRequired)
}

// 报警相关配置
