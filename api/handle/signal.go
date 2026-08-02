package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg"

	"github.com/hyahm/xmux"
)

func CanStop(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	svc, ok := cache.GetServer(name)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	go func() {
		svc.Mu.Lock()
		defer svc.Mu.Unlock()
		if !svc.Status.CanNotStop {
			return
		}
svc.Status.CanNotStop = false
	close(svc.Status.ChStop)
	svc.Status.ChStop = make(chan struct{})
	}()

}

func CanNotStop(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	svc, ok := cache.GetServer(name)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	go func() {
		svc.Mu.Lock()
		defer svc.Mu.Unlock()
		svc.Status.CanNotStop = true
	}()

}

func SetParameter(w http.ResponseWriter, r *http.Request) {
	name := xmux.Var(r)["name"]
	sr := xmux.GetInstance(r).Data.(*pkg.SignalRequest)
	controller.UpdateSignalRequest(name, *sr)
}
