package handle

import (
	"net/http"

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
		w.Write(pkg.WaitingConfigChanged())
		return
	}
	_, ok := store.Store.GetScriptByName(pname)
	if !ok {
		w.Write(pkg.NotFoundScript())
		return
	}
	svc, ok := store.Store.GetServerByName(name)
	if !ok {
		w.Write(pkg.NotFoundScript())
		return
	}
	go controller.RestartServer(svc)
	w.Write(pkg.Waiting("restart"))
}

func RestartPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	if global.CanReload != 0 {
		w.Write(pkg.WaitingConfigChanged())
		return
	}
	// for _, pname := range strings.Split(names, ",") {
	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		w.Write(pkg.NotFoundScript())
		return
	}
	controller.RestartScript(script)
	// }
	w.Write(pkg.Waiting("restart"))
}

func RestartAll(w http.ResponseWriter, r *http.Request) {
	// 删除所有的脚本
	if global.CanReload != 0 {
		w.Write(pkg.WaitingConfigChanged())
		return
	}
	names := xmux.GetInstance(r).Get("scriptname")
	if names == nil {
		controller.RestartAllServer()
	} else {
		controller.RestartAllServerFromScripts(names.(map[string]struct{}))
	}

	w.Write(pkg.Waiting("restart"))
}
