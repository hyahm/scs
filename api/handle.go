package api

import (
	"github.com/hyahm/scs/api/handle"
	"github.com/hyahm/scs/api/module"
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/xmux"
)

func simpleHandle() *xmux.GroupRoute {
	// 只是查看的权限
	simple := xmux.NewGroupRoute().AddPageKeys(scripts.Simple.ToString())
	simple.Post("/status/{pname}/{name}", handle.Status)
	simple.Post("/status/{pname}", handle.StatusPname)
	simple.Post("/start/{pname}", handle.StartPname)
	simple.Post("/start/{pname}/{name}", handle.Start)
	simple.Post("/update/{pname}/{name}", handle.Update)
	simple.Post("/update/{pname}", handle.UpdatePname)
	simple.Post("/start", handle.StartAll)   // complete
	simple.Post("/status", handle.AllStatus) // complete
	simple.Get("/log/{name}/{int:line}", handle.Log)
	simple.Post("/update", handle.UpdateAll) // complete
	simple.Post("/restart/{pname}/{name}", handle.Restart)
	simple.Post("/restart/{pname}", handle.RestartPname)
	simple.Post("/restart", handle.RestartAll) // complete
	return simple
}

func scriptHandle() *xmux.GroupRoute {
	script := xmux.NewGroupRoute().AddPageKeys(scripts.ScriptRole.ToString())
	script.DelModule(module.CheckAdminToken).AddModule()
	script.Post("/stop/{pname}/{name}", handle.Stop)
	script.Post("/stop/{pname}", handle.StopPname)
	script.Post("/kill/{pname}", handle.KillPname)
	script.Post("/kill/{pname}/{name}", handle.Kill)
	script.Post("/env/{name}", handle.GetEnvName)
	script.Post("/server/info/{name}", handle.ServerInfo)
	script.Post("/get/servers", handle.GetServers)    // complete
	script.Post("/get/index/{name}", handle.GetIndex) // complete
	script.Post("/cannotstop/{name}", handle.CanNotStop)
	script.Post("/canstop/{name}", handle.CanStop)
	script.Post("/get/scripts", handle.GetScripts)
	script.Post("/stop", handle.StopAll) // complete
	script.Post("/script", handle.AddScript).BindJson(&scripts.Script{})
	script.Post("/remove/{pname}/{name}", handle.Remove)
	script.Post("/remove/{pname}", handle.RemovePname)
	script.Post("/send/alert", handle.Alert)
	script.AddGroup(simpleHandle())
	return script
}

func AdminHandle() *xmux.GroupRoute {
	// 只能管理员操作 修改文件的操作
	admin := xmux.NewGroupRoute().AddPageKeys(scripts.AdminRole.ToString()).AddModule(module.CheckAdminToken)

	admin.Post("/get/alert", handle.GetAlert)   // 只能管理员用
	admin.Post("/-/reload", handle.Reload)      // 只能管理员用
	admin.Post("/-/fmt", handle.Fmt)            // 只能管理员用
	admin.Post("/get/alarms", handle.GetAlarms) // 只能管理员用
	admin.Post("/get/repo", handle.GetRepo)     // 只能管理员用

	admin.Post("/enable/{pname}", handle.Enable)   // 只能管理员用
	admin.Post("/disable/{pname}", handle.Disable) // 只能管理员用

	admin.AddGroup(scriptHandle())
	return admin
}
