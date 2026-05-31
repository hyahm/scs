package store

import (
	"fmt"
	"sync"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg"
)

type serverStore struct {
	mu      sync.RWMutex
	servers map[string]*server.Server
}

func NewServerStore() *serverStore {
	return &serverStore{
		mu:      sync.RWMutex{},
		servers: make(map[string]*server.Server),
	}
}

func (s *serverStore) Init(index int, pname, name string) *server.Server {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.servers[name] = &server.Server{
		Index:   index,
		Name:    pname,
		SubName: name,
	}
	return s.servers[name]
}

func (s *serverStore) set(name string, svc *server.Server) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.servers[name] = svc
}

func (s *serverStore) get(name string) (*server.Server, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.servers[name]
	return v, ok
}

func (s *serverStore) getAll() []*server.Server {
	s.mu.RLock()
	defer s.mu.RUnlock()
	servers := make([]*server.Server, 0, len(s.servers))
	for _, svc := range s.servers {
		servers = append(servers, svc)
	}
	return servers
}

func (s *serverStore) getAllMap() map[string]*server.Server {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.servers
}

func (s *serverStore) delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.servers, name)
}

func (s *serverStore) GetByScript(pname string, indexMap map[int]struct{}) map[string]*server.Server {
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
