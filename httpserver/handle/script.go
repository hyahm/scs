package handle

import (
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/script"
	"github.com/hyahm/xmux"
)

func AddScript(w http.ResponseWriter, r *http.Request) {
	s := xmux.GetData(r).Data.(*internal.Script)
	golog.Infof("%+v", s)
	if s.Name == "" {
		w.Write([]byte(`{"code": 201, "msg": "name require"}`))
		return
	}
	// 将时间转化为秒
	s.ContinuityInterval = s.ContinuityInterval * 1000000000
	if err := script.Cfg.AddScript(*s); err != nil {
		w.Write([]byte(`{"code": 201, "msg": "already exist script"}`))
		return
	}
	w.Write([]byte(`{"code": 200, "msg": "already add script"}`))
	return
}

func DelScript(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	if err := script.Cfg.DelScript(pname); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	// if reloadKey {
	// 	w.Write([]byte(`{"code": 201, "msg": "config file is reloading, waiting completed first"}`))
	// 	return
	// }

	// script.Reloadlocker.Lock()
	// defer script.Reloadlocker.Unlock()
	// reloadKey = true
	// // 拷贝一份到当前的脚本
	// script.Copy()
	// if err := config.Load(); err != nil {
	// 	w.Write([]byte(err.Error()))
	// 	reloadKey = false
	// 	return
	// }
	// script.RemoveUnUseScript()
	// reloadKey = false
	w.Write([]byte(`{"code": 200, "msg": "already delete script"}`))
	return
}
