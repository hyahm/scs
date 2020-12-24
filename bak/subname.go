package bak

import (
	"sync"
	"time"
)

// 记录加载配置文件前使用的脚本服务
// 重新加载的时候， 有多余的就会删除
var useScript map[string]string // 记录上次使用的, 这个先只添加就可以了
var delScript map[string]string // reload的时候，记录需要删除的
var uslocker sync.RWMutex
var Reloadlocker sync.RWMutex // reload 锁
var dellocker sync.RWMutex

func init() {
	useScript = make(map[string]string)
	uslocker = sync.RWMutex{}
	dellocker = sync.RWMutex{}
	Reloadlocker = sync.RWMutex{}
}

func Copy() {
	delScript = make(map[string]string)
	for k, v := range useScript {
		delScript[k] = v
	}
}

func DelDelScript(key string) {
	dellocker.Lock()
	delete(delScript, key)
	dellocker.Unlock()
}

func SetUseScript(key, value string) {
	uslocker.Lock()
	useScript[key] = value
	uslocker.Unlock()
}

func DelUseScript(key string) {
	uslocker.Lock()
	delete(useScript, key)
	uslocker.Unlock()
}

func StopUnUseScript() {
	// 停止无用的脚本， 并移除
	for subname, name := range delScript {
		// 停止服务, 先禁止always， 设置主动停止
		SS.Infos[name][subname].Status.Always = false
		time.Sleep(SS.Infos[name][subname].KillTime)
		SS.Infos[name][subname].exit = true

		if SS.Infos[name][subname].Status.Status != STOP {
			SS.Infos[name][subname].Stop()
		}

		// 阻塞服务停止才能继续
		for {
			if SS.Infos[name][subname].Status.Status == STOP {
				break
			}
			time.Sleep(SS.Infos[name][subname].KillTime)
		}
		// 从服务里面删除记录
		DelUseScript(subname)
		delete(SS.Infos[name], subname)
		if len(SS.Infos[name]) == 0 {
			// 清空name
			delete(SS.Infos, name)
		}
	}

}
