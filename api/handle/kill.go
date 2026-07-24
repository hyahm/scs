package handle

import (
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

// Kill 按 subname 强杀单个服务副本（未注册路由，保留供按副本名强杀时使用）
func Kill(w http.ResponseWriter, r *http.Request) {
	subname := r.FormValue("subname")
	svc, ok := cache.GetServer(subname)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	if svc.GetCanNotOperation() {
		golog.Info("服务暂时无法被操作: ", svc.SubName)
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 201
		return
	}
	if svc.Status.Status == pkg.RUNNING {
		go func(svc *cache.Server) {
			svc.SetCanNotOperation(true)
			svc.Kill()
			svc.SetCanNotOperation(false)
		}(svc)
	}
}

// KillPname 按 name 强杀：name 带后缀时强杀单个副本，否则强杀整组。
func KillPname(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if pkg.IsNameWithSuffix(name) {
		svc, ok := cache.GetServer(name)
		if !ok {
			xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
			return
		}
		if svc.GetCanNotOperation() {
			golog.Info("服务暂时无法被操作: ", svc.SubName)
			xmux.GetInstance(r).Response.(*pkg.Response).Code = 201
			return
		}
		if svc.Status.Status == pkg.RUNNING {
			go func(svc *cache.Server) {
				svc.SetCanNotOperation(true)
				svc.Kill()
				svc.SetCanNotOperation(false)
			}(svc)
		}
		return
	}
	for _, v := range cache.GetGroupServer(name) {
		if v.GetCanNotOperation() {
			golog.Info("服务暂时无法被操作: ", v.SubName)
			continue
		}
		if v.Status.Status == pkg.RUNNING {
			go func(svc *cache.Server) {
				v.SetCanNotOperation(true)
				v.Kill()
				v.SetCanNotOperation(false)
			}(v)
		}
	}
}
