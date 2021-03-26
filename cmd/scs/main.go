package main

import (
	"flag"
	"fmt"

	"github.com/hyahm/scs/alert"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/httpserver"
	"github.com/hyahm/scs/script"

	"github.com/hyahm/golog"
)

func main() {
	defer golog.Sync()
	go alert.GetIp()
	var cfg string
	var showversion bool
	flag.BoolVar(&showversion, "v", false, "get scs server version")
	flag.StringVar(&cfg, "f", "scs.yaml", "set config file")
	flag.Parse()
	if showversion {
		fmt.Println(global.VERSION)
		return
	}
	// 自动清除全局报警器的值
	go alert.SendNetAlert()
	script.Start(cfg)
	golog.Info("starting httpd")
	httpserver.HttpServer()
	// 依次启动
}
