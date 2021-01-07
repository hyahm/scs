package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"scs/global"
	"scs/httpserver/handle"
	"scs/httpserver/midware"
	"scs/internal"
	"scs/public"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func GetExecTime(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	handle(w, r)
	golog.Infof("url: %s -- addr: %s -- method: %s -- exectime: %f\n", r.URL.Path, r.RemoteAddr, r.Method, time.Since(start).Seconds())
}

// var dir := "key"
func HttpServer() {
	router := xmux.NewRouter()
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.SetHeader("Content-Type", "application/x-www-form-urlencoded,application/json; charset=UTF-8")
	router.SetHeader("Access-Control-Allow-Headers", "Content-Type")
	router.AddModule(midware.CheckToken)
	// 增加请求时间
	// router.MiddleWare(GetExecTime)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
		return
	})
	router.Post("/status/{pname}/{name}", handle.Status)

	router.Post("/status", handle.AllStatus)

	router.Post("/status/{pname}", handle.StatusPname)

	router.Post("/start/{pname}/{name}", handle.Start)
	router.Post("/start", handle.StartAll)
	router.Post("/start/{pname}", handle.StartPname)

	router.Post("/stop/{pname}/{name}", handle.Stop)
	router.Post("/stop/{pname}", handle.StopPname)
	router.Post("/stop", handle.StopAll)

	router.Post("/restart/{pname}/{name}", handle.Restart)
	router.Post("/restart/{pname}", handle.RestartPname)
	router.Post("/restart", handle.RestartAll)

	router.Post("/canstop/{name}", handle.CanStop)
	router.Post("/cannotstop/{name}", handle.CanNotStop)

	router.Post("/-/reload", handle.Reload)

	router.Post("/log/{name}", handle.Log)
	router.Post("/get/repo", handle.GetRepo)

	router.Post("/kill/{pname}", handle.KillPname)
	router.Post("/kill/{pname}/{name}", handle.Kill)
	router.ShowApi("/docs")
	// router.Post("/env/{pname}/{name}", handle.GetEnv)
	router.Post("/env/{name}", handle.GetEnvName)

	router.Post("/install/{name}", handle.InstallPackage)
	router.Post("/install", handle.InstallScript)
	router.Post("/set/alert", handle.Alert)
	router.Post("/get/alert", handle.Alert)
	// 监测点
	router.Get("/probe", handle.Probe).DelModule(midware.CheckToken)

	// router.Get("/version/{pname}/{name}", handle.Version)
	router.Post("/script", handle.AddScript).Bind(&internal.Script{}).AddModule(midware.Unmarshal)
	router.Post("/delete/{pname}", handle.DelScript)

	svc := &http.Server{
		ReadTimeout: 5 * time.Second,
		Addr:        global.Listen,
		Handler:     router,
	}

	public.CreateTLS()
	fmt.Println("listen on " + global.Listen + " over https")
	if err := svc.ListenAndServeTLS(filepath.Join("keys", "server.pem"), filepath.Join("keys", "server.key")); err != nil {
		log.Fatal(err)
	}
}
