package main

import (
	"github.com/hyahm/golog"
	"github.com/hyahm/scs/client/cliconfig"
	"github.com/hyahm/scs/client/command"
)

func main() {
	// 如果不是windows系统
	// 配置文件就放在 /etc/ 下面
	// cliconfig.NewClientConfig()
	defer golog.Sync()
	cliconfig.ReadConfig()
	command.Execute()
}
