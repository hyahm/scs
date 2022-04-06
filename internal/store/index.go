package store

func (s *store) SetScriptIndex(pname string, i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.serverIndex[pname]; !ok {
		s.serverIndex[pname] = make(map[int]struct{})
	}
	s.serverIndex[pname][i] = struct{}{}
}

func (s *store) DeleteScriptIndex(pname string, i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.serverIndex[pname]; !ok {
		return
	}
	delete(s.serverIndex[pname], i)
}

func (s *store) GetScriptIndex(pname string) []int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	indexs := make([]int, 0)
	for index := range s.serverIndex[pname] {
		indexs = append(indexs, index)
	}
	return indexs
}

func (s *store) GetScriptLength(pname string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.serverIndex[pname]; !ok {
		return 0
	}
	return len(s.serverIndex[pname])
}

// 判断是否存在这个script
func (s *store) HaveServerByIndex(pname string, i int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.serverIndex[pname]; !ok {
		return false
	}

	_, ok := s.serverIndex[pname][i]
	return ok
}
