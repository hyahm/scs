package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func GetAlarms(w http.ResponseWriter, r *http.Request) {
	xmux.GetInstance(r).Response.(*pkg.Response).Data = controller.GetAterts()
}

func GetServers(w http.ResponseWriter, r *http.Request) {
	namesInterface := xmux.GetInstance(r).Get("scriptname")
	if namesInterface == nil {
		xmux.GetInstance(r).Response.(*pkg.Response).Data = store.Store.GetAllServerMap()
	} else {
		xmux.GetInstance(r).Response.(*pkg.Response).Data = controller.GetServersFromScripts(namesInterface.(map[string]struct{}))
	}
}

func GetScripts(w http.ResponseWriter, r *http.Request) {
	names := xmux.GetInstance(r).Get("scriptname")
	if names == nil {
		xmux.GetInstance(r).Response.(*pkg.Response).Data = store.Store.GetAllScriptMap()
	} else {
		xmux.GetInstance(r).Response.(*pkg.Response).Data = store.Store.GetScriptMapFilterByName(names.(map[string]struct{}))
	}
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	xmux.GetInstance(r).Response.(*pkg.Response).Data = store.Store.GetScriptIndex(pname)
}
