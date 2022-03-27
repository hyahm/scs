package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/config"
	"github.com/hyahm/scs/internal/config/scripts"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func AddScript(w http.ResponseWriter, r *http.Request) {
	s := xmux.GetInstance(r).Data.(*scripts.Script)
	role := xmux.GetInstance(r).Get("role").(string)
	res := pkg.Response{
		Role: role,
	}
	if s.Name == "" {
		w.Write([]byte(`{"code": 201, "msg": "name require"}`))
		return
	}
	// 将时间转化为秒
	if s.ContinuityInterval != 0 {
		s.ContinuityInterval = s.ContinuityInterval * 1000000000
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
	w.Write(res.Sucess("already add script"))
}

func DelScript(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	role := xmux.GetInstance(r).Get("role").(string)
	res := pkg.Response{
		Role: role,
	}

	if err := controller.DelScript(pname); err != nil {
		w.Write(res.ErrorE(err))
		return
	}

	w.Write(res.Sucess("already add script"))
}
