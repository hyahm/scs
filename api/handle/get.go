package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func GetAlarms(w http.ResponseWriter, r *http.Request) {
	res := &pkg.Response{
		Data: controller.GetAterts(),
	}
	w.Write(res.Sucess(""))
}

func GetServers(w http.ResponseWriter, r *http.Request) {
	res := &pkg.Response{}
	namesInterface := xmux.GetInstance(r).Get("scriptname")
	if namesInterface == nil {
		res.Data = controller.GetServers()
	} else {
		res.Data = controller.GetServersFromScripts(namesInterface.(map[string]struct{}))
	}

	w.Write(res.Sucess(""))
}

func GetScripts(w http.ResponseWriter, r *http.Request) {
	res := &pkg.Response{}
	names := xmux.GetInstance(r).Get("scriptname")
	if names == nil {
		res.Data = controller.GetScripts()

	} else {
		res.Data = controller.GetScriptsFromScritps(names.(map[string]struct{}))
	}

	w.Write(res.Sucess(""))
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["name"]
	res := &pkg.Response{}
	res.Data = controller.GetIndexs(pname)
	w.Write(res.Sucess(""))
}
