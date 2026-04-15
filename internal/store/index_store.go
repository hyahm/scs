package store

import "sync"

type IndexStore struct {
	mu          sync.RWMutex
	serverIndex map[string]map[int]struct{}
}

func NewIndexStore() *IndexStore {
	return &IndexStore{
		serverIndex: make(map[string]map[int]struct{}),
	}
}

func (s *IndexStore) Set(pname string, i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.serverIndex[pname]; !ok {
		s.serverIndex[pname] = make(map[int]struct{})
	}
	s.serverIndex[pname][i] = struct{}{}
}

func (s *IndexStore) Delete(pname string, i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.serverIndex[pname]; !ok {
		return
	}
	delete(s.serverIndex[pname], i)
}

func (s *IndexStore) Get(pname string) []int {
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

func (s *IndexStore) Len(pname string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.serverIndex[pname]; !ok {
		return 0
	}
	return len(s.serverIndex[pname])
}

func (s *IndexStore) Has(pname string, i int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.serverIndex[pname]; !ok {
		return false
	}
	_, ok := s.serverIndex[pname][i]
	return ok
}
