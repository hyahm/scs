package handle

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/config/scripts/subname"
	"github.com/hyahm/xmux"
)

func Kill(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	svc, ok := controller.GetServerByNameAndSubname(pname, subname.Subname(name))
	if !ok {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, name)))
		return
	}
	go svc.Kill()
	w.Write([]byte(`{"code": 200, "msg": "already killed"}`))
}

func KillPname(w http.ResponseWriter, r *http.Request) {
	names := xmux.Var(r)["pname"]
	for _, pname := range strings.Split(names, ",") {
		s, ok := controller.GetScriptByPname(pname)
		if !ok {
			w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, pname)))
			return
		}
		controller.KillScript(s)
	}
	w.Write([]byte(`{"code": 200, "msg": "already killed"}`))
}
