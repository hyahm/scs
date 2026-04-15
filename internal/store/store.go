package store

import (
	"sync"

	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg/config/scripts"
)

type Store struct {
	Scripts *ScriptStore
	Servers *ServerStore
	Index   *IndexStore
	mu      sync.RWMutex
}

var storeInstance *Store

func init() {
	storeInstance = &Store{
		Scripts: NewScriptStore(),
		Servers: NewServerStore(),
		Index:   NewIndexStore(),
	}
}

// GetStore 返回全局 Store 实例
func GetStore() *Store {
	return storeInstance
}

func (s *Store) GetScriptByName(pname string) (*scripts.Script, bool) {
	return s.Scripts.Get(pname)
}

func (s *Store) SetScriptDisable(pname string, disable bool) bool {
	return s.Scripts.SetDisable(pname, disable)
}

func (s *Store) DeleteScriptByName(pname string) {
	s.Scripts.Delete(pname)
}

func (s *Store) GetAllScriptMap() map[string]*scripts.Script {
	return s.Scripts.GetAll()
}

func (s *Store) SetScript(script *scripts.Script) {
	s.Scripts.Set(script)
}

func (s *Store) GetScriptMapFilterByName(names map[string]struct{}) map[string]*scripts.Script {
	return s.Scripts.GetMapFilterByName(names)
}

func (s *Store) InitServer(index int, pname, name string) *server.Server {
	return s.Servers.Init(index, pname, name)
}

func (s *Store) SetServer(name string, srv *server.Server) *server.Server {
	return s.Servers.Set(name, srv)
}

func (s *Store) GetServerByName(name string) (*server.Server, bool) {
	return s.Servers.Get(name)
}

func (s *Store) GetAllServer() []*server.Server {
	return s.Servers.GetAll()
}

func (s *Store) DeleteServerByName(name string) {
	s.Servers.Delete(name)
}

func (s *Store) GetAllServerMap() map[string]*server.Server {
	return s.Servers.GetAllMap()
}

func (s *Store) GetServerMapFilterScripts(names map[string]struct{}) map[string]*server.Server {
	result := make(map[string]*server.Server)
	for name := range names {
		indexs := s.Index.Get(name)
		if len(indexs) == 0 {
			continue
		}
		indexMap := make(map[int]struct{})
		for _, i := range indexs {
			indexMap[i] = struct{}{}
		}
		for subname, srv := range s.Servers.GetByScript(name, indexMap) {
			result[subname] = srv
		}
	}
	return result
}

func (s *Store) SetScriptIndex(pname string, i int) {
	s.Index.Set(pname, i)
}

func (s *Store) DeleteScriptIndex(pname string, i int) {
	s.Index.Delete(pname, i)
}

func (s *Store) GetScriptIndex(pname string) []int {
	return s.Index.Get(pname)
}

func (s *Store) GetScriptLength(pname string) int {
	return s.Index.Len(pname)
}

func (s *Store) HaveServerByIndex(pname string, i int) bool {
	return s.Index.Has(pname, i)
}
