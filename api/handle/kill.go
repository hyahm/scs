package handle

import (
	"net/http"
	"strings"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/scripts/subname"
	"github.com/hyahm/xmux"
)

func Kill(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	svc, _, ok := controller.GetServerByNameAndSubname(pname, subname.Subname(name))
	if !ok {
		w.Write(pkg.NotFoundScript())
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
			w.Write(pkg.NotFoundScript())
			return
		}
		controller.KillScript(s)
	}
	w.Write([]byte(`{"code": 200, "msg": "already killed"}`))
}
