package handle

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/config/scripts/subname"

	"github.com/hyahm/xmux"
)

func Stop(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]

	svc, ok := controller.GetServerByNameAndSubname(pname, subname.Subname(name))

	if !ok {
		w.Write([]byte(`{"code": 404, "msg": "not found this script"}`))
		return
	}
	go svc.Stop()
	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
}

func StopPname(w http.ResponseWriter, r *http.Request) {
	names := xmux.Var(r)["pname"]
	for _, pname := range strings.Split(names, ",") {
		s, ok := controller.GetScriptByPname(pname)
		if !ok {
			w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s}`, pname)))
			return
		}
		err := controller.StopScript(s)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())))
			return
		}
	}

	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
}

func StopAll(w http.ResponseWriter, r *http.Request) {

	controller.StopAllServer()
	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
}
