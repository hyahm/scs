package controller

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/config"
	"github.com/hyahm/scs/internal/config/scripts"
	"github.com/hyahm/scs/internal/config/scripts/subname"
)

func getTempScript(temp map[string]struct{}) {
	mu.RLock()
	defer mu.RUnlock()
	for name := range ss {
		temp[name] = struct{}{}
	}
}

func Reload() error {
	cfg, err := config.Start("")
	if err != nil {
		// 第一次报错直接退出
		return err
	}
	err = cfg.WriteConfigFile(true)
	if err != nil {
		// 写进配置文件
		return err
	}
	temp := make(map[string]struct{})
	getTempScript(temp)
	for index := range cfg.SC {

		// 	// 将数据填充至 SS, 返回是否存在此脚本
		if !config.CheckScriptNameRule(cfg.SC[index].Name) {
			golog.Error("script name must be a word, have been ignore: " + cfg.SC[index].Name)
			continue
		}
		// 删除之前存在的name
		delete(temp, cfg.SC[index].Name)
		// 	// 修改配置
		ReloadScripts(cfg.SC[index], false)
	}
	mu.Lock()
	defer mu.Unlock()
	// 删除已删除的 script
	for name := range temp {
		if _, ok := ss[name]; ok {
			replicate := ss[name].Replicate
			if replicate == 0 {
				replicate = 1
			}

			for i := 0; i < replicate; i++ {
				subname := subname.NewSubname(name, i)
				atomic.AddInt64(&global.CanReload, 1)
				go Remove(servers[subname.String()], false)
			}

		}
	}
	return nil
}

func AddScript(script *scripts.Script) {
	if ss == nil {
		ss = make(map[string]*scripts.Script)
	}
	ss[script.Name] = script
	replicate := script.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		makeReplicateServerAndStart(script, replicate)

	}
}

func UpdateScript(script *scripts.Script, update bool) {

	oldReplicate := ss[script.Name].Replicate
	if oldReplicate == 0 {
		oldReplicate = 1
	}
	if script.Replicate == 1 {
		script.Replicate = 0
	}
	newReplicate := script.Replicate
	if newReplicate == 0 {
		newReplicate = 1
	}

	if ss == nil {
		ss = make(map[string]*scripts.Script)
	}
	ss[script.Name] = script
	for i := 0; i < newReplicate; i++ {
		// subname := subname.NewSubname(script.Name, i)
		// servers[subname.String()].Start()
		go func(i int) {
			golog.Info("remove ", subname.NewSubname(script.Name, i).String())
			Remove(servers[subname.NewSubname(script.Name, i).String()], false)
			golog.Info("update")
			makeReplicateServerAndStart(ss[script.Name], newReplicate)
			golog.Info("update success")
		}(i)

	}

	// 删除多余的
	for i := newReplicate; i < oldReplicate; i++ {
		golog.Info("remove " + script.Name + fmt.Sprintf("_%d", i))
		Remove(servers[subname.NewSubname(script.Name, i).String()], false)
	}

}

func ReloadScripts(script *scripts.Script, update bool) {
	// script: 配置文件新读取出来的
	// 处理存在的
	if _, ok := ss[script.Name]; ok {
		// 对比启动的副本

		oldReplicate := ss[script.Name].Replicate
		if oldReplicate == 0 {
			oldReplicate = 1
		}
		if script.Replicate == 1 {
			script.Replicate = 0
		}
		newReplicate := script.Replicate
		if newReplicate == 0 {
			newReplicate = 1
		}
		// 对比脚本是否修改
		ss[script.Name].Replicate = newReplicate
		if !scripts.EqualScript(script, ss[script.Name]) {
			// 如果不一样， 那么 就需要重新启动服务

			ss[script.Name] = script
			ss[script.Name].EnvLocker = &sync.RWMutex{}
			for i := 0; i < newReplicate; i++ {
				// subname := subname.NewSubname(script.Name, i)
				// servers[subname.String()].Start()
				go func(i int) {
					Remove(servers[subname.NewSubname(script.Name, i).String()], update)
					makeReplicateServerAndStart(script, newReplicate)
				}(i)

			}

			// 删除多余的
			for i := newReplicate; i < oldReplicate; i++ {
				golog.Info("remove " + script.Name + fmt.Sprintf("_%d", i))
				Remove(servers[subname.NewSubname(script.Name, i).String()], update)
			}

			return
		}

		if oldReplicate == newReplicate {
			// 如果一样的名字， 副本数一样的就直接跳过
			return
		}
		if oldReplicate > newReplicate {
			// 如果大于的话， 那么就删除多余的
			for i := newReplicate; i < oldReplicate; i++ {
				delete(serverIndex[script.Name], i)
				atomic.AddInt64(&global.CanReload, 1)
				go Remove(servers[subname.NewSubname(script.Name, i).String()], update)
			}
		} else {
			makeReplicateServerAndStart(script, newReplicate)
		}

	} else {
		// 不存在的脚本直接启动即可
		ss[script.Name] = script
		replicate := script.Replicate
		if replicate == 0 {
			replicate = 1
		}
		makeReplicateServerAndStart(script, replicate)
	}
}
