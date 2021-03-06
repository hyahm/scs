package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/hyahm/scs/alert"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/httpserver"
	"github.com/hyahm/scs/message"
	"github.com/hyahm/scs/server"

	"github.com/hyahm/golog"
)

func main() {
	defer golog.Sync()
	go message.GetIp()
	var cfg string
	var showversion bool
	flag.BoolVar(&showversion, "v", false, "get scs server version")
	flag.StringVar(&cfg, "f", "scs.yaml", "set config file")
	flag.Parse()
	if showversion {
		fmt.Println(global.VERSION)
		return
	}
	single := make(chan os.Signal, 1)
	signal.Notify(single, os.Interrupt, os.Kill)
	go func() {
		select {
		case <-single:
			// 确保删除了server
			fmt.Println("waiting stop all")
			server.WaitKillAllServer()
			os.Exit(1)
		}
	}()
	// 自动清除全局报警器的值
	go alert.SendNetAlert()
	server.Start(cfg)
	httpserver.HttpServer()

	// 依次启动
}
