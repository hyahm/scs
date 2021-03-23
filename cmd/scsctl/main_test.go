package main

import (
	"testing"

	"github.com/hyahm/scs/client/cliconfig"
)

func BenchmarkMain(t *testing.B) {
	// 如果不是windows系统
	// 配置文件就放在 /etc/ 下面
	cliconfig.NewClientConfig()

	// command.Execute()
}

func TestMain(t *testing.T) {
	// 如果不是windows系统
	// 配置文件就放在 /etc/ 下面
	cliconfig.NewClientConfig()

	// command.Execute()
}
