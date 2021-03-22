package node

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/script"

	"github.com/hyahm/golog"
)

var UseNodes string
var GroupName string

type Node struct {
	Name  string `yaml:"-"`
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
	// Sc    *client.SCSClient
	Wg *sync.WaitGroup
}

func (node *Node) NewSCSClient() *client.SCSClient {
	return &client.SCSClient{
		Domain: node.Url,
		Token:  node.Token,
	}
}

func (node *Node) Reload() {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	b, err := node.NewSCSClient().Reload()
	// b, err := Requests("POST", fmt.Sprintf("%s/-/reload", node.Url), node.Token, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func (node *Node) Restart(args ...string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	// fmt.Println(string(node.crud("restart", args...)))
	b, err := node.NewSCSClient().Restart(args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

type SearchInfo struct {
	Name string `json:"name"`
	Info string `json:"info"`
}

func (node *Node) Search(args string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	b, err := node.NewSCSClient().Repo()
	// b, err := Requests("POST", fmt.Sprintf("%s/get/repo", node.Url), node.Token, nil)
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
		node.NewSCSClient().Domain = url
		b, err := node.NewSCSClient().Search(resp.Derivative, args)
		// b, err := Requests("POST", fmt.Sprintf("%s/search/%s/%s", url, resp.Derivative, args), node.Token, nil)
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
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	b, err := node.NewSCSClient().Start(args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	// fmt.Println(string(node.crud("start", args...)))
}

func (node *Node) Status(args ...string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	b, err := node.NewSCSClient().Status(args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	sl := &script.StatusList{}
	// fmt.Println(string(b))
	// var s status = make([]*script.ServiceStatus, 0)
	err = json.Unmarshal(b, sl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if sl.Code == 203 {
		fmt.Println("token error")
		return
	}
	status(sl.Data).sortAndPrint(node.Name, node.Url)
}

func (node *Node) Kill(args ...string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	b, err := node.NewSCSClient().Kill(args...)
	// b, err := Requests("POST", fmt.Sprintf("%s/kill/%s/%s", node.Url, args[0], args[1]), node.Token, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	// switch len(args) {
	// case 1:
	// 	b, err := Requests("POST", fmt.Sprintf("%s/kill/%s", node.Url, args[0]), node.Token, nil)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	fmt.Println(string(b))
	// default:
	// 	b, err := Requests("POST", fmt.Sprintf("%s/kill/%s/%s", node.Url, args[0], args[1]), node.Token, nil)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	fmt.Println(string(b))
	// }
}

func (node *Node) Env(args string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	// var b []byte
	// var err error
	// switch len(args) {
	// case 1:
	b, err := node.NewSCSClient().Env(args)
	// b, err = Requests("POST", fmt.Sprintf("%s/env/%s", node.Url, args[0]), node.Token, nil)
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

// func (node *Node) getDependEnv(args ...string) {
// 	// 获取依赖的env
// 	b, err := Requests("POST", fmt.Sprintf("%s/env/%s", node.Url, args[0]), node.Token, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(string(b))
// }

func (node *Node) Install(scripts []*internal.Script, env map[string]string) {
	// 先读取配置文件
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	for _, script := range scripts {
		b, err := node.NewSCSClient().Script(script)
		// b, err := Requests("POST", fmt.Sprintf("%s/script", node.Url), node.Token, bytes.NewReader(body))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
	}

}

func (node *Node) Log(args string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	b, err := node.NewSCSClient().Log(args)
	// b, err := Requests("POST", fmt.Sprintf("%s/log/%s", node.Url, args), node.Token, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func (node *Node) Stop(args ...string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	b, err := node.NewSCSClient().Stop(args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	// fmt.Println(string(node.crud("stop", args...)))
}

// func (node *Node) crud(operate string, args ...string) []byte {
// 	if node.Wg != nil {
// 		defer node.Wg.Done()
// 	}
// 	var url string
// 	switch len(args) {
// 	case 0:
// 		url = fmt.Sprintf("%s/%s", node.Url, operate)
// 	case 1:
// 		url = fmt.Sprintf("%s/%s/%s", node.Url, operate, args[0])
// 	default:
// 		url = fmt.Sprintf("%s/%s/%s/%s", node.Url, operate, args[0], args[1])
// 	}
// 	b, err := Requests("POST", url, node.Token, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil
// 	}
// 	return b
// }

func (node *Node) Update(args ...string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	b, err := node.NewSCSClient().Update(args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	// return string(b)
	// return node.Update(args...)
	// fmt.Println(string(node.crud("update", args...)))
}
