package handle

import (
	"net/http"

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
	// 先停止，再删除， 异步处理
	go func(svc *cache.Server) {
		svc.SetCanNotOperation(true)
		svc.Stop()
		svc.Start()
		svc.SetCanNotOperation(false)
	}(svc)
}

func RestartPname(w http.ResponseWriter, r *http.Request) {
	// pname := xmux.Var(r)["pname"]
	name := r.FormValue("name")
	for _, v := range cache.GetGroupServer(name) {
		// 先停止，再删除， 异步处理
		go func(svc *cache.Server) {
			svc.SetCanNotOperation(true)
			svc.Stop()
			svc.Start()
			svc.SetCanNotOperation(false)
		}(v)
	}

}

func RestartAll(w http.ResponseWriter, r *http.Request) {
	// 删除所有的脚本
	for _, v := range cache.GetAllServer() {
		// 先停止，再删除， 异步处理
		go func(svc *cache.Server) {
			svc.SetCanNotOperation(true)
			svc.Stop()
			svc.Start()
			svc.SetCanNotOperation(false)
		}(v)
	}
}
