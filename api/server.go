package api

import (
	"net/http"
	"time"

	"github.com/hyahm/scs/api/handle"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/pkg"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

// 为了兼容之前版本，无视解析失败的问题
func unmarshalError(w http.ResponseWriter, r *http.Request, err error) bool {
	return true
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

	router.UnmarshalError = unmarshalError
	router.Post("/probe", handle.Probe)

	router.AddGroup(AdminHandle())

	router.SetAddr(internal.GetListen()).SetTimeout(time.Second * 5)
	if !internal.GetEnableTLS() {
		err := router.Run()
		if err != nil {
			golog.Error(err)
		}
		// os.Exit(1)
		return
	}

	if internal.GetKey() == "" || internal.GetCert() == "" {
		panic("use tls. but key or cert not specified")
	}

	golog.Info("listen on " + internal.GetListen() + " over https")
	err := router.RunTLS(internal.GetCert(), internal.GetKey())
	if err != nil {
		golog.Error(err)
	}
	// if err := svc.ListenAndServeTLS(filepath.Join("keys", "server.pem"), filepath.Join("keys", "server.key")); err != nil {
	// 	golog.Fatal(err)
	// }
}
