package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/store"

	"github.com/hyahm/xmux"
)

func Stop(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	_, ok := store.Store.GetScriptByName(pname)
	if !ok {
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 404)
		return
	}
	svc, ok := store.Store.GetServerByName(name)
	if !ok {
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 404)
		return
	}
	go svc.Stop()
}

func StopPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 404)
		return
	}
	err := controller.StopScript(script)
	if err != nil {
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 500)
		return
	}
}

func StopAll(w http.ResponseWriter, r *http.Request) {
	scriptname := xmux.GetInstance(r).Get("scriptname")
	if scriptname == nil {
		controller.StopAllServer()
	} else {
		controller.StopScriptFromName(scriptname.(map[string]struct{}))
	}

}
