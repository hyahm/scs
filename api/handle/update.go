package handle

import (
	"net/http"
	"strings"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/config/scripts/subname"

	"github.com/hyahm/xmux"
)

func Update(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	role := xmux.GetInstance(r).Get("role").(string)
	svc, ok := controller.GetServerByNameAndSubname(pname, subname.Subname(name))
	if !ok {
		w.Write(NotFoundScript(role))
		return
	}
	svc.UpdateAndRestart()

	w.Write(Waiting("update", role))
}

func UpdatePname(w http.ResponseWriter, r *http.Request) {
	names := xmux.Var(r)["pname"]
	role := xmux.GetInstance(r).Get("role").(string)
	for _, pname := range strings.Split(names, ",") {
		s, ok := controller.GetScriptByPname(pname)
		if !ok {
			w.Write(NotFoundScript(role))
			return
		}
		controller.UpdateAndRestartScript(s)
	}
	w.Write(Waiting("update", role))
}

func UpdateAll(w http.ResponseWriter, r *http.Request) {
	role := xmux.GetInstance(r).Get("role").(string)
	controller.UpdateAndRestartAllServer()
	w.Write(Waiting("update", role))
}
