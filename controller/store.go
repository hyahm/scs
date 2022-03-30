package controller

// 操作
import (
	"sync"

	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg/config/scripts"
)

type Store struct {
	servers     map[string]*server.Server
	serverIndex map[string]map[int]struct{}
	ss          map[string]*scripts.Script
	mu          sync.RWMutex
}

var store *Store

// // 保存的servers
// var servers map[string]*server.Server

// // 保存 server的index还有哪些
// var serverIndex map[string]map[int]struct{}

// // 保存的scripts
// var ss map[string]*scripts.Script

// var mu sync.RWMutex

func init() {
	store = &Store{
		mu:          sync.RWMutex{},
		servers:     make(map[string]*server.Server),
		ss:          make(map[string]*scripts.Script),
		serverIndex: make(map[string]map[int]struct{}),
	}
}

func GetServerBySubname(subname string) (*server.Server, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	v, ok := store.servers[subname]
	return v, ok
}
