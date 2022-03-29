package api

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hyahm/scs/api/handle"
	"github.com/hyahm/scs/api/midware"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/pkg/config/scripts"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func enter(w http.ResponseWriter, r *http.Request) bool {
	var addr string
	if global.ProxyHeader == "" {
		addr = strings.Split(r.RemoteAddr, ":")[0]
	} else {
		addr = r.Header.Get(global.ProxyHeader)
	}

	golog.Infof("url: %s -- addr: %s -- method: %s", r.URL.Path, addr, r.Method)
	return false
}

func LookHandle() *xmux.GroupRoute {
	// 只是查看的权限
	group := xmux.NewGroupRoute().AddPageKeys("scripts")
	group.Post("/status/{pname}/{name}", handle.Status)
	group.Post("/status/{pname}", handle.StatusPname)
	group.Post("/start/{pname}", handle.StartPname)
	group.Post("/start/{pname}/{name}", handle.Start)
	group.Post("/stop/{pname}/{name}", handle.Stop)
	group.Post("/stop/{pname}", handle.StopPname)
	group.Post("/restart/{pname}/{name}", handle.Restart)
	group.Post("/restart/{pname}", handle.RestartPname)
	group.Post("/update/{pname}/{name}", handle.Update)
	group.Post("/update/{pname}", handle.UpdatePname)
	group.Post("/kill/{pname}", handle.KillPname)
	group.Post("/kill/{pname}/{name}", handle.Kill)
	group.Post("/env/{name}", handle.GetEnvName)
	group.Post("/server/info/{name}", handle.ServerInfo)
	group.Post("/get/servers", handle.GetServers) // complete

	group.Get("/log/{name}/{int:line}", handle.Log)

	group.Post("/cannotstop/{name}", handle.CanNotStop)
	group.Post("/canstop/{name}", handle.CanStop)

	group.Post("/get/scripts", handle.GetScripts)
	group.Post("/start", handle.StartAll)   // complete
	group.Post("/status", handle.AllStatus) // complete
	group.Post("/stop", handle.StopAll)     // complete

	group.Post("/restart", handle.RestartAll) // complete
	group.Post("/update", handle.UpdateAll)   // complete
	return group
}

func FileHandle() *xmux.GroupRoute {
	// 修改文件的操作
	group := xmux.NewGroupRoute()
	// group.Post("/get/info", handle.GetOS)
	group.Post("/remove/{pname}/{name}", handle.Remove)
	group.Post("/remove/{pname}", handle.RemovePname)
	// router.Post("/remove", handle.RemoveAll)
	group.Post("/get/alert", handle.GetAlert)   // 只能管理员用
	group.Post("/-/reload", handle.Reload)      // 只能管理员用
	group.Post("/get/alarms", handle.GetAlarms) // 只能管理员用
	group.Post("/get/repo", handle.GetRepo)     // 只能管理员用

	group.Post("/enable/{pname}", handle.Enable)   // 只能管理员用
	group.Post("/disable/{pname}", handle.Disable) // 只能管理员用

	group.Post("/set/alert", handle.Alert) // 只能管理员用
	// 监测点

	// router.Get("/version/{pname}/{name}", handle.Version)
	group.Post("/script", handle.AddScript).BindJson(&scripts.Script{})
	// group.Post("/delete/{pname}", handle.DelScript)
	return group
}

// var dir := "key"
func HttpServer() {
	router := xmux.NewRouter()
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.SetHeader("Content-Type", "application/x-www-form-urlencoded,application/json; charset=UTF-8")
	router.SetHeader("Access-Control-Allow-Headers", "Content-Type")
	router.AddModule(midware.CheckToken)
	router.Post("/probe", handle.Probe).DelModule(midware.CheckToken)
	// 增加请求时间
	router.Enter = enter
	// router.MiddleWare(GetExecTime)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	router.AddGroup(LookHandle())
	router.AddGroup(FileHandle())
	if global.GetDisableTls() {
		golog.Info("listen on " + global.GetListen() + " over http")
		golog.Fatal(router.Run(global.GetListen()))
	}

	svc := &http.Server{
		ReadTimeout: 5 * time.Second,
		Addr:        global.GetListen(),
		Handler:     router,
	}

	// 如果key文件不存在那么就自动生成
	keyfile := filepath.Join("keys", "server.key")
	pemfile := filepath.Join("keys", "server.pem")
	_, err1 := os.Stat(keyfile)
	_, err2 := os.Stat(pemfile)
	if global.GetKey() == "" || global.GetPem() == "" && os.IsNotExist(err1) || os.IsNotExist(err2) {

		internal.CreateTLS()
		golog.Info(global.GetListen())
		on := "listen on " + global.GetListen() + " over https"
		golog.Info(on)
		if err := svc.ListenAndServeTLS(filepath.Join("keys", "server.pem"), filepath.Join("keys", "server.key")); err != nil {
			golog.Fatal(err)
		}
	}
	golog.Info("listen on " + global.GetListen() + " over http")
	if err := svc.ListenAndServeTLS(global.GetPem(), global.GetKey()); err != nil {
		log.Fatal(err)
	}
}
