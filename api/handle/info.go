package handle

import (
	"net/http"

	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func ServerInfo(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	svc, ok := cache.GetServer(name)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	xmux.GetInstance(r).Response.(*pkg.Response).Data = svc.GetStatus()
}
