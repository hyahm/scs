package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/xmux"
)

// 这是一个添加script的handle
func AddScript(w http.ResponseWriter, r *http.Request) {
	s := xmux.GetInstance(r).Data.(*scripts.Script)

	if s.Name == "" {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	msg, ok := global.SetReLoading(fmt.Sprintf("add script %s ", s.Name))
	if !ok {
		pkg.Error(r, msg)
		return
	}
	_, ok = store.Store.GetScriptByName(s.Name)
	if ok {
		// 存在的话，需要对比配置文件的修改
		// 需要判断是否相等
		if !controller.NeedStop(s) {
			// 如果没有修改，那么就返回已经add了
			return
		}

		// 更新配置文件
		err := config.UpdateScriptToConfigFile(s, true)
		if err != nil {
			xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
			return
		}
		// 否则删除原来的, 需不需要删除
		controller.UpdateScriptApi(s)
	} else {
		// 添加
		err := config.AddScriptToConfigFile(s)
		if err != nil {
			xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
			return
		}

		controller.AddScript(s)

	}
	global.SetCanReLoad()
	xmux.GetInstance(r).Response.(*pkg.Response).Data = s.ScriptToken
}
