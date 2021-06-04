package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/server"
	"github.com/hyahm/xmux"
)

func AddScript(w http.ResponseWriter, r *http.Request) {
	s := xmux.GetData(r).Data.(*server.Script)
	golog.Infof("%#v", *s)
	if s.Name == "" {
		w.Write([]byte(`{"code": 201, "msg": "name require"}`))
		return
	}
	// 将时间转化为秒
	if s.ContinuityInterval != 0 {
		s.ContinuityInterval = s.ContinuityInterval * 1000000000
	}
	if server.HaveScript(s.Name) {
		// 修改
		// 需要判断是否相等
		golog.Info("aaaaa")
		err := server.UpdateScriptToConfigFile(s)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())))
			return
		}
		if !server.NeedStop(s) {
			w.Write([]byte(`{"code": 200, "msg": "already add script"}`))
			return
		}
		server.RemoveScript(s.Name)
	} else {
		golog.Info("bbbb")
		// 添加
		err := server.AddScriptToConfigFile(s)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())))
			return
		}

	}
	golog.Info("start and add server")
	server.AddScript(s)
	w.Write([]byte(`{"code": 200, "msg": "already add script"}`))
}

func DelScript(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	if err := server.Cfg.DelScript(pname); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(`{"code": 200, "msg": "already delete script"}`))
}
