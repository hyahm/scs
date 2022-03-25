package main

import (
	"os"
	"path/filepath"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/client"
	"github.com/hyahm/scs/client/command"
)

func main() {
	defer golog.Sync()
	// 如果不是windows系统
	// 配置文件就放在 /etc/ 下面
	// cliconfig.NewClientConfig()
	root, err := os.UserHomeDir()
	if err != nil {
		// 找不到就报错
		panic(err)
	}
	configfile := filepath.Join(root, ".scsctl.yaml")

	client.ReadClientConfig(configfile)
	command.Execute()
}
