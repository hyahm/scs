package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/script"
	"github.com/hyahm/xmux"
)

func AddScript(w http.ResponseWriter, r *http.Request) {
	s := xmux.GetData(r).Data.(*script.Script)
	if s.Name == "" {
		w.Write([]byte(`{"code": 201, "msg": "name require"}`))
		return
	}
	// 将时间转化为秒
	if s.ContinuityInterval != 0 {
		s.ContinuityInterval = s.ContinuityInterval * 1000000000
	}
	golog.Info("11111")
	if script.HaveScript(s.Name) {
		// 修改
		// 需要判断是否相等

		err := script.UpdateScriptToConfigFile(s)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())))
			return
		}
		if !s.NeedReStop() {
			w.Write([]byte(`{"code": 200, "msg": "already add script"}`))
			return
		}
		s.RemoveScript()
	} else {
		// 添加
		err := script.AddScriptToConfigFile(s)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())))
			return
		}

	}
	s.AddScript()
	w.Write([]byte(`{"code": 200, "msg": "already add script"}`))
}

func DelScript(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	if err := script.Cfg.DelScript(pname); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(`{"code": 200, "msg": "already delete script"}`))
	return
}
