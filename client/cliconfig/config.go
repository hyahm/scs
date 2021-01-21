package cliconfig

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/hyahm/scs/client/node"

	"gopkg.in/yaml.v2"
)

var Cfg *ClientConfig

type ClientConfig struct {
	Nodes       map[string]node.Node `yaml:"nodes"`
	Group       map[string][]string  `yaml:"group"`
	ReadTimeout time.Duration        `yaml:"readTimeout"`
}

func ReadConfig() {

	root, err := os.UserHomeDir()
	if err != nil {
		// 找不到就报错
		panic(err)
	}
	configfile := filepath.Join(root, "scsctl.yaml")
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
	Cfg = &ClientConfig{}
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
