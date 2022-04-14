package handle

import (
	"net/http"

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
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	svc, ok := store.Store.GetServerByName(name)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	go controller.UpdateAndRestart(svc)

}

func UpdatePname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	controller.UpdateAndRestartScript(script)
}

func UpdateAll(w http.ResponseWriter, r *http.Request) {
	names := xmux.GetInstance(r).Get("scriptname")
	if names == nil {
		controller.UpdateAllServer()
	} else {
		controller.UpdateAllServerFromScript(names.(map[string]struct{}))
	}

}
