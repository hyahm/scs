package handle

import (
	"net/http"

	"github.com/hyahm/scs/api/module"
	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"

	"github.com/hyahm/xmux"
)

func Update(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	_, ok := store.Store.GetScriptByName(pname)
	if !ok {
		module.Write(w, r, pkg.NotFoundScript())
		return
	}
	svc, ok := store.Store.GetServerByName(name)
	if !ok {
		module.Write(w, r, pkg.NotFoundScript())
		return
	}
	go controller.UpdateAndRestart(svc)

	module.Write(w, r, pkg.Waiting("update"))
}

func UpdatePname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		module.Write(w, r, pkg.NotFoundScript())
		return
	}
	controller.UpdateAndRestartScript(script)
	module.Write(w, r, pkg.Waiting("update"))
}

func UpdateAll(w http.ResponseWriter, r *http.Request) {
	names := xmux.GetInstance(r).Get("scriptname")
	if names == nil {
		controller.UpdateAllServer()
	} else {
		controller.UpdateAllServerFromScript(names.(map[string]struct{}))
	}

	module.Write(w, r, pkg.Waiting("update"))
}
