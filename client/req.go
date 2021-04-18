package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/hyahm/scs"
)

type ClientConfig struct {
	Nodes       map[string]*scs.Node `yaml:"nodes"`
	Group       map[string][]string  `yaml:"group"`
	ReadTimeout time.Duration        `yaml:"timeout"`
	mu          *sync.RWMutex
}

func NewClientConfig() *ClientConfig {
	return &ClientConfig{mu: &sync.RWMutex{}}
}

func (cc *ClientConfig) GetNode(name string) (*scs.Node, bool) {
	if cc.mu == nil {
		cc.mu = &sync.RWMutex{}
	}
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	if v, ok := cc.Nodes[name]; ok {
		v.Name = name
		return v, true
	} else {
		return nil, false
	}
}

func (cc *ClientConfig) GetNodes() []*scs.Node {
	ns := make([]*scs.Node, 0)
	if cc.mu == nil {
		cc.mu = &sync.RWMutex{}
	}
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	for name, v := range cc.Nodes {
		v.Name = name
		ns = append(ns, v)
	}
	return ns
}

func (cc *ClientConfig) GetNodesInGroup(group string) []*scs.Node {
	ns := make([]*scs.Node, 0)
	if cc.mu == nil {
		cc.mu = &sync.RWMutex{}
	}
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	if v, ok := cc.Group[group]; ok {
		for _, name := range v {
			if node, ok := cc.GetNode(name); ok {
				ns = append(ns, node)
			}
		}
		return ns
	} else {
		return nil
	}
}

func (cc *ClientConfig) PrintNodes() {
	if cc.mu == nil {
		cc.mu = &sync.RWMutex{}
	}
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	for name, v := range cc.Nodes {
		fmt.Printf("name: %s \t url: %s \t token: %s \n", name, v.Url, v.Token)
	}

}
