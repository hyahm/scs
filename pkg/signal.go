package pkg

// type atomSignal struct {
// 	Cancel map[string]context.CancelFunc
// 	sync.RWMutex
// }

type SignalRequest struct {
	Timeout            int64  `json:"timeout"`            // 超时时间， 默认s
	Restart            bool   `json:"restart"`            // 如果超时了是否重启
	Notice             bool   `json:"notice"`             // 如果超时了是否报警通知
	ContinuityInterval int    `json:"continuityInterval"` // 下次报警时间
	Parameter          string `json:"parameter"`          // 重启后的传参
}

// var atom = &atomSignal{
// 	Cancel: make(map[string]context.CancelFunc),
// }

// // 删除原子操作的超时处理
// func DeleteAtomSignal1(name string) {
// 	atom.Lock()
// 	delete(atom.Cancel, name)
// 	atom.Unlock()
// }

// // 设置信号
// func SetAtomSignal1(name string, cancel context.CancelFunc) {
// 	atom.Lock()
// 	atom.Cancel[name] = cancel
// 	atom.Unlock()
// }

// func CancelAtomSignal1(name string) {
// 	atom.Lock()
// 	if cancel, ok := atom.Cancel[name]; ok {
// 		cancel()
// 		delete(atom.Cancel, name)
// 	}
// 	atom.Unlock()
// }
