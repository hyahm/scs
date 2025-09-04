package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func Reload(w http.ResponseWriter, r *http.Request) {
	// 关闭上次监控的goroutine

	// 拷贝一份到当前运行的脚本列表
	if err := controller.Reload(); err != nil {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
		return
	}
}

func Fmt(w http.ResponseWriter, r *http.Request) {
	// 关闭上次监控的goroutine

	// 拷贝一份到当前运行的脚本列表
	if err := controller.Fmt(); err != nil {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
		return
	}
}
