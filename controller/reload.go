package controller

import (
	"fmt"
	"sync/atomic"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/scs/pkg/config/scripts/subname"
)

func getTempScript(temp map[string]struct{}) {
	store.mu.Lock()
	defer store.mu.Unlock()
	for name := range store.ss {
		temp[name] = struct{}{}
	}
}

func Fmt() error {
	c, err := config.ReadConfig("")
	if err != nil {
		golog.Error(err)
		// 第一次报错直接退出
		return err
	}
	// 配置文件是对的， 那么直接写进配置文件
	return c.WriteConfig(true)
}

/*
重载：  备份旧的scripts的replicate

如果有多余的副本或scripts就删除, 新的scripts将会启动， 自动扩缩容
更新所有store.ss

*/
func Reload() error {
	c, err := config.ReadConfig("")
	if err != nil {
		// 第一次报错直接退出
		return err
	}
	// 配置文件是对的， 那么直接写进配置文件， 后面所有的操作都取消更新配置文件
	cfg = c
	err = cfg.WriteConfig(true)
	if err != nil {
		// 写进配置文件
		return err
	}
	// 取出之前的scripts
	temp := make(map[string]struct{})
	// 备份旧的scripts
	getTempScript(temp)
	for index := range cfg.SC {

		// 删除之前存在的name
		delete(temp, cfg.SC[index].Name)
		// 查看副本是不是对的， 不会对存在的脚本有影响
		reloadScripts(cfg.SC[index], false)
	}

	// 删除已删除的 script
	for name := range temp {
		if _, ok := store.ss[name]; ok {
			golog.Info("remove ", name)
			replicate := store.ss[name].Replicate
			if replicate == 0 {
				replicate = 1
			}

			for i := 0; i < replicate; i++ {
				subname := subname.NewSubname(name, i)
				atomic.AddInt64(&global.CanReload, 1)
				go Remove(store.servers[subname.String()], false)
			}

		}
	}
	return nil
}

// 添加script并启动

func AddScript(script *scripts.Script) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if script.Token == "" {
		script.Token = pkg.RandomToken()
	}
	// 将scripts填充到store中
	store.ss[script.Name] = script
	replicate := script.Replicate
	if replicate == 0 {
		replicate = 1
	}
	// 初始化脚本的副本数
	if _, ok := store.serverIndex[script.Name]; !ok {
		store.serverIndex[script.Name] = make(map[int]struct{})
	}
	// 生成环境变量, 填充到script.tempenv里面

	// 假设设置的端口是可用的
	availablePort := script.Port
	for i := 0; i < replicate; i++ {
		subname := fmt.Sprintf("%s_%d", script.Name, i)
		store.servers[subname] = &server.Server{
			Index:     i,
			Replicate: replicate,
			SubName:   subname,
			Name:      script.Name,
		}
		store.serverIndex[script.Name][i] = struct{}{}
		availablePort = store.servers[subname].MakeServer(script, availablePort)
		availablePort++
		if script.Disable {
			// 如果是禁用的 ，那么不用生成多个副本，直接执行下一个script
			return
		}

		store.servers[subname].Start()
	}
}

func UpdateScript(script *scripts.Script, update bool) {
	store.mu.Lock()

	oldReplicate := store.ss[script.Name].Replicate
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

	store.ss[script.Name] = script
	availablePort := script.Port
	for i := 0; i < newReplicate; i++ {
		if _, ok := store.serverIndex[script.Name][i]; !ok {
			subname := fmt.Sprintf("%s_%d", script.Name, i)
			store.servers[subname] = &server.Server{
				Index:     i,
				Replicate: newReplicate,
				SubName:   subname,
				Name:      script.Name,
			}
			store.serverIndex[script.Name][i] = struct{}{}
			availablePort = store.servers[subname].MakeServer(script, availablePort)
			availablePort++
			if script.Disable {
				// 如果是禁用的 ，那么不用生成多个副本，直接执行下一个script
				store.mu.Unlock()
				return
			}
			store.servers[subname].Start()
		}
	}
	store.mu.Unlock()
	// 删除多余的
	for i := newReplicate; i < oldReplicate; i++ {
		golog.Info("remove " + script.Name + fmt.Sprintf("_%d", i))
		atomic.AddInt64(&global.CanReload, 1)
		go Remove(store.servers[subname.NewSubname(script.Name, i).String()], false)
	}

}

// todo:
func reloadScripts(script *scripts.Script, update bool) {
	// script: 配置文件新读取出来的
	// 处理存在的

	// 对比启动的副本
	if _, ok := store.ss[script.Name]; !ok {
		// 如果不存在，说明要新增
		AddScript(script)
		return
	}
	oldReplicate := store.ss[script.Name].Replicate
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
	store.ss[script.Name] = script
	if oldReplicate == newReplicate {
		// 如果一样的名字， 副本数一样的就直接跳过
		return
	}
	if oldReplicate > newReplicate {
		// 如果大于的话， 那么就删除多余的
		for i := newReplicate; i < oldReplicate; i++ {
			atomic.AddInt64(&global.CanReload, 1)
			golog.Info("remove " + script.Name + fmt.Sprintf("_%d", i))
			go Remove(store.servers[subname.NewSubname(script.Name, i).String()], update)
		}
	} else {
		// 小于的话，就增加
		availablePort := script.Port
		for i := oldReplicate; i < newReplicate; i++ {
			subname := fmt.Sprintf("%s_%d", script.Name, i)
			store.servers[subname] = &server.Server{
				Index:     i,
				Replicate: newReplicate,
				SubName:   subname,
				Name:      script.Name,
			}
			store.serverIndex[script.Name][i] = struct{}{}
			availablePort = store.servers[subname].MakeServer(script, availablePort)
			availablePort++
			if script.Disable {
				// 如果是禁用的 ，那么不用生成多个副本，直接执行下一个script
				return
			}

			store.servers[subname].Start()
		}
	}

}
