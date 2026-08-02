package handle

import (
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func Update(w http.ResponseWriter, r *http.Request) {
	subname := r.FormValue("subname")
	svc, ok := cache.GetServer(subname)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	go func() {
		svc.Mu.Lock()
		if svc.CanNotOpration {
			svc.Mu.Unlock()
			return
		}
		svc.CanNotOpration = true
		svc.Mu.Unlock()
		svc.UpdateServer()
		svc.Mu.Lock()
		svc.CanNotOpration = false
		svc.Mu.Unlock()
	}()
}

func UpdatePname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	_, ok := cache.GetScript(pname)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	for _, svc := range cache.GetGroupServer(pname) {
		go func(s *cache.Server) {
			s.Mu.Lock()
			if s.CanNotOpration {
				s.Mu.Unlock()
				return
			}
			s.CanNotOpration = true
			s.Mu.Unlock()
			s.UpdateServer()
			s.Mu.Lock()
			s.CanNotOpration = false
			s.Mu.Unlock()
		}(svc)
	}
}

func UpdateAll(w http.ResponseWriter, r *http.Request) {
	for _, svc := range cache.GetAllServer() {
		go func(s *cache.Server) {
			s.Mu.Lock()
			if s.CanNotOpration {
				s.Mu.Unlock()
				return
			}
			s.CanNotOpration = true
			s.Mu.Unlock()
			s.UpdateServer()
			s.Mu.Lock()
			s.CanNotOpration = false
			s.Mu.Unlock()
		}(svc)
	}
	golog.Info("update all servers done")
}
