package handle

import (
	"net/http"
	"strings"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/scripts/subname"

	"github.com/hyahm/xmux"
)

func Update(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	svc, _, ok := controller.GetServerByNameAndSubname(pname, subname.Subname(name))
	if !ok {
		w.Write(pkg.NotFoundScript())
		return
	}
	go controller.UpdateAndRestart(svc)

	w.Write(pkg.Waiting("update"))
}

func UpdatePname(w http.ResponseWriter, r *http.Request) {
	names := xmux.Var(r)["pname"]
	for _, pname := range strings.Split(names, ",") {
		s, ok := controller.GetScriptByPname(pname)
		if !ok {
			w.Write(pkg.NotFoundScript())
			return
		}
		controller.UpdateAndRestartScript(s)
	}
	w.Write(pkg.Waiting("update"))
}

func UpdateAll(w http.ResponseWriter, r *http.Request) {
	names := xmux.GetInstance(r).Get("scriptname")
	if names == nil {
		controller.UpdateAllServer()
	} else {
		controller.UpdateAllServerFromScript(names.(map[string]struct{}))
	}

	w.Write(pkg.Waiting("update"))
}
