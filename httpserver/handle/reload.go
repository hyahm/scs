package handle

import (
	"net/http"

	"github.com/hyahm/scs/probe"
	"github.com/hyahm/scs/script"
)

var reloadKey bool

func Reload(w http.ResponseWriter, r *http.Request) {
	// 关闭上次监控的goroutine

	if reloadKey {
		w.Write([]byte(`{"code": 201, "msg": "config file is reloading, waiting completed first"}`))
		return
	}
	probe.Exit <- struct{}{}

	// script.Reloadlocker.Lock()
	// defer script.Reloadlocker.Unlock()
	reloadKey = true
	// 拷贝一份到当前运行的脚本列表
	script.Copy()
	if err := script.Load(true); err != nil {
		w.Write([]byte(err.Error()))
		reloadKey = false
		return
	}
	reloadKey = false
	w.Write([]byte(`{"code": 200, "msg": "config file reloaded"}`))
}
