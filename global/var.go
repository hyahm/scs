package global

import (
	"sync"
	"time"
)

// 配置文件重载
var cr *canreloads

type canreloads struct {
	busy bool
	msg  string
	mu   sync.Mutex
}

func init() {
	cr = &canreloads{
		mu: sync.Mutex{},
	}
}

// 是否可以重载
func IsCanReload() (string, bool) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	return cr.msg, !cr.busy
}

// 设置重载， 如果false 就是加载中，等待设置完成
func SetReLoading(msg string) (string, bool) {
	_, ok := IsCanReload()
	if ok {
		cr.busy = true
		cr.msg = msg
	}
	return cr.msg, ok
}

func SetCanReLoad() {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.busy = false
	cr.msg = ""
}

type ConfigStore struct {
	Listen             string
	Token              string
	LogCount           int
	ReadTime           time.Duration
	IgNoreToken        []string
	Monitored          []string
	EnableTLS          bool
	Key                string
	Cert               string
	ContinuityInterval time.Duration
	LogDir             string
	CleanLog           int
}

var CS ConfigStore
var ProxyHeader string

const VERSION = "v3.8.5"

// const FORMAT = "{{ .Ctime }} - [{{ .Level }}] - {{.Hostname}} - {{ .Msg }}"
