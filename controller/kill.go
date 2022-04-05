package controller

func WaitKillAllServer() {
	store.mu.Lock()
	defer store.mu.Unlock()
	for _, svc := range store.servers {
		svc.Kill()
	}
}
