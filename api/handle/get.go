package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/xmux"
)

func GetAlarms(w http.ResponseWriter, r *http.Request) {
	xmux.GetInstance(r).Response.(*pkg.Response).Data = controller.GetAterts()
}

func GetServers(w http.ResponseWriter, r *http.Request) {
	validAuths := xmux.GetInstance(r).Get("validAuths").([]controller.Auth)
	svc := store.Store.GetAllServerMap()
	validServer := make(map[string]*server.Server)
	for _, auth := range validAuths {
		if v, ok := svc[auth.ServerName]; ok {
			validServer[auth.ServerName] = v
		}

	}
	xmux.GetInstance(r).Response.(*pkg.Response).Data = validServer
}

func GetScripts(w http.ResponseWriter, r *http.Request) {
	validAuths := xmux.GetInstance(r).Get("validAuths").([]controller.Auth)
	ss := store.Store.GetAllScriptMap()
	validScript := make(map[string]*scripts.Script)
	for _, auth := range validAuths {
		if v, ok := ss[auth.ScriptName]; ok {
			validScript[auth.ScriptName] = v
		}
	}
	xmux.GetInstance(r).Response.(*pkg.Response).Data = validScript
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	xmux.GetInstance(r).Response.(*pkg.Response).Data = store.Store.GetScriptIndex(pname)
}
