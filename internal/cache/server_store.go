package cache

import (
	"sync"
)

type serverStore struct {
	sync.RWMutex
	// Server 只有一个会
	server map[string]*Server // 保存脚本状态// 保存的 subname  中间有_{index}
}

var storeInstance = &serverStore{
	server: make(map[string]*Server),
}

func SetServer(name string, srv *Server) {
	storeInstance.Lock()
	defer storeInstance.Unlock()
	storeInstance.server[name] = srv
}

func GetServer(name string) (*Server, bool) {
	storeInstance.RLock()
	defer storeInstance.RUnlock()
	v, ok := storeInstance.server[name]
	return v, ok
}

func GetGroupServer(name string) []*Server {
	storeInstance.RLock()
	defer storeInstance.RUnlock()
	ss := make([]*Server, 0)
	for _, v := range storeInstance.server {
		if v.Name == name {
			ss = append(ss, v)
		}
	}
	return ss
}

func GetAllServer() []*Server {
	storeInstance.RLock()
	defer storeInstance.RUnlock()
	ss := make([]*Server, 0)
	for _, v := range storeInstance.server {
		ss = append(ss, v)
	}
	return ss
}

func RemoveServer(name string) {
	storeInstance.Lock()
	defer storeInstance.Unlock()
	delete(storeInstance.server, name)
}

func RemoveGroupServer(pname string) {
	storeInstance.Lock()
	defer storeInstance.Unlock()
	for k, v := range storeInstance.server {
		if v.Name == pname {
			delete(storeInstance.server, k)
		}
	}
}

func GetAllServerMap() map[string]*Server {
	storeInstance.RLock()
	defer storeInstance.RUnlock()
	result := make(map[string]*Server, len(storeInstance.server))
	for k, v := range storeInstance.server {
		result[k] = v
	}
	return result
}

// GetStore 返回全局 Store 实例

// func (s *Store) GetServerBySubName(subname string) (Server, bool) {
// 	s.RLock()
// 	defer s.RUnlock()
// 	v, ok := s.server[subname]
// 	return v, ok
// }

// func (s *store) SetScriptDisable(pname string, disable bool) bool {
// 	return s.Scripts.SetDisable(pname, disable)
// }

// func (s *store) DeleteScriptByName(pname string) {
// 	s.Scripts.Delete(pname)
// }

// func (s *store) GetAllScriptMap() map[string]config.Script {
// 	return s.Scripts.GetAll()
// }

// func (s *store) SetScript(script config.Script) {
// 	s.Scripts.Set(script)
// }

// func (s *store) GetScriptMapFilterByName(names map[string]struct{}) map[string]config.Script {
// 	return s.Scripts.GetMapFilterByName(names)
// }

// func (s *store) InitServer(index int, pname, name string) *server.Server {
// 	return s.Servers.Init(index, pname, name)
// }

// func (s *store) SetServer(name string, svc *server.Server) {
// 	s.Servers.set(name, svc)
// }

// func (s *store) GetServerByName(name string) (*server.Server, bool) {
// 	return s.Servers.get(name)
// }

// func (s *store) GetAllServer() []*server.Server {
// 	return s.Servers.getAll()
// }

// func (s *store) DeleteServerByName(name string) {
// 	s.Servers.delete(name)
// }

// func (s *store) GetAllServerMap() map[string]*server.Server {
// 	return s.Servers.getAllMap()
// }

// func (s *store) GetServerMapFilterScripts(names map[string]struct{}) map[string]*server.Server {
// 	result := make(map[string]*server.Server)
// 	for name := range names {
// 		indexs := s.Index.Get(name)
// 		if len(indexs) == 0 {
// 			continue
// 		}
// 		indexMap := make(map[int]struct{})
// 		for _, i := range indexs {
// 			indexMap[i] = struct{}{}
// 		}
// 		for subname, srv := range s.Servers.GetByScript(name, indexMap) {
// 			result[subname] = srv
// 		}
// 	}
// 	return result
// }

// func (s *store) SetScriptIndex(pname string, i int) {
// 	s.Index.Set(pname, i)
// }

// func (s *store) DeleteScriptIndex(pname string, i int) {
// 	s.Index.Delete(pname, i)
// }

// func (s *store) GetScriptIndex(pname string) []int {
// 	return s.Index.Get(pname)
// }

// func (s *store) GetScriptLength(pname string) int {
// 	return s.Index.Len(pname)
// }

// func (s *store) HaveServerByIndex(pname string, i int) bool {
// 	return s.Index.Has(pname, i)
// }

// func (s *serverStore) Init(index int, pname, name string) *server.Server {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	s.servers[name] = &server.Server{
// 		Index:   index,
// 		Name:    pname,
// 		SubName: name,
// 	}
// 	return s.servers[name]
// }

// func (s *serverStore) set(name string, svc *server.Server) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	s.servers[name] = svc
// }

// func (s *serverStore) get(name string) (*server.Server, bool) {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
// 	v, ok := s.servers[name]
// 	return v, ok
// }

// func (s *serverStore) getAll() []*server.Server {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
// 	servers := make([]*server.Server, 0, len(s.servers))
// 	for _, svc := range s.servers {
// 		servers = append(servers, svc)
// 	}
// 	return servers
// }

// func (s *serverStore) getAllMap() map[string]*server.Server {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
// 	return s.servers
// }

// func (s *serverStore) delete(name string) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	delete(s.servers, name)
// }

// func (s *serverStore) GetByScript(pname string, indexMap map[int]struct{}) map[string]*server.Server {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
// 	sm := make(map[string]*server.Server)
// 	for index := range indexMap {
// 		subname := fmt.Sprintf("%s_%d", pname, index)
// 		if v, ok := s.servers[subname]; ok {
// 			sm[subname] = v
// 		} else {
// 			golog.Error(pkg.ErrBugMsg)
// 		}
// 	}
// 	return sm
// }

// func (s *indexStore) Set(pname string, i int) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	if _, ok := s.serverIndex[pname]; !ok {
// 		s.serverIndex[pname] = make(map[int]struct{})
// 	}
// 	s.serverIndex[pname][i] = struct{}{}
// }

// func (s *indexStore) Delete(pname string, i int) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	if _, ok := s.serverIndex[pname]; !ok {
// 		return
// 	}
// 	delete(s.serverIndex[pname], i)
// }

// func (s *indexStore) Get(pname string) []int {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
// 	indexs := make([]int, 0)
// 	if _, ok := s.serverIndex[pname]; !ok {
// 		return indexs
// 	}
// 	for index := range s.serverIndex[pname] {
// 		indexs = append(indexs, index)
// 	}
// 	return indexs
// }

// func (s *indexStore) Len(pname string) int {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
// 	if _, ok := s.serverIndex[pname]; !ok {
// 		return 0
// 	}
// 	return len(s.serverIndex[pname])
// }

// func (s *indexStore) Has(pname string, i int) bool {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
// 	if _, ok := s.serverIndex[pname]; !ok {
// 		return false
// 	}
// 	_, ok := s.serverIndex[pname][i]
// 	return ok
// }
