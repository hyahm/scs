package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/scripts/subname"

	"github.com/hyahm/xmux"
)

func CanStop(w http.ResponseWriter, r *http.Request) {

	// golog.Info(string(res))
	name := xmux.Var(r)["name"]
	role := xmux.GetInstance(r).Get("role").(string)
	svc, ok := controller.GetServerBySubname(subname.Subname(name).String())
	if !ok {
		w.Write(pkg.NotFoundScript(role))
		return
	}
	svc.Status.CanNotStop = false
	w.Write([]byte(`{"code": 200, "msg": "now can stop"}`))
}

func CanNotStop(w http.ResponseWriter, r *http.Request) {

	// golog.Info(string(res))
	name := xmux.Var(r)["name"]
	role := xmux.GetInstance(r).Get("role").(string)
	svc, ok := controller.GetServerBySubname(subname.Subname(name).String())
	if !ok {
		w.Write(pkg.NotFoundScript(role))
		return
	}
	svc.Status.CanNotStop = true

	w.Write([]byte(`{"code": 200, "msg": "now can not stop"}`))
}
