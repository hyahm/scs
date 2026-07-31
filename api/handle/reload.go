package handle

import (
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func Reload(w http.ResponseWriter, r *http.Request) {
	if err := internal.ReadConfig(); err != nil {
		golog.Error(err)
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
		return
	}
	oldScripts := cache.GetAllScript()
	newScripts := internal.GetAllScript()
	newMap := make(map[string]struct{})
	for _, s := range newScripts {
		newMap[s.Name] = struct{}{}
		if _, ok := oldScripts[s.Name]; !ok {
			cache.AddScript(s)
		}
	}
	for name := range oldScripts {
		if _, ok := newMap[name]; !ok {
			removeScriptByName(name)
		}
	}
}

func Fmt(w http.ResponseWriter, r *http.Request) {
	if err := internal.ReadConfig(); err != nil {
		golog.Error(err)
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
		return
	}
	cfg := internal.GetConfig()
	if err := internal.WriteConfig(cfg); err != nil {
		golog.Error(err)
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
		return
	}
}

func removeScriptByName(name string) {
	for _, svc := range cache.GetGroupServer(name) {
		svc.Remove()
	}
	cache.RemoveGroupServer(name)
	cache.RemoveScript(name)
}
