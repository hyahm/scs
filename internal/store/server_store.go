package store

import (
	"fmt"
	"sync"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg"
)

type ServerStore struct {
	mu      sync.RWMutex
	servers map[string]*server.Server
}

func NewServerStore() *ServerStore {
	return &ServerStore{
		servers: make(map[string]*server.Server),
	}
}

func (s *ServerStore) Init(index int, pname, name string) *server.Server {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.servers[name] = &server.Server{
		Index:   index,
		Name:    pname,
		SubName: name,
	}
	return s.servers[name]
}

func (s *ServerStore) Set(name string, srv *server.Server) *server.Server {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.servers[name] = srv
	return s.servers[name]
}

func (s *ServerStore) Get(name string) (*server.Server, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.servers[name]
	return v, ok
}

func (s *ServerStore) GetAll() []*server.Server {
	s.mu.RLock()
	defer s.mu.RUnlock()
	servers := make([]*server.Server, 0, len(s.servers))
	for _, svc := range s.servers {
		servers = append(servers, svc)
	}
	return servers
}

func (s *ServerStore) GetAllMap() map[string]*server.Server {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.servers
}

func (s *ServerStore) Delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.servers, name)
}

func (s *ServerStore) GetByScript(pname string, indexMap map[int]struct{}) map[string]*server.Server {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sm := make(map[string]*server.Server)
	for index := range indexMap {
		subname := fmt.Sprintf("%s_%d", pname, index)
		if v, ok := s.servers[subname]; ok {
			sm[subname] = v
		} else {
			golog.Error(pkg.ErrBugMsg)
		}
	}
	return sm
}
