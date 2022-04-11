package handle

import (
	"net/http"

	"github.com/hyahm/scs/api/module"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"

	"github.com/hyahm/xmux"
)

func CanStop(w http.ResponseWriter, r *http.Request) {
	name := xmux.Var(r)["name"]
	svc, ok := store.Store.GetServerByName(name)
	if !ok {
		module.Write(w, r, pkg.NotFoundScript())
		return
	}
	svc.Status.CanNotStop = false
	module.Write(w, r, []byte(`{"code": 200, "msg": "now can stop"}`))
}

func CanNotStop(w http.ResponseWriter, r *http.Request) {

	// golog.Info(string(res))
	name := xmux.Var(r)["name"]
	svc, ok := store.Store.GetServerByName(name)
	if !ok {
		module.Write(w, r, pkg.NotFoundScript())
		return
	}
	svc.Status.CanNotStop = true

	module.Write(w, r, []byte(`{"code": 200, "msg": "now can not stop"}`))
}
