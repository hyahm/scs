package main

import (
	"flag"
	"fmt"
	"scs/alert"
	"scs/config"
	"scs/global"
	"scs/httpserver"
)

func main() {
	var cfg string
	var showversion bool
	flag.BoolVar(&showversion, "v", false, "get scs server version")
	flag.StringVar(&cfg, "f", "scs.yaml", "set config file")
	flag.Parse()
	if showversion {
		fmt.Println("version:", global.VERSION)
		return
	}
	// 报警器
	go alert.SendNetAlert()
	config.Start(cfg)

	httpserver.HttpServer()
	// 依次启动
}
