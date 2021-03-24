package script

import "sync"

// 记录加载配置文件前使用的脚本服务
// 重新加载的时候， 有多余的就会删除

// 记录这次所有使用的, 这个先只添加就可以了， key 是subname， value是name
var delScript map[string]string // reload的时候，记录需要删除的
var isreloading bool
var slock sync.RWMutex

func Copy() {
	// delScript = make(map[string]string)
	delScript = SS.Copy()
	slock = sync.RWMutex{}
	// delScript = make(map[string]string)
	// for pname, v := range SS.Infos {
	// 	for name := range v {
	// 		delScript[name] = pname
	// 	}
	// }
}

func DelDelScript(key string) {
	slock.Lock()
	defer slock.Unlock()
	delete(delScript, key)
}

func RemoveUnUseScript() {
	slock.RLock()
	defer slock.RUnlock()
	// 停止无用的脚本， 并移除
	for subname, name := range delScript {
		SS.GetScriptFromPnameAndSubname(subname, name).Remove()
	}
	delScript = nil
}
