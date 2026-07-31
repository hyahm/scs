package handle

import (
	"net/http"

	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/xmux"
)

func GetAlarms(w http.ResponseWriter, r *http.Request) {
	xmux.GetInstance(r).Response.(*pkg.Response).Data = config.GetDispatcherList()
}

func GetServers(w http.ResponseWriter, r *http.Request) {
	xmux.GetInstance(r).Response.(*pkg.Response).Data = cache.GetAllServerMap()
}

func GetScripts(w http.ResponseWriter, r *http.Request) {
	xmux.GetInstance(r).Response.(*pkg.Response).Data = internal.GetAllScript()
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	servers := cache.GetGroupServer(name)
	indices := make([]int, 0, len(servers))
	for _, svc := range servers {
		indices = append(indices, svc.Index)
	}
	xmux.GetInstance(r).Response.(*pkg.Response).Data = indices
}
