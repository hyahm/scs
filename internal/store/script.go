package store

import "github.com/hyahm/scs/pkg/config/scripts"

func (s *store) GetScriptByName(pname string) (*scripts.Script, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.ss[pname]
	return v, ok
}

// 返回是否被修改
func (s *store) SetScriptDisable(pname string, disable bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.ss[pname]
	if !ok {
		return false
	}
	if s.ss[pname].Disable == disable {
		return false
	}
	s.ss[pname].Disable = disable
	return true
}

func (s *store) DeleteScriptByName(pname string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.ss, pname)
}

func (s *store) GetAllScriptMap() map[string]*scripts.Script {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ss
}

func (s *store) SetScript(script *scripts.Script) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ss[script.Name] = script
}

func (s *store) GetScriptMapFilterByName(names map[string]struct{}) map[string]*scripts.Script {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ss := make(map[string]*scripts.Script)
	for name := range names {
		if _, ok := s.ss[name]; ok {
			ss[name] = s.ss[name]
		}

	}
	return ss
}
