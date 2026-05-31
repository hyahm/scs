package store

import (
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg/config"
)

type store struct {
	Scripts *scriptStore
	Servers *serverStore
	Index   *indexStore
}

var storeInstance *store

func init() {
	storeInstance = &store{
		Scripts: NewScriptStore(),
		Servers: NewServerStore(),
		Index:   NewIndexStore(),
	}
}

// GetStore 返回全局 Store 实例
func GetStore() *store {
	return storeInstance
}

func (s *store) GetScriptByName(pname string) (config.Script, bool) {
	return s.Scripts.Get(pname)
}

func (s *store) SetScriptDisable(pname string, disable bool) bool {
	return s.Scripts.SetDisable(pname, disable)
}

func (s *store) DeleteScriptByName(pname string) {
	s.Scripts.Delete(pname)
}

func (s *store) GetAllScriptMap() map[string]config.Script {
	return s.Scripts.GetAll()
}

func (s *store) SetScript(script config.Script) {
	s.Scripts.Set(script)
}

func (s *store) GetScriptMapFilterByName(names map[string]struct{}) map[string]config.Script {
	return s.Scripts.GetMapFilterByName(names)
}

func (s *store) InitServer(index int, pname, name string) *server.Server {
	return s.Servers.Init(index, pname, name)
}

func (s *store) SetServer(name string, svc *server.Server) {
	s.Servers.set(name, svc)
}

func (s *store) GetServerByName(name string) (*server.Server, bool) {
	return s.Servers.get(name)
}

func (s *store) GetAllServer() []*server.Server {
	return s.Servers.getAll()
}

func (s *store) DeleteServerByName(name string) {
	s.Servers.delete(name)
}

func (s *store) GetAllServerMap() map[string]*server.Server {
	return s.Servers.getAllMap()
}

func (s *store) GetServerMapFilterScripts(names map[string]struct{}) map[string]*server.Server {
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

func (s *store) SetScriptIndex(pname string, i int) {
	s.Index.Set(pname, i)
}

func (s *store) DeleteScriptIndex(pname string, i int) {
	s.Index.Delete(pname, i)
}

func (s *store) GetScriptIndex(pname string) []int {
	return s.Index.Get(pname)
}

func (s *store) GetScriptLength(pname string) int {
	return s.Index.Len(pname)
}

func (s *store) HaveServerByIndex(pname string, i int) bool {
	return s.Index.Has(pname, i)
}
