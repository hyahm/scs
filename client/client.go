package client

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

var CCfg *ClientConfig

func ReadClientConfig() {
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
	CCfg = NewClientConfig()
	err = yaml.Unmarshal(b, &CCfg)
	if err != nil {
		panic(err)
	}
	if CCfg.ReadTimeout == 0 {
		CCfg.ReadTimeout = time.Second * 3
	}
}
