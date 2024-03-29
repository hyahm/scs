package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"

	"github.com/hyahm/xmux"
)

// 只有状态为 stop的 才会启动

func Start(w http.ResponseWriter, r *http.Request) {
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
	// rp := xmux.GetInstance(r).Data.(*pkg.ReStartParameter)
	svc.Start()

}

func StartPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	_, ok := store.Store.GetScriptByName(pname)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	controller.StartExsitScript(pname)
}

func StartAll(w http.ResponseWriter, r *http.Request) {
	validAuths := xmux.GetInstance(r).Get("validAuths").([]controller.Auth)
	validName := make(map[string]struct{})
	for _, auth := range validAuths {
		validName[auth.ScriptName] = struct{}{}
	}
	controller.StartAllServerFromScript(validName)

}
