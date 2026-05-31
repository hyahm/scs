package store

import (
	"sync"

	"github.com/hyahm/scs/pkg/config"
)

type scriptStore struct {
	mu sync.RWMutex
	ss map[string]config.Script
}

func NewScriptStore() *scriptStore {
	return &scriptStore{
		mu: sync.RWMutex{},
		ss: make(map[string]config.Script),
	}
}

func (s *scriptStore) Get(name string) (config.Script, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.ss[name]
	return v, ok
}

func (s *scriptStore) Set(script config.Script) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ss[script.Name] = script
}

func (s *scriptStore) Delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.ss, name)
}

func (s *scriptStore) GetAll() map[string]config.Script {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ss
}

func (s *scriptStore) GetMapFilterByName(names map[string]struct{}) map[string]config.Script {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ss := make(map[string]config.Script)
	for name := range names {
		if v, ok := s.ss[name]; ok {
			ss[name] = v
		}
	}
	return ss
}

func (s *scriptStore) SetDisable(name string, disable bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.ss[name]; ok {
		if v.Disable == disable {
			return false
		}
		v.Disable = disable
		return true
	}
	return false
}
