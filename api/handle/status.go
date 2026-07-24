package handle

import (
	"net/http"

	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

// func Status(w http.ResponseWriter, r *http.Request) {
// 	// pname := xmux.Var(r)["pname"]
// 	subname := r.FormValue("subname")
// 	svc, ok := cache.GetServer(subname)
// 	// status, err := controller.ScriptName(pname, name)
// 	if !ok {
// 		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
// 		return
// 	}

// 	xmux.GetInstance(r).Response.(*pkg.Response).Data = svc.GetStatus()
// }

func StatusPname(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if pkg.IsNameWithSuffix(name) {
		svc, ok := cache.GetServer(name)
		// status, err := controller.ScriptName(pname, name)
		if !ok {
			xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
			return
		}

		xmux.GetInstance(r).Response.(*pkg.Response).Data = svc.GetStatus()
		return
	}
	list := make([]pkg.ServiceStatus, 0)
	for _, v := range cache.GetAllServer() {
		if v.Name == name {
			list = append(list, v.GetStatus())
		}
	}
	xmux.GetInstance(r).Response.(*pkg.Response).Data = list
	// pname := xmux.Var(r)["pname"]
	// status, err := controller.ScriptPname(pname)
	// if err != nil {
	// 	xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
	// 	return
	// }
	// xmux.GetInstance(r).Response.(*pkg.Response).Data = status
}

func AllStatus(w http.ResponseWriter, r *http.Request) {
	// validAuths := xmux.GetInstance(r).Get("validAuths").([]controller.Auth)
	// validName := make(map[string]struct{})
	// for _, auth := range validAuths {
	// 	validName[auth.ScriptName] = struct{}{}
	// }
	list := make([]pkg.ServiceStatus, 0)
	for _, v := range cache.GetAllServer() {
		list = append(list, v.GetStatus())
	}
	// status, err := controller.ScriptName(pname, name)

	xmux.GetInstance(r).Response.(*pkg.Response).Data = list
}
