package controller

import (
	"fmt"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
)

// 启动存在的脚本
func StartExsitScript(name string) {
	for index := range store.Store.GetScriptIndex(name) {
		subname := fmt.Sprintf("%s_%d", name, index)
		svc, ok := store.Store.GetServerByName(subname)
		if ok {
			svc.Start()
			continue
		}
		golog.Error(pkg.ErrBugMsg)
	}
}

// 第一次启动
func FirstStartAllScript() {
	cfg, err := config.ReadConfig()
	if err != nil && cfg == nil {
		// 第一次报错直接退出
		golog.Fatal(err)
	}
	for _, script := range cfg.SC {
		// 如果没设置token， 默认生成一个脚本的token
		AddScript(script)
	}
}

// // 启动脚本, 也有可能是重载
// func StartAllScript() {
// 	// 先将配置文件填充到 store
// 	for _, script := range cfg.SC {
// 		// 如果没设置token， 默认生成一个脚本的token
// 		AddScript(script)
// 	}
// }
