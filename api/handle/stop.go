package handle

import (
	"net/http"

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
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	svc, ok := store.Store.GetServerByName(name)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	go svc.Stop()
}

func StopPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	err := controller.StopScript(script)
	if err != nil {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
		return
	}
}

func StopAll(w http.ResponseWriter, r *http.Request) {
	validAuths := xmux.GetInstance(r).Get("validAuths").([]controller.Auth)
	validName := make(map[string]struct{})
	for _, auth := range validAuths {
		validName[auth.ScriptName] = struct{}{}
	}
	controller.StopScriptFromName(validName)

}
