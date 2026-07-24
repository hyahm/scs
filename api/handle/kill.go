package handle

import (
	"net/http"
)

func Kill(w http.ResponseWriter, r *http.Request) {
	// pname := xmux.Var(r)["pname"]
	// name := xmux.Var(r)["name"]
	// _, ok := store.GetStore().GetScriptByName(pname)
	// if !ok {
	// 	xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
	// 	return
	// }
	// svc, ok := store.GetStore().GetServerByName(name)
	// if !ok {
	// 	xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
	// 	return
	// }
	// svc.Kill()
}

func KillPname(w http.ResponseWriter, r *http.Request) {
	// pname := xmux.Var(r)["pname"]
	// script, ok := store.GetStore().GetScriptByName(pname)
	// if !ok {
	// 	xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
	// 	return
	// }

	// controller.KillScript(script)
}
