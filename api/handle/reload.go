package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func Reload(w http.ResponseWriter, r *http.Request) {
	// 关闭上次监控的goroutine
	role := xmux.GetInstance(r).Get("role").(string)
	if global.CanReload != 0 {
		w.Write(pkg.WaitingConfigChanged(role))
		return
	}
	res := pkg.Response{
		Role: role,
	}

	// 拷贝一份到当前运行的脚本列表
	if err := controller.Reload(); err != nil {
		w.Write(res.ErrorE(err))
		return
	}
	w.Write([]byte(`{"code": 200, "msg": "config file reloaded"}`))
}
