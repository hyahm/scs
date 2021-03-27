package cliconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hyahm/scs/client/node"

	"gopkg.in/yaml.v2"
)

var Cfg *ClientConfig

type ClientConfig struct {
	Nodes       map[string]*node.Node `yaml:"nodes"`
	Group       map[string][]string   `yaml:"group"`
	ReadTimeout time.Duration         `yaml:"readTimeout"`
	mu          *sync.RWMutex
}

func NewClientConfig() *ClientConfig {
	return &ClientConfig{mu: &sync.RWMutex{}}
}

func (cc *ClientConfig) GetNode(name string) (*node.Node, bool) {
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

func (cc *ClientConfig) GetNodes() []*node.Node {
	ns := make([]*node.Node, 0)
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

func (cc *ClientConfig) GetNodesInGroup(group string) []*node.Node {
	ns := make([]*node.Node, 0)
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

func ReadConfig() {
	root, err := os.UserHomeDir()
	if err != nil {
		// 找不到就报错
		panic(err)
	}
	configfile := filepath.Join(root, ".scsctl.yaml")
	_, err = os.Stat(configfile)
	if err != nil {
		_, err = os.Create(configfile)
		if err != nil {
			panic(err)
		}
	}
	b, err := ioutil.ReadFile(configfile)
	if err != nil {
		panic(err)
	}
	if len(b) == 0 {
		x := `nodes:
  local: 
    url: "https://127.0.0.1:11111"
    token:  
group: `
		b = []byte(x)
		err := ioutil.WriteFile(configfile, b, 0644)
		if err != nil {
			panic(err)
		}
	}
	Cfg = NewClientConfig()
	err = yaml.Unmarshal(b, Cfg)
	if err != nil {
		panic(err)
	}
	if Cfg.ReadTimeout == 0 {
		node.ReadTimeout = time.Second * 3
	} else {
		node.ReadTimeout = Cfg.ReadTimeout
	}
}
