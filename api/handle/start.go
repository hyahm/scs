package handle

import (
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

// 只有状态为 stop的 才会启动

// func Start(w http.ResponseWriter, r *http.Request) {
// 	subname := r.FormValue("subname")
// 	golog.Debug(subname)
// 	svc, ok := cache.GetServer(subname)
// 	if !ok {
// 		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
// 		return
// 	}
// 	if svc.GetCanNotOperation() {
// 		golog.Info("服务暂时无法被操作: ", svc.SubName)
// 		return
// 	}
// 	go func() {
// 		svc.SetCanNotOperation(true)
// 		svc.Start()
// 		svc.SetCanNotOperation(false)
// 	}()

// }

func StartPname(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if pkg.IsNameWithSuffix(name) {
		svc, ok := cache.GetServer(name)
		if !ok {
			xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
			return
		}
		if svc.GetCanNotOperation() {
			golog.Info("服务暂时无法被操作: ", name)
			return
		}
		go func() {
			svc.SetCanNotOperation(true)
			svc.Start()
			svc.SetCanNotOperation(false)
		}()
		return
	}
	for _, v := range cache.GetGroupServer(name) {
		if v.GetCanNotOperation() {
			golog.Info("服务暂时无法被操作: ", v.SubName)
			continue
		}
		if v.Status.Status == pkg.STOP {
			go func(svc *cache.Server) {
				svc.SetCanNotOperation(true)
				svc.Start()
				svc.SetCanNotOperation(false)
			}(v)
		}

	}
	// controller.StartExsitScript(pname)
}

func StartAll(w http.ResponseWriter, r *http.Request) {
	for _, v := range cache.GetAllServer() {
		if v.GetCanNotOperation() {
			golog.Info("服务暂时无法被操作: ", v.SubName)
			continue
		}
		if v.Status.Status == pkg.STOP {
			go func(svc *cache.Server) {
				svc.SetCanNotOperation(true)
				svc.Start()
				svc.SetCanNotOperation(false)
			}(v)
		}
	}
}
