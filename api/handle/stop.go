package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs/api/module"
	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"

	"github.com/hyahm/xmux"
)

func Stop(w http.ResponseWriter, r *http.Request) {
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
	go svc.Stop()
	module.Write(w, r, pkg.Waiting("stop"))
}

func StopPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		module.Write(w, r, pkg.NotFoundScript())
		return
	}
	err := controller.StopScript(script)
	if err != nil {
		module.Write(w, r, []byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())))
		return
	}

	module.Write(w, r, pkg.Waiting("stop"))
}

func StopAll(w http.ResponseWriter, r *http.Request) {
	scriptname := xmux.GetInstance(r).Get("scriptname")
	if scriptname == nil {
		controller.StopAllServer()
	} else {
		controller.StopScriptFromName(scriptname.(map[string]struct{}))
	}

	module.Write(w, r, pkg.Waiting("stop"))
}
