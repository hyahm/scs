package controller

func WaitKillAllServer() {
	store.mu.Lock()
	defer store.mu.Unlock()
	// ss.ScriptLocker.RLock()
	// defer ss.ScriptLocker.RUnlock()
	for _, svc := range store.servers {
		svc.Kill()
		// replicate := s.Replicate
		// if replicate == 0 {
		// 	replicate = 1
		// }
		// for i := 0; i < replicate; i++ {
		// 	subname := subname.NewSubname(s.Name, i)
		// 	servers[subname.String()].Kill()
		// }
	}
}

// 同步杀掉
// func WaitKillScript(s *scripts.Script) {
// 	// ss.ServerLocker.RLock()
// 	// defer ss.ServerLocker.RUnlock()
// 	// 禁用 script 所在的所有server
// 	// mu.RLock()
// 	// defer mu.RUnlock()
// 	// 禁用 script 所在的所有server
// 	replicate := s.Replicate
// 	if replicate == 0 {
// 		replicate = 1
// 	}
// 	for i := 0; i < replicate; i++ {
// 		subname := subname.NewSubname(s.Name, i)
// 		servers[subname.String()].Kill()
// 	}
// }
