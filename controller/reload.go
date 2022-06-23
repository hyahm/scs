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
		reloadScripts(cfg.SC[index])
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

// 新增script并启动
func AddScript(s *scripts.Script) {
	if s.ScriptToken == "" {
		s.ScriptToken = pkg.RandomToken()
	}
	// 将scripts填充到store中
	store.Store.SetScript(s)
	// 初始化脚本的副本数
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	s.MakeTempEnv()
	// 假设设置的端口是可用的
	// 对于每个script 都生成对应的
	availablePort := s.Port
	for i := 0; i < replicate; i++ {
		subname := fmt.Sprintf("%s_%d", s.Name, i)
		svc := store.Store.InitServer(i, s.Name, subname)
		store.Store.SetScriptIndex(s.Name, i)
		svc.Port = availablePort
		svc.MakeServer(s)
		availablePort = svc.Port + 1
		if s.Disable {
			// 如果是禁用的 ，那么不用生成多个副本
			return
		}

		svc.Start()
	}
}

// 脚本更新操作
func UpdateScriptApi(s *scripts.Script) {
	// 既然是更新操作，那么这个必定存在
	script, _ := store.Store.GetScriptByName(s.Name)

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
			atomic.AddInt64(&global.CanReload, 1)
			go Remove(svc, false)
		}
	} else {
		// 小于的话，就增加
		availablePort := s.Port
		for i := oldReplicate; i < newReplicate; i++ {
			subname := fmt.Sprintf("%s_%d", s.Name, i)
			svc := store.Store.InitServer(i, s.Name, subname)
			store.Store.SetScriptIndex(s.Name, i)
			svc.Port = availablePort
			svc.MakeServer(s)
			availablePort = svc.Port + 1
			if s.Disable {
				// 如果是禁用的 ，那么不用生成多个副本，直接执行下一个script
				return
			}
			svc.Start()
		}
	}

}

// 配置文件直接reload
func reloadScripts(s *scripts.Script) {
	// script: 配置文件新读取出来的
	// 处理存在的
	script, ok := store.Store.GetScriptByName(s.Name)
	// 对比启动的副本
	if !ok {
		// 如果不存在，说明要新增
		AddScript(s)
		return
	}
	UpdateScriptApi(script)
}
