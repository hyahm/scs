package handle

import (
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/xmux"
)

// 这是一个添加script的handle
func AddScript(w http.ResponseWriter, r *http.Request) {
	s := xmux.GetInstance(r).Data.(*config.Script)
	if s.Name == "" || s.Command == "" {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	script, ok := cache.GetScript(s.Name)
	if ok {
		if config.EqualScript(script, *s) {
			return
		}
		err := internal.UpdateScriptToConfigFile(*s, true)
		if err != nil {
			golog.Error(err)
			xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
			return
		}
		for _, svc := range cache.GetGroupServer(s.Name) {
			svc.Remove()
		}
		cache.RemoveGroupServer(s.Name)
		cache.RemoveScript(s.Name)
		cache.AddScript(*s)
		return
	}
	err := internal.AddScriptToConfigFile(s)
	if err != nil {
		golog.Error(err)
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
		return
	}
	cache.AddScript(*s)
}
