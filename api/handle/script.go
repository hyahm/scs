package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/config"
	"github.com/hyahm/scs/internal/config/scripts"
	"github.com/hyahm/xmux"
)

func AddScript(w http.ResponseWriter, r *http.Request) {
	s := xmux.GetInstance(r).Data.(*scripts.Script)
	golog.Infof("%#v", *s)
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
			w.Write([]byte(`{"code": 200, "msg": "already have script and the same"}`))
			return
		}
		// 更新配置文件

		// 更新这个script

		err := config.UpdateScriptToConfigFile(s)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())))
			return
		}
		// 否则删除原来的, 需不需要删除
		controller.UpdateScript(s)
	} else {
		// 添加
		err := config.AddScriptToConfigFile(s)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())))
			return
		}

		controller.AddScript(s)

	}

	w.Write([]byte(`{"code": 200, "msg": "already add script"}`))
}

func DelScript(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	if err := controller.DelScript(pname); err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())))
		return
	}

	w.Write([]byte(`{"code": 200, "msg": "already delete script"}`))
}
