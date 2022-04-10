package handle

import (
	"net/http"

	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"

	"github.com/hyahm/xmux"
)

func CanStop(w http.ResponseWriter, r *http.Request) {

	// golog.Info(string(res))
	name := xmux.Var(r)["name"]
	svc, ok := store.Store.GetServerByName(name)
	if !ok {
		w.Write(pkg.NotFoundScript())
		return
	}
	svc.Status.CanNotStop = false
	w.Write([]byte(`{"code": 200, "msg": "now can stop"}`))
}

func CanNotStop(w http.ResponseWriter, r *http.Request) {

	// golog.Info(string(res))
	name := xmux.Var(r)["name"]
	svc, ok := store.Store.GetServerByName(name)
	if !ok {
		w.Write(pkg.NotFoundScript())
		return
	}
	svc.Status.CanNotStop = true

	w.Write([]byte(`{"code": 200, "msg": "now can not stop"}`))
}
