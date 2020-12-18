package handle

import (
	"net/http"
	"scs/config"
	"scs/probe"
	"scs/script"
)

var reloadKey bool

func Reload(w http.ResponseWriter, r *http.Request) {
	// 关闭上次监控的goroutine
	if reloadKey {
		w.Write([]byte(`{"code": 201, "msg": "config file is reloading, waiting completed first"}`))
		return
	}
	probe.GlobalProbe.Exit <- true

	script.Reloadlocker.Lock()
	defer script.Reloadlocker.Unlock()
	reloadKey = true
	// 拷贝一份到当前的脚本
	script.Copy()
	if err := config.Load(); err != nil {
		w.Write([]byte(err.Error()))
		reloadKey = false
		return
	}
	script.StopUnUseScript()
	reloadKey = false
	w.Write([]byte(`{"code": 200, "msg": "config file reloaded"}`))
}
