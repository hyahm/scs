package controller

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/scs/pkg/config/scripts"
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
	store.mu.Lock()
	defer store.mu.Unlock()
	for _, script := range cfg.SC {

		// 如果没设置token， 默认生成一个脚本的token
		if script.Token == "" {
			script.Token = pkg.RandomToken()
		}
		// 将scripts填充到store中
		store.ss[script.Name] = script
		replicate := script.Replicate
		if replicate == 0 {
			replicate = 1
		}

		if store.serverIndex[script.Name] == nil {
			store.serverIndex[script.Name] = make(map[int]struct{})
		}
		// serverIndex[cfg.SC[index].Name] = newServerIndex
		// 生成server， 填充进来
		// 生成环境变量, 填充到script.tempenv里面
		script.MakeEnv()
		// 假设设置的端口是可用的
		availablePort := script.Port
		for i := 0; i < replicate; i++ {
			subname := fmt.Sprintf("%s_%d", script.Name, i)
			store.servers[subname] = makeServer(script, replicate, availablePort, i)
			if script.Disable {
				// 如果是禁用的 ，那么不用生成多个副本，直接执行下一个script
				break
			}

			store.servers[subname].Start()
		}
		// makeReplicateServerAndStart(ss[cfg.SC[index].Name], replicate)
	}

}

func makeServer(script *scripts.Script, replicate, availablePort, i int) *server.Server {
	svc := &server.Server{}
	subname := fmt.Sprintf("%s_%d", script.Name, i)
	// 将索引添加serverIndex里面
	store.serverIndex[script.Name][i] = struct{}{}
	// 将环境变量填充到server中
	env := make(map[string]string)
	for k, v := range script.TempEnv {
		env[k] = v
	}
	if script.Port > 0 {
		// 顺序拿到可用端口
		availablePort = pkg.GetAvailablePort(availablePort)
		env["PORT"] = strconv.Itoa(availablePort)
		svc = script.Add(availablePort, replicate, i, subname)
		availablePort++
	} else {
		env["PORT"] = "0"
		svc = script.Add(0, replicate, i, subname)
	}
	env["OS"] = runtime.GOOS
	// 格式化 SCS_TPL 开头的环境变量
	for k := range env {
		if len(k) > 8 && k[:7] == "SCS_TPL" {
			env[k] = internal.Format(env[k], env)
		}
	}
	svc.Env = env
	svc.Port = availablePort
	return svc

}
