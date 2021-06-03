package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs/server"
	"github.com/hyahm/xmux"
)

var reloadKey bool

func Disable(w http.ResponseWriter, r *http.Request) {
	if reloadKey {
		w.Write([]byte(`{"code": 201, "msg": "config file is reloading, waiting completed first"}`))
		return
	}
	reloadKey = true
	defer func() {
		reloadKey = false
	}()
	pname := xmux.Var(r)["pname"]

	s, err := server.GetScriptByPname(pname)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s}`, pname)))
		return
	}
	// 上面已经判断过是否存在了， 这里就忽略
	server.DisableScript(s)
	err = server.UpdateScriptToConfigFile(s)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s}`, err.Error())))
		return
	}
	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))

}

func Enable(w http.ResponseWriter, r *http.Request) {
	if reloadKey {
		w.Write([]byte(`{"code": 201, "msg": "config file is reloading, waiting completed first"}`))
		return
	}
	reloadKey = true
	defer func() {
		reloadKey = false
	}()
	pname := xmux.Var(r)["pname"]
	s, err := server.GetScriptByPname(pname)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s}`, pname)))
		return
	}
	// 上面已经判断过是否存在了， 这里就忽略
	server.EnableScript(s)
	err = server.UpdateScriptToConfigFile(s)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s}`, err.Error())))
		return
	}

	w.Write([]byte(`{"code": 200, "msg": "waiting start"}`))
}
