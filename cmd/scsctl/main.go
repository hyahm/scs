package main

import (
	"scs/client/cliconfig"
	"scs/client/command"
)

func main() {
	// 如果不是windows系统
	// 配置文件就放在 /etc/ 下面
	cliconfig.ReadConfig()

	command.Execute()

}
