package handle

// func ShowConfig(w http.ResponseWriter, r *http.Request) {
// 	// 关闭上次监控的goroutine
// 	if reloadKey {
// 		w.Write([]byte("config file is reloading, waiting completed first"))
// 		return
// 	} else {
// 		w.Write([]byte("waiting config file reloaded"))
// 	}
// 	hardware.VarAT.Exit <- true

// 	script.Reloadlocker.Lock()
// 	defer script.Reloadlocker.Unlock()
// 	reloadKey = true
// 	// 拷贝一份到当前的脚本
// 	script.Copy()
// 	if err := config.Load(); err != nil {
// 		w.Write([]byte(err.Error()))
// 		reloadKey = false
// 		return
// 	}
// 	script.StopUnUseScript()
// 	reloadKey = false
// }
// 	defer script.Reloadlocker.Unlock()
// 	reloadKey = true
// 	// 拷贝一份到当前的脚本
// 	script.Copy()
// 	if err := config.Load(); err != nil {
// 		w.Write([]byte(err.Error()))
// 		reloadKey = false
// 		return
// 	}
// 	script.StopUnUseScript()
// 	reloadKey = false
// }
