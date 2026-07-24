package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/api/handle"
	"github.com/hyahm/scs/api/module"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/xmux"
)

func simpleHandle() *xmux.RouteGroup {
	// 只是调试的权限
	simple := xmux.NewRouteGroup().AddPageKeys(config.SimpleRole.ToString())
	simple.Get("/status/name", handle.StatusPname)
	simple.Get("/start/name", handle.StartPname)
	simple.Get("/update/name", handle.UpdatePname)
	simple.Get("/update", handle.UpdateAll) // complete
	simple.Get("/start", handle.StartAll)   // complete
	simple.Get("/status", handle.AllStatus) // complete
	simple.Get("/log/name", handle.Log).BindResponse(nil)

	simple.Post("/restart/name", handle.RestartPname)
	simple.Post("/restart", handle.RestartAll) // complete
	return simple
}

func ScriptHandle() *xmux.RouteGroup {
	script := xmux.NewRouteGroup().AddPageKeys(config.ScriptRole.ToString())
	script.Get("/stop/name", handle.StopPname)
	script.Get("/stop", handle.StopAll) // complete
	script.Get("/kill/name", handle.KillPname)
	script.Get("/server/info/name", handle.ServerInfo) // 获取某个server信息
	script.Get("/get/servers", handle.GetServers)      // 获取所有server信息 complete
	script.Get("/get/index/name", handle.GetIndex)     // 获取某script 对应副本的index
	script.Get("/cannotstop/name", handle.CanNotStop)
	// .BindJson(pkg.SignalRequest{})
	script.Get("/parameter/name", handle.SetParameter).BindJson(pkg.SignalRequest{})
	script.Get("/canstop/name", handle.CanStop)
	script.Get("/get/scripts", handle.GetScripts) // complete

	// script.Get("/remove/name", handle.Remove).AddModule(module.UpdateConfig)
	script.Get("/remove/name", handle.RemovePname).AddModule(module.UpdateConfig)
	script.Get("/send/alert", handle.Alert)
	script.AddGroup(simpleHandle())
	return script
}

func AdminHandle() *xmux.RouteGroup {
	// 只能管理员操作 修改文件的操作
	admin := xmux.NewRouteGroup().AddPageKeys(config.AdminRole.ToString()).AddModule(module.CheckToken)

	admin.Get("/get/alert", handle.GetAlert)                             // 只能管理员用
	admin.Get("/-/reload", handle.Reload).AddModule(module.UpdateConfig) // 只能管理员用
	admin.Get("/-/fmt", handle.Fmt).AddModule(module.UpdateConfig)       // 只能管理员用
	admin.Get("/get/alarms", handle.GetAlarms)                           // 只能管理员用
	admin.Get("/get/repo", handle.GetRepo)                               // 只能管理员用
	admin.Get("/script", handle.AddScript).BindJson(&config.Script{}).AddModule(module.UpdateConfig)
	admin.Get("/enable/name", handle.Enable).AddModule(module.UpdateConfig)   // 只能管理员用
	admin.Get("/disable/name", handle.Disable).AddModule(module.UpdateConfig) // 只能管理员用
	admin.AddGroup(ScriptHandle())
	return admin
}

var statusMsg map[int]string

func init() {
	statusMsg = make(map[int]string)
	statusMsg[200] = "ok"
	statusMsg[201] = "config is reloading, please wait"
	statusMsg[203] = "token error"
	statusMsg[404] = "pname or name not found"
	statusMsg[406] = "没有找到对应运行的信号"
	statusMsg[500] = "system error"
}

func exit(w http.ResponseWriter, r *http.Request, start time.Time) {
	var send []byte
	var err error

	if xmux.GetInstance(r).Response != nil && xmux.GetInstance(r).StatusCode == 200 {
		response := xmux.GetInstance(r).Response.(*pkg.Response)
		// response.Code = xmux.GetInstance(r).Get(xmux.STATUSCODE).(int)
		response.Msg = statusMsg[response.Code]
		send, err = json.Marshal(response)
		if err != nil {
			golog.Error(err)
		}
		w.Write(send)
	}
	// golog.Debugf("method: %s\turl: %s\ttime: %f\t status_code: %v, body: %s, response: %s",
	// 	r.Method,
	// 	r.URL.Path, time.Since(start).Seconds(), xmux.GetInstance(r).StatusCode,
	// 	string(xmux.GetInstance(r).Body),
	// 	string(send))
}
