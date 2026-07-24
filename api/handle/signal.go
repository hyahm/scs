package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg"

	"github.com/hyahm/xmux"
)

func CanStop(w http.ResponseWriter, r *http.Request) {
	subname := r.FormValue("subname")
	svc, ok := cache.GetServer(subname)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	if !svc.GetCanNotStop() {
		return
	}
	svc.SetCanNotStop(false)
	svc.Status.ChStop <- struct{}{}
}

func CanNotStop(w http.ResponseWriter, r *http.Request) {
	subname := r.FormValue("subname")
	svc, ok := cache.GetServer(subname)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	svc.SetCanNotStop(true)
}

func SetParameter(w http.ResponseWriter, r *http.Request) {
	name := xmux.Var(r)["name"]
	sr := xmux.GetInstance(r).Data.(*pkg.SignalRequest)
	controller.UpdateSignalRequest(name, *sr)
}
