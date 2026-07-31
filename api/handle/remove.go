package handle

import (
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func RemovePname(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	_, ok := cache.GetScript(name)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	go func() {
		for _, svc := range cache.GetGroupServer(name) {
			svc.Remove()
		}
		cache.RemoveGroupServer(name)
		cache.RemoveScript(name)
		if err := internal.DeleteScriptFromConfigFile(name); err != nil {
			golog.Error(err)
		}
	}()
}
