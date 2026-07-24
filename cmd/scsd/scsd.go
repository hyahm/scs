package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hyahm/scs/api"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/scs/pkg/message"

	"github.com/hyahm/golog"
)

var showversion bool

func main() {
	defer golog.Sync()
	// golog.Level = golog.DEBUG
	// golog.Format = "{{ .Ctime }} - [{{ .Level }}]- {{.Msg}}"
	// 异步获取ip，防止阻塞
	go message.GetIp()
	// 设置limit值
	internal.Setrlimit()
	flag.BoolVar(&showversion, "v", false, "get scs server version")
	flag.StringVar(&internal.ConfigFile, "f", "scs.yaml", "set config file")
	flag.Parse()
	if showversion {
		fmt.Println(global.VERSION)
		return
	}
	single := make(chan os.Signal, 1)
	signal.Notify(single, os.Interrupt, syscall.SIGTERM, syscall.SIGPIPE)
	go func() {
		for range single {
			// 确保删除了server
			fmt.Println("waiting stop all")
			// controller.WaitKillAllServer()
			os.Exit(1)
		}

	}()
	golog.Info("config file path: ", internal.ConfigFile)
	err := internal.ReadConfig()
	if err != nil {
		// 第一次报错直接退出
		golog.Fatal(err)
	}
	golog.SetLevel(golog.DEBUG)
	// if internal.GetDebug() {
	// 	golog.SetLevel(golog.DEBUG)
	// } else {
	// 	golog.SetDir(internal.GetLogPath())
	// 	golog.InitLogger("scs.log", 0, true)
	// }

	// 自动清除全局报警器的值
	go config.CleanAlert()
	// 启动脚本
	cache.FirstStartAllScript()
	// 回写
	err = internal.WriteConfig(internal.GetConfig())
	if err != nil {
		// 第一次报错直接退出
		golog.Error(err)
	}
	api.HttpServer()

}
