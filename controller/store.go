package controller

// 操作
import (
	"sync"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/config"
	"github.com/hyahm/scs/internal/config/scripts"
	"github.com/hyahm/scs/internal/server"
)

var cfg *config.Config

// 保存的servers
var servers map[string]*server.Server

// 保存 server的index还有哪些
var serverIndex map[string]map[int]struct{}

// 保存的scripts
var ss map[string]*scripts.Script

var mu sync.RWMutex

func init() {
	mu = sync.RWMutex{}
	servers = make(map[string]*server.Server)
	ss = make(map[string]*scripts.Script)
	serverIndex = make(map[string]map[int]struct{})
}

// 刚启动
func Start(filename string) {
	cfg, err := config.Start(filename)
	if err != nil {
		// 第一次报错直接退出
		golog.Fatal(err)
	}
	if cfg.SC == nil {
		cfg.SC = make([]*scripts.Script, 0)
	}
	for index := range cfg.SC {

		ss[cfg.SC[index].Name] = cfg.SC[index]
		ss[cfg.SC[index].Name].EnvLocker = &sync.RWMutex{}
		replicate := ss[cfg.SC[index].Name].Replicate
		if replicate == 0 {
			replicate = 1
		}
		newServerIndex := make(map[int]struct{})
		serverIndex[cfg.SC[index].Name] = newServerIndex

		makeReplicateServerAndStart(ss[cfg.SC[index].Name], replicate)
	}

}

func GetServerBySubname(subname string) (*server.Server, bool) {
	mu.RLock()
	defer mu.RUnlock()
	v, ok := servers[subname]
	return v, ok
}
