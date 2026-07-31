package handle

import (
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func Disable(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	script, ok := cache.GetScript(name)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	if script.Disable {
		return
	}
	script.Disable = true
	cache.SetScript(name, script)
	if err := internal.SetScriptDisableInConfig(name, true); err != nil {
		golog.Error(err)
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
		return
	}
	for _, svc := range cache.GetGroupServer(name) {
		go stop(svc)
	}
}

func Enable(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	script, ok := cache.GetScript(name)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	if !script.Disable {
		return
	}
	script.Disable = false
	cache.SetScript(name, script)
	if err := internal.SetScriptDisableInConfig(name, false); err != nil {
		golog.Error(err)
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
		return
	}
	for _, svc := range cache.GetGroupServer(name) {
		svc.Disable = false
		go start(svc)
	}
}
