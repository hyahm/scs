package controller

import (
	"fmt"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/pkg/config"
)

// 启动存在的脚本
func StartExsitScript(name string) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for index := range store.serverIndex[name] {
		store.servers[fmt.Sprintf("%s_%d", name, index)].Start()
	}
}

// 启动服务
func Start(filename string) {
	config, err := config.ReadConfig(filename)
	if err != nil && cfg == nil {
		// 第一次报错直接退出
		golog.Fatal(err)
	}
	cfg = config
	startScripts()
}

// 启动脚本, 也有可能是重载
func startScripts() {
	// 先将配置文件填充到 store
	for _, script := range cfg.SC {
		// 如果没设置token， 默认生成一个脚本的token
		AddScript(script)
	}

}
