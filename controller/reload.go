package controller

import (
	"fmt"
	"sync/atomic"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/scs/pkg/config/scripts"
)

func getTempScript(temp map[string]struct{}) {
	for name := range store.Store.GetAllScriptMap() {
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
		for index := range store.Store.GetScriptIndex(name) {
			subname := fmt.Sprintf("%s_%d", name, index)
			svc, ok := store.Store.GetServerByName(subname)
			if !ok {
				golog.Error(pkg.ErrBugMsg)
				continue
			}
			atomic.AddInt64(&global.CanReload, 1)
			go Remove(svc, false)
		}
	}
	return nil
}

// 添加script并启动

func AddScript(s *scripts.Script) {
	if s.Token == "" {
		s.Token = pkg.RandomToken()
	}
	if s.Role == "" {
		s.Role = scripts.ScriptRole
	}
	// 将scripts填充到store中
	store.Store.SetScript(s)
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	// 初始化脚本的副本数

	// 生成环境变量, 填充到script.tempenv里面

	// 假设设置的端口是可用的
	availablePort := s.Port
	for i := 0; i < replicate; i++ {
		subname := fmt.Sprintf("%s_%d", s.Name, i)
		store.Store.InitServer(i, replicate, s.Name, subname)
		store.Store.SetScriptIndex(s.Name, i)
		svc, _ := store.Store.GetServerByName(subname)
		availablePort = svc.MakeServer(s, availablePort)
		availablePort++
		if s.Disable {
			// 如果是禁用的 ，那么不用生成多个副本，直接执行下一个script
			return
		}

		svc.Start()
	}
}

func UpdateScript(s *scripts.Script, update bool) {
	script, ok := store.Store.GetScriptByName(s.Name)
	if !ok {
		return
	}
	oldReplicate := script.Replicate
	if oldReplicate == 0 {
		oldReplicate = 1
	}
	if s.Replicate == 1 {
		script.Replicate = 0
	}
	newReplicate := s.Replicate
	if newReplicate == 0 {
		newReplicate = 1
	}

	script = s
	availablePort := s.Port
	for i := 0; i < newReplicate; i++ {
		if !store.Store.HaveServerByIndex(s.Name, i) {
			subname := fmt.Sprintf("%s_%d", s.Name, i)
			store.Store.InitServer(i, newReplicate, s.Name, subname)
			store.Store.SetScriptIndex(s.Name, i)
			svc, _ := store.Store.GetServerByName(subname)
			availablePort = svc.MakeServer(s, availablePort)
			availablePort++
			if script.Disable {
				// 如果是禁用的 ，那么不用生成多个副本，直接执行下一个script
				return
			}
			svc.Start()
		}
	}
	// 删除多余的
	for i := newReplicate; i < oldReplicate; i++ {
		subname := fmt.Sprintf("%s_%d", s.Name, i)
		svc, ok := store.Store.GetServerByName(subname)
		if !ok {
			golog.Error(pkg.ErrBugMsg)
			continue
		}
		golog.Info("remove " + script.Name + fmt.Sprintf("_%d", i))
		atomic.AddInt64(&global.CanReload, 1)
		go Remove(svc, false)
	}

}

// todo:
func reloadScripts(s *scripts.Script, update bool) {
	// script: 配置文件新读取出来的
	// 处理存在的
	script, ok := store.Store.GetScriptByName(s.Name)
	// 对比启动的副本
	if !ok {
		// 如果不存在，说明要新增
		AddScript(s)
		return
	}
	oldReplicate := script.Replicate
	if oldReplicate == 0 {
		oldReplicate = 1
	}
	if s.Replicate == 1 {
		s.Replicate = 0
	}
	newReplicate := s.Replicate
	if newReplicate == 0 {
		newReplicate = 1
	}

	// 对比脚本是否修改
	golog.Error(oldReplicate)
	golog.Error(newReplicate)
	golog.Error(script.Command)
	golog.Error(s.Command)
	if oldReplicate == newReplicate {
		if !scripts.EqualScript(s, script) {
			golog.Info(s.Name)
			store.Store.SetScript(s)
		}

		// 如果一样的名字， 副本数一样的就直接跳过
		return
	}
	if oldReplicate > newReplicate {
		// 如果大于的话， 那么就删除多余的
		for i := newReplicate; i < oldReplicate; i++ {
			atomic.AddInt64(&global.CanReload, 1)
			subname := fmt.Sprintf("%s_%d", s.Name, i)
			golog.Info("remove " + s.Name + fmt.Sprintf("_%d", i))
			svc, ok := store.Store.GetServerByName(subname)
			if !ok {
				golog.Error(pkg.ErrBugMsg)
				continue
			}
			go Remove(svc, update)
		}
	} else {
		// 小于的话，就增加
		availablePort := s.Port
		for i := oldReplicate; i < newReplicate; i++ {
			subname := fmt.Sprintf("%s_%d", s.Name, i)
			store.Store.InitServer(i, newReplicate, s.Name, subname)
			store.Store.SetScriptIndex(s.Name, i)
			svc, _ := store.Store.GetServerByName(subname)
			availablePort = svc.MakeServer(s, availablePort)
			availablePort++
			if s.Disable {
				// 如果是禁用的 ，那么不用生成多个副本，直接执行下一个script
				return
			}
			svc.Start()
		}
	}

}
