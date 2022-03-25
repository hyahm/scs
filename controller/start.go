package controller

import "fmt"

// 启动存在的脚本
func StartExsitScript(name string) {
	mu.RLock()
	defer mu.RUnlock()
	replicate := servers[name+"_0"].Replicate
	for i := 0; i < replicate; i++ {
		servers[name+fmt.Sprintf("_%d", i)].Start()
	}

}
