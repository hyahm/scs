package pkg

import (
	"context"
	"sync"
)

type atomSignal struct {
	Cancel map[string]context.CancelFunc
	mu     *sync.RWMutex
}

type SignalRequest struct {
	Timeout            int64  `json:"timeout"`            // 超时时间， 默认s
	Restart            bool   `json:"restart"`            // 如果超时了是否重启
	Notice             bool   `json:"notice"`             // 如果超时了是否报警通知
	ContinuityInterval int    `json:"continuityInterval"` // 下次报警时间
	Parameter          string `json:"parameter"`          // 重启后的传参
}

var atom *atomSignal

func init() {
	atom = &atomSignal{
		Cancel: make(map[string]context.CancelFunc),
		mu:     &sync.RWMutex{},
	}
}

// 删除原子操作的超时处理
func DeleteAtomSignal(name string) {
	atom.mu.Lock()
	delete(atom.Cancel, name)
	atom.mu.Unlock()
}

// 设置信号
func SetAtomSignal(name string, cancel context.CancelFunc) {
	atom.mu.Lock()
	atom.Cancel[name] = cancel
	atom.mu.Unlock()
}

func CancelAtomSignal(name string) {
	atom.mu.Lock()
	if cancel, ok := atom.Cancel[name]; ok {
		cancel()
		delete(atom.Cancel, name)
	}
	atom.mu.Unlock()
}
