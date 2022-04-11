package handle

import (
	"net/http"

	"github.com/hyahm/scs/api/module"
	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func Kill(w http.ResponseWriter, r *http.Request) {
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
	go svc.Kill()
	module.Write(w, r, []byte(`{"code": 200, "msg": "already killed"}`))
}

func KillPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		module.Write(w, r, pkg.NotFoundScript())
		return
	}

	controller.KillScript(script)
	module.Write(w, r, []byte(`{"code": 200, "msg": "already killed"}`))
}
