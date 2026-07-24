package handle

import (
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

// func Stop(w http.ResponseWriter, r *http.Request) {
// 	// pname := xmux.Var(r)["pname"]
// 	// name := xmux.Var(r)["name"]
// 	name := r.FormValue("subname")
// 	// _, ok := store.GetStore().GetScriptByName(pname)
// 	// if !ok {
// 	// 	xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
// 	// 	return
// 	// }
// 	svc, ok := cache.GetServer(name)
// 	if !ok {
// 		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
// 		return
// 	}
// 	if svc.GetCanNotOperation() {
// 		golog.Info("服务暂时无法被操作: ", svc.SubName)
// 		xmux.GetInstance(r).Response.(*pkg.Response).Code = 201
// 		return

// 	}
// 	if svc.Status.Status == pkg.RUNNING {
// 		go func(svc *cache.Server) {
// 			svc.SetCanNotOperation(true)
// 			svc.Stop()
// 			svc.SetCanNotOperation(false)
// 		}(svc)
// 	}
// }

func StopPname(w http.ResponseWriter, r *http.Request) {
	// pname := xmux.Var(r)["pname"]
	name := r.FormValue("name")
	// script, ok := store.GetStore().GetScriptByName(pname)
	// if !ok {
	// 	xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
	// 	return
	// }
	if pkg.IsNameWithSuffix(name) {
		svc, ok := cache.GetServer(name)
		if !ok {
			xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
			return
		}
		if svc.GetCanNotOperation() {
			golog.Info("服务暂时无法被操作: ", svc.SubName)
			xmux.GetInstance(r).Response.(*pkg.Response).Code = 201
			return

		}
		if svc.Status.Status == pkg.RUNNING {
			go func(svc *cache.Server) {
				svc.SetCanNotOperation(true)
				svc.Stop()
				svc.SetCanNotOperation(false)
			}(svc)
		}
		return
	}
	for _, v := range cache.GetGroupServer(name) {
		if v.GetCanNotOperation() {
			golog.Info("服务暂时无法被操作: ", v.SubName)
			continue

		}
		if v.Status.Status == pkg.RUNNING {
			go func(svc *cache.Server) {
				v.SetCanNotOperation(true)
				v.Stop()
				v.SetCanNotOperation(false)
			}(v)
		}
	}
	// err := controller.StopScript(script)
	// if err != nil {
	// 	xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
	// 	return
	// }
}

func StopAll(w http.ResponseWriter, r *http.Request) {
	// validAuths := xmux.GetInstance(r).Get("validAuths").([]controller.Auth)
	// validName := make(map[string]struct{})
	// for _, auth := range validAuths {
	// 	validName[auth.ScriptName] = struct{}{}
	// }
	// controller.StopScriptFromName(validName)
	for _, v := range cache.GetAllServer() {
		if v.GetCanNotOperation() {
			golog.Info("服务暂时无法被操作: ", v.SubName)
			continue
		}
		if v.Status.Status == pkg.RUNNING {
			go func(svc *cache.Server) {
				v.SetCanNotOperation(true)
				v.Stop()
				v.SetCanNotOperation(false)
			}(v)
		}
	}

}
