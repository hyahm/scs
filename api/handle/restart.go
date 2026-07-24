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
	if svc.GetCanNotOperation() {
		golog.Info("服务暂时无法被操作: ", svc.SubName)
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 201
		return
	}
	// 重启交由 svc.Restart 处理（内含 cannotstop 等待 + stop->start 语义）
	go func(svc *cache.Server) {
		svc.SetCanNotOperation(true)
		svc.Restart()
		svc.SetCanNotOperation(false)
	}(svc)
}

func RestartPname(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	for _, v := range cache.GetGroupServer(name) {
		if v.GetCanNotOperation() {
			golog.Info("服务暂时无法被操作: ", v.SubName)
			continue
		}
		go func(svc *cache.Server) {
			svc.SetCanNotOperation(true)
			svc.Restart()
			svc.SetCanNotOperation(false)
		}(v)
	}
}

func RestartAll(w http.ResponseWriter, r *http.Request) {
	for _, v := range cache.GetAllServer() {
		if v.GetCanNotOperation() {
			golog.Info("服务暂时无法被操作: ", v.SubName)
			continue
		}
		go func(svc *cache.Server) {
			svc.SetCanNotOperation(true)
			svc.Restart()
			svc.SetCanNotOperation(false)
		}(v)
	}
}
