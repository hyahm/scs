package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"scs/internal"
	"scs/pkg/script"
	"sync"

	"github.com/hyahm/golog"
)

var UseNodes string
var GroupName string

type Node struct {
	Name  string
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
	Wg    *sync.WaitGroup
}

func (node *Node) Reload() {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	b, err := Requests("POST", fmt.Sprintf("%s/-/reload", node.Url), node.Token, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func (node *Node) Restart(args ...string) {
	fmt.Println(string(node.crud("restart", args...)))

}

type SearchInfo struct {
	Name string `json:"name"`
	Info string `json:"info"`
}

func (node *Node) Search(args string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	b, err := Requests("POST", fmt.Sprintf("%s/get/repo", node.Url), node.Token, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp := &struct {
		Url        []string `json:"url"`
		Derivative string   `json:"derivative"`
	}{}
	err = json.Unmarshal(b, resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	sl := make([]*SearchInfo, 0)
	for _, url := range resp.Url {
		b, err := Requests("POST", fmt.Sprintf("%s/search/%s/%s", url, resp.Derivative, args), node.Token, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		si := make([]*SearchInfo, 0)
		err = json.Unmarshal(b, &si)
		if err != nil {
			golog.Error(err)
			continue
		}
		sl = append(sl, si...)
	}
	for _, searchinfo := range sl {
		fmt.Printf("name: %s \t info: %s \n", searchinfo.Name, searchinfo.Info)
	}

}

func (node *Node) Start(args ...string) {
	fmt.Println(string(node.crud("start", args...)))
}

func (node *Node) Status(args ...string) {
	b := node.crud("status")
	if len(b) == 0 {
		return
	}
	var s status = make([]*script.ServiceStatus, 0)
	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println(string(b))
		fmt.Println(err.Error() + " or token error")
		return
	}
	s.sortAndPrint(node.Name, node.Url)
}

func (node *Node) Kill(args ...string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	switch len(args) {
	case 1:
		b, err := Requests("POST", fmt.Sprintf("%s/kill/%s", node.Url, args[0]), node.Token, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
	default:
		b, err := Requests("POST", fmt.Sprintf("%s/kill/%s/%s", node.Url, args[0], args[1]), node.Token, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
	}
}

func (node *Node) Env(args ...string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	var b []byte
	var err error
	// switch len(args) {
	// case 1:
	b, err = Requests("POST", fmt.Sprintf("%s/env/%s", node.Url, args[0]), node.Token, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	l := make(map[string]string, 0)
	err = json.Unmarshal(b, &l)
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range l {
		fmt.Println(k + ": " + v)
	}
}

func (node *Node) getDependEnv(args ...string) {
	// 获取依赖的env
	b, err := Requests("POST", fmt.Sprintf("%s/env/%s", node.Url, args[0]), node.Token, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func (node *Node) Install(script *internal.Script, env map[string]string) {
	// 先读取配置文件
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	body, err := json.Marshal(script)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	b, err := Requests("POST", fmt.Sprintf("%s/script", node.Url), node.Token, bytes.NewReader(body))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))
}

func (node *Node) Log(args string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	b, err := Requests("POST", fmt.Sprintf("%s/log/%s", node.Url, args), node.Token, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func (node *Node) Stop(args ...string) {
	fmt.Println(string(node.crud("stop", args...)))
}

func (node *Node) crud(operate string, args ...string) []byte {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	var url string
	switch len(args) {
	case 0:
		url = fmt.Sprintf("%s/%s", node.Url, operate)
	case 1:
		url = fmt.Sprintf("%s/%s/%s", node.Url, operate, args[0])
	default:
		url = fmt.Sprintf("%s/%s/%s/%s", node.Url, operate, args[0], args[1])

	}
	b, err := Requests("POST", url, node.Token, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}

func (node *Node) Update(args ...string) {
	fmt.Println(string(node.crud("update", args...)))
}
