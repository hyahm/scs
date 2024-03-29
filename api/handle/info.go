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
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	xmux.GetInstance(r).Response.(*pkg.Response).Data = svc
}
