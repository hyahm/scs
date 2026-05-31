package store

import "sync"

type indexStore struct {
	mu          sync.RWMutex
	serverIndex map[string]map[int]struct{}
}

func NewIndexStore() *indexStore {
	return &indexStore{
		mu:          sync.RWMutex{},
		serverIndex: make(map[string]map[int]struct{}),
	}
}

func (s *indexStore) Set(pname string, i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.serverIndex[pname]; !ok {
		s.serverIndex[pname] = make(map[int]struct{})
	}
	s.serverIndex[pname][i] = struct{}{}
}

func (s *indexStore) Delete(pname string, i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.serverIndex[pname]; !ok {
		return
	}
	delete(s.serverIndex[pname], i)
}

func (s *indexStore) Get(pname string) []int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	indexs := make([]int, 0)
	if _, ok := s.serverIndex[pname]; !ok {
		return indexs
	}
	for index := range s.serverIndex[pname] {
		indexs = append(indexs, index)
	}
	return indexs
}

func (s *indexStore) Len(pname string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.serverIndex[pname]; !ok {
		return 0
	}
	return len(s.serverIndex[pname])
}

func (s *indexStore) Has(pname string, i int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.serverIndex[pname]; !ok {
		return false
	}
	_, ok := s.serverIndex[pname][i]
	return ok
}
