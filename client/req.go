package client

import (
	"fmt"
	"time"
)

type ClientConfig struct {
	Nodes       map[string]*Node    `yaml:"nodes"`
	Group       map[string][]string `yaml:"group"`
	ReadTimeout time.Duration       `yaml:"timeout"`
}

func NewClientConfig() *ClientConfig {
	return &ClientConfig{}
}

func (cc *ClientConfig) GetNode(name string) (*Node, bool) {
	if v, ok := cc.Nodes[name]; ok {
		v.Name = name
		v.Timeout = cc.ReadTimeout
		return v, true
	} else {
		return nil, false
	}
}

func (cc *ClientConfig) GetNodes() []*Node {
	ns := make([]*Node, 0)
	for name, v := range cc.Nodes {
		v.Name = name
		v.Timeout = cc.ReadTimeout
		ns = append(ns, v)
	}
	return ns
}

func (cc *ClientConfig) GetNodesInGroup(group string) []*Node {
	ns := make([]*Node, 0)
	if v, ok := cc.Group[group]; ok {
		for _, name := range v {
			if node, ok := cc.GetNode(name); ok {
				node.Timeout = cc.ReadTimeout
				ns = append(ns, node)
			}
		}
		return ns
	} else {
		return nil
	}
}

func (cc *ClientConfig) PrintNodes() {
	for name, v := range cc.Nodes {
		fmt.Printf("name: %s \t url: %s \t token: %s \n", name, v.Url, v.Token)
	}

}
