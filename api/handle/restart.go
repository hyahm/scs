package handle

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/config/scripts/subname"

	"github.com/hyahm/xmux"
)

func Restart(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	svc, ok := controller.GetServerByNameAndSubname(pname, subname.Subname(name))
	if !ok {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, name)))
		return
	}
	go controller.RestartServer(svc)
	w.Write([]byte(`{"code": 200, "msg": "waiting restart"}`))
}

func RestartPname(w http.ResponseWriter, r *http.Request) {
	names := xmux.Var(r)["pname"]
	for _, pname := range strings.Split(names, ",") {
		s, ok := controller.GetScriptByPname(pname)
		if !ok {
			w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s"}`, pname)))
			return
		}
		controller.RestartScript(s)
	}
	w.Write([]byte(`{"code": 200, "msg": "waiting restart"}`))
}

func RestartAll(w http.ResponseWriter, r *http.Request) {
	// 删除所有的脚本
	controller.RestartAllServer()

	w.Write([]byte(`{"code": 200, "msg": "waiting restart"}`))
}
