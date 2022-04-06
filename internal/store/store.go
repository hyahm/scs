package store

import (
	"sync"

	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg/config/scripts"
)

type store struct {
	servers     map[string]*server.Server
	serverIndex map[string]map[int]struct{}
	ss          map[string]*scripts.Script
	mu          sync.RWMutex
}

var Store *store

func init() {
	Store = &store{
		mu:          sync.RWMutex{},
		servers:     make(map[string]*server.Server),
		ss:          make(map[string]*scripts.Script),
		serverIndex: make(map[string]map[int]struct{}),
	}
}
