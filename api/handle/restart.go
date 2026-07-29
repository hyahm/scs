package handle

import (
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func Restart(w http.ResponseWriter, r *http.Request) {
	subname := r.FormValue("subname")
	svc, ok := cache.GetServer(subname)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	go restart(svc)
}

func RestartPname(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	for _, svc := range cache.GetGroupServer(name) {
		go restart(svc)
	}
}

func RestartAll(w http.ResponseWriter, r *http.Request) {
	for _, svc := range cache.GetAllServer() {
		go restart(svc)
	}
}

func restart(svc *cache.Server) {
	svc.Mu.Lock()
	defer svc.Mu.Unlock()
	if svc.CanNotOpration {
		golog.Info("服务暂时无法被操作: ", svc.SubName)
		return
	}
	svc.CanNotOpration = true
	svc.Restart()
	svc.CanNotOpration = false
}
