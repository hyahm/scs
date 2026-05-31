package api

import (
	"net/http"
	"time"

	"github.com/hyahm/scs/api/handle"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"

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

	router.UnmarshalError = unmarshalError
	router.Post("/probe", handle.Probe)

	router.AddGroup(AdminHandle())

	router.SetAddr(config.Cfg.Listen).SetTimeout(time.Second * 5)
	if !config.Cfg.EnableTLS {
		err := router.Run()
		if err != nil {
			golog.Error(err)
		}
		// os.Exit(1)
		return
	}

	if config.Cfg.Key == "" || config.Cfg.Cert == "" {
		panic("use tls. but key or cert not specified")
	}

	golog.Info("listen on " + config.Cfg.Listen + " over https")
	err := router.RunTLS(config.Cfg.Cert, config.Cfg.Key)
	if err != nil {
		golog.Error(err)
	}
	// if err := svc.ListenAndServeTLS(filepath.Join("keys", "server.pem"), filepath.Join("keys", "server.key")); err != nil {
	// 	golog.Fatal(err)
	// }
}
