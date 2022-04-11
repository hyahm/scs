package handle

import (
	"net/http"

	"github.com/hyahm/scs/api/module"
	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"

	"github.com/hyahm/xmux"
)

func Restart(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	if global.CanReload != 0 {
		module.Write(w, r, pkg.WaitingConfigChanged())
		return
	}
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
	go controller.RestartServer(svc)
	module.Write(w, r, pkg.Waiting("restart"))
}

func RestartPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	if global.CanReload != 0 {
		module.Write(w, r, pkg.WaitingConfigChanged())
		return
	}
	// for _, pname := range strings.Split(names, ",") {
	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		module.Write(w, r, pkg.NotFoundScript())
		return
	}
	controller.RestartScript(script)
	// }
	module.Write(w, r, pkg.Waiting("restart"))
}

func RestartAll(w http.ResponseWriter, r *http.Request) {
	// 删除所有的脚本
	if global.CanReload != 0 {
		module.Write(w, r, pkg.WaitingConfigChanged())
		return
	}
	names := xmux.GetInstance(r).Get("scriptname")
	if names == nil {
		controller.RestartAllServer()
	} else {
		controller.RestartAllServerFromScripts(names.(map[string]struct{}))
	}

	module.Write(w, r, pkg.Waiting("restart"))
}
