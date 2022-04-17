package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/store"

	"github.com/hyahm/xmux"
)

func Restart(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	if global.CanReload != 0 {
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 201)
		return
	}
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
	go controller.RestartServer(svc)
}

func RestartPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	if global.CanReload != 0 {
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 201)
		return
	}
	// for _, pname := range strings.Split(names, ",") {
	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 404)
		return
	}
	controller.RestartScript(script)
	// }
}

func RestartAll(w http.ResponseWriter, r *http.Request) {
	// 删除所有的脚本
	if global.CanReload != 0 {
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 201)
		return
	}
	names := xmux.GetInstance(r).Get("scriptname")
	if names == nil {
		controller.RestartAllServer()
	} else {
		controller.RestartAllServerFromScripts(names.(map[string]struct{}))
	}

}
