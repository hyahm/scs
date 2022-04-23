package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/hyahm/scs/api"
	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/pkg/config/alert"
	"github.com/hyahm/scs/pkg/message"

	"github.com/hyahm/golog"
)

var configfile string
var showversion bool

func main() {
	defer golog.Sync()
	// golog.Level = golog.DEBUG
	golog.Format = "{{ .Ctime }} - [{{ .Level }}]- {{.Msg}}"
	// 异步获取ip，防止阻塞
	go message.GetIp()
	// 设置limit值
	internal.Setrlimit()
	flag.BoolVar(&showversion, "v", false, "get scs server version")
	flag.StringVar(&configfile, "f", "scs.yaml", "set config file")
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
			controller.WaitKillAllServer()
			os.Exit(1)
		}
	}()
	// 自动清除全局报警器的值
	go alert.CleanAlert()
	golog.Info("config file path: ", configfile)

	controller.Start(configfile)
	api.HttpServer()

}
