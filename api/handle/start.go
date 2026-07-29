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
		go start(svc)
		return
	}
	for _, svc := range cache.GetGroupServer(name) {
		go start(svc)

	}
	// controller.StartExsitScript(pname)
}

func StartAll(w http.ResponseWriter, r *http.Request) {
	for _, svc := range cache.GetAllServer() {
		go start(svc)
	}
}

func start(svc *cache.Server) {
	svc.Mu.Lock()
	defer svc.Mu.Unlock()
	if svc.CanNotOpration {
		golog.Info("服务暂时无法被操作: ", svc.SubName)
		return
	}
	svc.CanNotOpration = true
	svc.Start()
	svc.CanNotOpration = false
}
