package handle

import (
	"net/http"
)

func Disable(w http.ResponseWriter, r *http.Request) {

	// pname := xmux.Var(r)["pname"]
	// golog.Info("disable ", pname)
	// script, ok := store.GetStore().GetScriptByName(pname)
	// if !ok {
	// 	xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
	// 	return
	// }
	// // 上面已经判断过是否存在了， 这里就忽略
	// // msg, ok := global.SetReLoading(fmt.Sprintf("enable %s is running", pname))
	// // if !ok {
	// // 	pkg.Error(r, msg)
	// // 	return
	// // }
	// // defer global.SetCanReLoad()
	// // 上面已经判断过是否存在了， 这里就忽略
	// if controller.DisableScript(script, false) {
	// 	// err := config.UpdateScriptToConfigFile(script, true)
	// 	// if err != nil {
	// 	// 	xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
	// 	// 	return
	// 	// }
	// }

}

func Enable(w http.ResponseWriter, r *http.Request) {

	// pname := xmux.Var(r)["pname"]
	// script, ok := store.GetStore().GetScriptByName(pname)
	// if !ok {
	// 	xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
	// 	return
	// }
	// // 上面已经判断过是否存在了， 这里就忽略
	// // msg, ok := global.SetReLoading(fmt.Sprintf("enable %s is running", pname))
	// // if !ok {
	// // 	pkg.Error(r, msg)
	// // 	return
	// // }
	// // defer global.SetCanReLoad()
	// if controller.EnableScript(script) {
	// 	err := config.UpdateScriptToConfigFile(script, true)
	// 	if err != nil {
	// 		golog.Error(err)
	// 		xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
	// 		return
	// 	}
	// 	controller.AddScript(script)
	// }

}
