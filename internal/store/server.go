package store

import (
	"fmt"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg"
)

// 判断是否存在这个script
func (s *store) InitServer(index int, pname, name string) *server.Server {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.servers[name] = &server.Server{
		Index:   index,
		Name:    pname,
		SubName: name,
	}
	return s.servers[name]
}

func (s *store) SetServer(name string, srv *server.Server) *server.Server {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.servers[name] = srv
	return s.servers[name]
}

func (s *store) GetServerByName(name string) (*server.Server, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.servers[name]
	return v, ok
}

func (s *store) GetAllServer() []*server.Server {
	s.mu.RLock()
	defer s.mu.RUnlock()
	servers := make([]*server.Server, 0)
	for _, svc := range s.servers {
		servers = append(servers, svc)
	}
	return servers
}

func (s *store) DeleteServerByName(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.servers, name)
}

func (s *store) GetAllServerMap() map[string]*server.Server {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.servers
}

func (s *store) GetServerMapFilterScripts(names map[string]struct{}) map[string]*server.Server {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sm := make(map[string]*server.Server)
	for pname := range names {
		for index := range s.serverIndex[pname] {
			subname := fmt.Sprintf("%s_%d", pname, index)
			if _, ok := s.servers[subname]; !ok {
				golog.Error(pkg.ErrBugMsg)
				continue
			}
			sm[subname] = s.servers[subname]
		}
	}
	return sm
}
