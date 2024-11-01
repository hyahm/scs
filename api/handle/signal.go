package handle

import (
	"context"
	"net/http"
	"time"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"

	"github.com/hyahm/xmux"
)

func CanStop(w http.ResponseWriter, r *http.Request) {
	name := xmux.Var(r)["name"]
	svc, ok := store.Store.GetServerByName(name)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	pkg.CancelAtomSignal(name)
	svc.Status.CanNotStop = false
}

func CanNotStop(w http.ResponseWriter, r *http.Request) {
	name := xmux.Var(r)["name"]
	sr := xmux.GetInstance(r).Data.(*pkg.SignalRequest)
	if sr.Timeout > 0 {
		c, cancel := context.WithCancel(context.Background())
		pkg.SetAtomSignal(name, cancel)
		// 如果大于0，创建一个goroutine 来监听超时
		controller.AddSignalRequest(name, sr)
		go controller.UnStop(c, name, time.Duration(sr.Timeout))
	}
	svc, ok := store.Store.GetServerByName(name)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	svc.Status.CanNotStop = true
}

func SetParameter(w http.ResponseWriter, r *http.Request) {
	name := xmux.Var(r)["name"]
	sr := xmux.GetInstance(r).Data.(*pkg.SignalRequest)
	if !controller.UpdateSignalRequest(name, sr) {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 406
	}
}
