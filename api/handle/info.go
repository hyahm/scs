package handle

import (
	"net/http"

	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func ServerInfo(w http.ResponseWriter, r *http.Request) {
	name := xmux.Var(r)["name"]
	svc, ok := store.Store.GetServerByName(name)
	if !ok {
		w.Write(pkg.NotFoundScript())
		return
	}
	res := pkg.Response{
		Data: svc,
	}
	w.Write(res.Sucess(""))
}
