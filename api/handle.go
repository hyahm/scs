package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/api/handle"
	"github.com/hyahm/scs/api/module"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/xmux"
)

func simpleHandle() *xmux.RouteGroup {
	// 只是调试的权限
	simple := xmux.NewRouteGroup().AddPageKeys(scripts.SimpleRole.ToString())
	simple.Post("/status/{pname}/{name}", handle.Status)
	simple.Post("/status/{pname}", handle.StatusPname)
	simple.Post("/start/{pname}", handle.StartPname)
	simple.Post("/start/{pname}/{name}", handle.Start)
	simple.Post("/update/{pname}/{name}", handle.Update)
	simple.Post("/update/{pname}", handle.UpdatePname)
	simple.Post("/update", handle.UpdateAll) // complete
	simple.Post("/start", handle.StartAll)   // complete
	simple.Post("/status", handle.AllStatus) // complete
	simple.Get("/log/{name}/{int:line}", handle.Log).BindResponse(nil)

	simple.Post("/restart/{pname}/{name}", handle.Restart)
	simple.Post("/restart/{pname}", handle.RestartPname)
	simple.Post("/restart", handle.RestartAll) // complete
	return simple
}

func ScriptHandle() *xmux.RouteGroup {
	script := xmux.NewRouteGroup().AddPageKeys(scripts.ScriptRole.ToString())
	script.Post("/stop/{pname}/{name}", handle.Stop)
	script.Post("/stop/{pname}", handle.StopPname)
	script.Post("/kill/{pname}", handle.KillPname)
	script.Post("/kill/{pname}/{name}", handle.Kill)
	script.Post("/env/{name}", handle.GetEnvName)
	script.Post("/server/info/{name}", handle.ServerInfo) // 获取某个server信息
	script.Post("/get/servers", handle.GetServers)        // 获取所有server信息 complete
	script.Post("/get/index/{pname}", handle.GetIndex)    // 获取某script 对应副本的index
	script.Post("/cannotstop/{name}", handle.CanNotStop).BindJson(pkg.SignalRequest{})
	script.Post("/parameter/{name}", handle.SetParameter).BindJson(pkg.SignalRequest{})
	script.Post("/canstop/{name}", handle.CanStop)
	script.Post("/get/scripts", handle.GetScripts) // complete
	script.Post("/stop", handle.StopAll)           // complete
	script.Post("/remove/{pname}/{name}", handle.Remove)
	script.Post("/remove/{pname}", handle.RemovePname)
	script.Post("/send/alert", handle.Alert)
	script.AddGroup(simpleHandle())
	return script
}

func AdminHandle() *xmux.RouteGroup {
	// 只能管理员操作 修改文件的操作
	admin := xmux.NewRouteGroup().AddPageKeys(scripts.AdminRole.ToString()).AddModule(module.CheckToken)

	admin.Post("/get/alert", handle.GetAlert)   // 只能管理员用
	admin.Post("/-/reload", handle.Reload)      // 只能管理员用
	admin.Post("/-/fmt", handle.Fmt)            // 只能管理员用
	admin.Post("/get/alarms", handle.GetAlarms) // 只能管理员用
	admin.Post("/get/repo", handle.GetRepo)     // 只能管理员用
	admin.Post("/script", handle.AddScript).BindJson(&scripts.Script{})
	admin.Post("/enable/{pname}", handle.Enable)   // 只能管理员用
	admin.Post("/disable/{pname}", handle.Disable) // 只能管理员用
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

func exit(start time.Time, w http.ResponseWriter, r *http.Request) {
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
	golog.Infof("method: %s\turl: %s\ttime: %f\t status_code: %v, body: %s, response: %s\n",
		r.Method,
		r.URL.Path, time.Since(start).Seconds(), xmux.GetInstance(r).StatusCode,
		string(xmux.GetInstance(r).Body),
		string(send))
}
