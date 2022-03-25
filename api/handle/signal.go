package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/config/scripts/subname"

	"github.com/hyahm/xmux"
)

func CanStop(w http.ResponseWriter, r *http.Request) {

	// golog.Info(string(res))
	name := xmux.Var(r)["name"]
	svc, ok := controller.GetServerBySubname(subname.Subname(name).String())
	if !ok {
		w.Write([]byte(`{"code": 200, "msg": "not found this name"}`))
		return
	}
	svc.Status.CanNotStop = false
	w.Write([]byte(`{"code": 200, "msg": "now can stop"}`))
}

func CanNotStop(w http.ResponseWriter, r *http.Request) {

	// golog.Info(string(res))
	name := xmux.Var(r)["name"]
	svc, ok := controller.GetServerBySubname(subname.Subname(name).String())
	if !ok {
		w.Write([]byte(`{"code": 200, "msg": "not found this name"}`))
		return
	}
	svc.Status.CanNotStop = true

	w.Write([]byte(`{"code": 200, "msg": "now can not stop"}`))
}
