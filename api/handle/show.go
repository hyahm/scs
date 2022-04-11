package handle

// func ShowConfig(w http.ResponseWriter, r *http.Request) {
// 	// 关闭上次监控的goroutine
// 	if reloadKey {
// 		Write(w, r,[]byte("config file is reloading, waiting completed first"))
// 		return
// 	} else {
// 		Write(w, r,[]byte("waiting config file reloaded"))
// 	}
// 	hardware.VarAT.Exit <- true

// 	script.Reloadlocker.Lock()
// 	defer script.Reloadlocker.Unlock()
// 	reloadKey = true
// 	// 拷贝一份到当前的脚本
// 	script.Copy()
// 	if err := config.Load(); err != nil {
// 		Write(w, r,[]byte(err.Error()))
// 		reloadKey = false
// 		return
// 	}
// 	reloadKey = false
// }
// 	defer script.Reloadlocker.Unlock()
// 	reloadKey = true
// 	// 拷贝一份到当前的脚本
// 	script.Copy()
// 	if err := config.Load(); err != nil {
// 		Write(w, r,[]byte(err.Error()))
// 		reloadKey = false
// 		return
// 	}
// 	reloadKey = false
// }
