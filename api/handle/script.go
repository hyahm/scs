package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/xmux"
)

// 这是一个添加script的handle
func AddScript(w http.ResponseWriter, r *http.Request) {
	s := xmux.GetInstance(r).Data.(*scripts.Script)
	if global.CanReload != 0 {
		w.Write(pkg.WaitingConfigChanged())
		return
	}
	res := pkg.Response{}
	if s.Name == "" {
		w.Write([]byte(`{"code": 404, "msg": "name not found"}`))
		return
	}
	if controller.HaveScript(s.Name) {
		// 存在的话，需要对比配置文件的修改
		// 需要判断是否相等
		if !controller.NeedStop(s) {
			// 如果没有修改，那么就返回已经add了
			w.Write(res.Sucess("already have script and the same"))
			return
		}
		// 更新配置文件

		// 更新这个script
		err := config.UpdateScriptToConfigFile(s, true)
		if err != nil {
			w.Write(res.ErrorE(err))
			return
		}
		// 否则删除原来的, 需不需要删除
		controller.UpdateScript(s, false)
	} else {
		// 添加
		err := config.AddScriptToConfigFile(s, true)
		if err != nil {
			w.Write(res.ErrorE(err))
			return
		}

		controller.AddScript(s)

	}
	res.Data = s.Token
	w.Write(res.Sucess("already add script"))
}
