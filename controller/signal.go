package controller

// 保存是否可停止的状态

// type SNS struct {
// 	timeout int64           // 保存超时时间
// 	ctx     context.Context // 上下文管理
// 	cancel  context.CancelFunc
// }

// var StopSignal map[string]SNS
// var signalMu sync.RWMutex
// var signal chan string

// func init() {
// 	StopSignal = make(map[string]SNS)
// 	signalMu = sync.RWMutex{}
// 	signal = make(chan string, 1000)
// }

// // 启动主线程的时候就启动这个
// func StartSignalTime() {
// 	for {
// 		select {
// 		case name := <-signal:
// 			// 把这个key删掉，并且结束对应的goroutine
// 			removeSignal(name)
// 		}
// 	}
// }

// func removeSignal(name string) {
// 	signalMu.Lock()
// 	defer signalMu.Unlock()
// 	// 为了避免短暂的操作， 还是先判断下是否存在这个key
// 	if v, ok := StopSignal[name]; ok {
// 		// 先停止goroutine
// 		v.cancel()
// 	}

// }

// // 启动一个不能停止的信号
// func StartCanNotStop(name string, timeout int64) {
// 	// name: name
// 	// timeout: 超时时间， 超过这个时间没收到停止信号就发送报警
// 	if timeout == 0 {
// 		return
// 	}
// 	ctx, cancel := context.WithCancel(context.Background())
// 	sns := SNS{timeout: timeout, ctx: ctx, cancel: cancel}
// 	signalMu.Lock()
// 	defer signalMu.Unlock()

// 	if _, ok := StopSignal[name]; ok {
// 		// 如果存在这个key就是重复提交
// 		return
// 	}

// 	StopSignal[name] = sns
// 	go func(sns SNS) {
// 		// 开一个协程
// 		select {
// 		case <-time.After(time.Second * time.Duration(sns.timeout)):
// 			//  超时了发报警
// 			// 还是要判断一下，这个是不是运行状态
// 			fmt.Println("超时了")
// 		case <-sns.ctx.Done():
// 			// 然后删掉这个key
// 			delete(StopSignal, name)
// 		}
// 	}(sns)

// }

// // 启动一个可用停止的信号
// func StartCanStop(name string) {
// 	signalMu.RLock()
// 	defer signalMu.RUnlock()
// 	if _, ok := StopSignal[name]; ok {
// 		signal <- name
// 	}

// }
