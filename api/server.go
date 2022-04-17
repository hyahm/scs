package api

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/hyahm/scs/api/handle"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/pkg"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

// var dir := "key"
func HttpServer() {
	response := &pkg.Response{
		Code:    200,
		Msg:     "ok",
		Version: global.VERSION,
	}
	router := xmux.NewRouter().BindResponse(response)
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.SetHeader("Content-Type", "application/x-www-form-urlencoded,application/json; charset=UTF-8")
	router.SetHeader("Access-Control-Allow-Headers", "Content-Type")
	router.Exit = exit
	router.Post("/probe", handle.Probe)

	router.AddGroup(AdminHandle())
	router.AddGroup(ScriptHandle())
	router.DebugAssignRoute("/status")
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
	}
	on := "listen on " + global.GetListen() + " over https"
	golog.Info(on)
	if err := svc.ListenAndServeTLS(filepath.Join("keys", "server.pem"), filepath.Join("keys", "server.key")); err != nil {
		golog.Fatal(err)
	}
}
