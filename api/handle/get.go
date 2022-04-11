package handle

import (
	"net/http"

	"github.com/hyahm/scs/api/module"
	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func GetAlarms(w http.ResponseWriter, r *http.Request) {
	res := &pkg.Response{
		Data: controller.GetAterts(),
	}
	module.Write(w, r, res.Sucess(""))
}

func GetServers(w http.ResponseWriter, r *http.Request) {
	res := &pkg.Response{}
	namesInterface := xmux.GetInstance(r).Get("scriptname")
	if namesInterface == nil {
		res.Data = store.Store.GetAllServerMap()
	} else {
		res.Data = controller.GetServersFromScripts(namesInterface.(map[string]struct{}))
	}
	module.Write(w, r, res.Sucess(""))
}

func GetScripts(w http.ResponseWriter, r *http.Request) {
	res := &pkg.Response{}
	names := xmux.GetInstance(r).Get("scriptname")
	if names == nil {
		res.Data = store.Store.GetAllScriptMap()
	} else {
		res.Data = store.Store.GetScriptMapFilterByName(names.(map[string]struct{}))
	}
	module.Write(w, r, res.Sucess(""))
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["name"]
	res := &pkg.Response{}
	res.Data = store.Store.GetScriptIndex(pname)
	module.Write(w, r, res.Sucess(""))
}
