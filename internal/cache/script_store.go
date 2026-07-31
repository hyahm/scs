package cache

import (
	"sync"

	"github.com/hyahm/scs/pkg/config"
)

// 保存运行中的脚本
type Script struct {
	sync.RWMutex
	script map[string]config.Script // 保存脚本状态// 保存的  name 中间有_
}

var scriptInstance = &Script{
	script: make(map[string]config.Script),
	// index:   make(map[string]int),
}

func GetAllScript() map[string]config.Script {
	scriptInstance.RLock()
	defer scriptInstance.RUnlock()
	result := make(map[string]config.Script, len(scriptInstance.script))
	for k, v := range scriptInstance.script {
		result[k] = v
	}
	return result
}

func RemoveScript(name string) {
	scriptInstance.Lock()
	defer scriptInstance.Unlock()
	delete(scriptInstance.script, name)
}

func GetScript(name string) (config.Script, bool) {
	scriptInstance.RLock()
	defer scriptInstance.RUnlock()
	v, ok := scriptInstance.script[name]
	return v, ok
}

func SetScript(name string, script config.Script) {
	scriptInstance.Lock()
	defer scriptInstance.Unlock()
	scriptInstance.script[name] = script
}
