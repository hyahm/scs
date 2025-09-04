package api

import (
	"net/http"
	"time"

	"github.com/hyahm/scs/api/handle"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/pkg"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

// 为了兼容之前版本，无视解析失败的问题
func unmarshalError(err error, w http.ResponseWriter, r *http.Request) bool {
	return false
}

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
	router.SetHeader("Access-Control-Max-Age", "1728000")
	router.Exit = exit

	router.AddGroup(xmux.Pprof())

	router.UnmarshalError = unmarshalError
	router.Post("/probe", handle.Probe)

	router.AddGroup(AdminHandle())

	router.SetAddr(global.CS.Listen).SetTimeout(time.Second * 5)
	if !global.CS.EnableTLS {
		golog.Info("listen on " + global.CS.Listen + " over http")
		golog.Fatal(router.Run())
	}

	if global.CS.Key == "" || global.CS.Cert == "" {
		panic("use tls. but key or cert not specified")
	}

	golog.Info("listen on " + global.CS.Listen + " over http")
	golog.Fatal(router.RunTLS(global.CS.Cert, global.CS.Key))
	// if err := svc.ListenAndServeTLS(filepath.Join("keys", "server.pem"), filepath.Join("keys", "server.key")); err != nil {
	// 	golog.Fatal(err)
	// }
}
