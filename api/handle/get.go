package handle

import (
	"net/http"
)

func GetAlarms(w http.ResponseWriter, r *http.Request) {
	// xmux.GetInstance(r).Response.(*pkg.Response).Data = controller.GetAterts()
}

func GetServers(w http.ResponseWriter, r *http.Request) {
	// validAuths := xmux.GetInstance(r).Get("validAuths").([]controller.Auth)
	// svc := server.GetStore().GetAllServerMap()
	// validServer := make(map[string]*server.Server)
	// for _, auth := range validAuths {
	// 	if v, ok := svc[auth.ServerName]; ok {
	// 		validServer[auth.ServerName] = v
	// 	}

	// }
	// xmux.GetInstance(r).Response.(*pkg.Response).Data = validServer
}

func GetScripts(w http.ResponseWriter, r *http.Request) {
	// validAuths := xmux.GetInstance(r).Get("validAuths").([]controller.Auth)
	// ss := server.GetStore().GetAllScriptMap()
	// validScript := make(map[string]config.Script)
	// for _, auth := range internal.Cfg.Scripts {
	// 	if v, ok := ss[auth.ScriptName]; ok {
	// 		validScript[auth.ScriptName] = v
	// 	}
	// }
	// xmux.GetInstance(r).Response.(*pkg.Response).Data = internal.Cfg.Scripts
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	// pname := xmux.Var(r)["pname"]
	// xmux.GetInstance(r).Response.(*pkg.Response).Data = server.GetStore().GetScriptIndex(pname)
}
