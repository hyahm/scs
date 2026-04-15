package store

import (
	"sync"

	"github.com/hyahm/scs/pkg/config/scripts"
)

type ScriptStore struct {
	mu sync.RWMutex
	ss map[string]*scripts.Script
}

func NewScriptStore() *ScriptStore {
	return &ScriptStore{
		ss: make(map[string]*scripts.Script),
	}
}

func (s *ScriptStore) Get(name string) (*scripts.Script, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.ss[name]
	return v, ok
}

func (s *ScriptStore) Set(script *scripts.Script) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ss[script.Name] = script
}

func (s *ScriptStore) Delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.ss, name)
}

func (s *ScriptStore) GetAll() map[string]*scripts.Script {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ss
}

func (s *ScriptStore) GetMapFilterByName(names map[string]struct{}) map[string]*scripts.Script {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ss := make(map[string]*scripts.Script)
	for name := range names {
		if v, ok := s.ss[name]; ok {
			ss[name] = v
		}
	}
	return ss
}

func (s *ScriptStore) SetDisable(name string, disable bool) bool {
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
