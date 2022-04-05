package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/pkg"
)

func Reload(w http.ResponseWriter, r *http.Request) {
	// 关闭上次监控的goroutine
	if global.CanReload != 0 {
		w.Write(pkg.WaitingConfigChanged())
		return
	}
	res := pkg.Response{}

	// 拷贝一份到当前运行的脚本列表
	if err := controller.Reload(); err != nil {
		w.Write(res.ErrorE(err))
		return
	}
	w.Write([]byte(`{"code": 200, "msg": "config file reloaded"}`))
}

func Fmt(w http.ResponseWriter, r *http.Request) {
	// 关闭上次监控的goroutine
	if global.CanReload != 0 {
		w.Write(pkg.WaitingConfigChanged())
		return
	}
	res := pkg.Response{}

	// 拷贝一份到当前运行的脚本列表
	if err := controller.Fmt(); err != nil {
		w.Write(res.ErrorE(err))
		return
	}
	w.Write([]byte(`{"code": 200, "msg": "config file reloaded"}`))
}
