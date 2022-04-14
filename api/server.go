package api

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/hyahm/scs/api/handle"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

// var dir := "key"
func HttpServer() {
	router := xmux.NewRouter()
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.SetHeader("Content-Type", "application/x-www-form-urlencoded,application/json; charset=UTF-8")
	router.SetHeader("Access-Control-Allow-Headers", "Content-Type")

	router.Post("/probe", handle.Probe)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	router.AddGroup(AdminHandle())
	router.AddGroup(scriptHandle())
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
	if err := svc.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
