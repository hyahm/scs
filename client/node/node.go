package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hyahm/scs/client"
	"github.com/hyahm/scs/script"
)

var UseNodes string
var GroupName string
var ReadTimeout time.Duration
var ErrorToken = errors.New("Token error")

type Node struct {
	Name   string `yaml:"-"`
	Url    string `yaml:"url"`
	Token  string `yaml:"token"`
	Filter []string
	Result *ScriptStatusNode
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
	cli := node.NewSCSClient()
	var b []byte
	var err error
	switch len(args) {
	case 0:
		b, err = cli.RestartAll()
	case 1:
		cli.Pname = args[0]
		b, err = cli.RestartPname()
	default:
		cli.Pname = args[0]
		cli.Name = args[1]
		b, err = cli.RestartName()
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	// fmt.Println(string(node.crud("stop", args...)))
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
	cli := node.NewSCSClient()
	var b []byte
	var err error
	switch len(args) {
	case 0:
		b, err = cli.StartAll()
	case 1:
		cli.Pname = args[0]
		b, err = cli.StartPname()
	default:
		cli.Pname = args[0]
		cli.Name = args[1]
		b, err = cli.StartName()
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	// fmt.Println(string(node.crud("stop", args...)))
}

func (node *Node) Status(args ...string) error {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	cli := node.NewSCSClient()
	var b []byte
	var err error
	switch len(args) {
	case 0:
		b, err = cli.StatusAll()
	case 1:
		cli.Pname = args[0]
		b, err = cli.StatusPname()
	default:
		cli.Pname = args[0]
		cli.Name = args[1]
		b, err = cli.StartName()
	}
	resp := &script.StatusList{}
	// fmt.Println(string(b))
	err = json.Unmarshal(b, resp)
	if err != nil {
		fmt.Printf("node: %s, url: %s %v \n", node.Name, node.Url, err)
		return err
	}
	if resp.Code == 203 {
		fmt.Printf("node: %s, url: %s %v \n", node.Name, node.Url, ErrorToken)
		return ErrorToken
	}

	if len(node.Filter) > 0 {
		resp.Filter(node.Filter)
	}
	node.Result = &ScriptStatusNode{}
	node.Result.Nodes = resp.Data
	node.Result.Name = node.Name
	node.Result.Url = node.Url
	node.Result.Filter = node.Filter
	return nil
}

func (node *Node) Kill(args ...string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	cli := node.NewSCSClient()
	var b []byte
	var err error

	switch len(args) {
	case 2:
		b, err = cli.KillName()
	case 1:
		cli.Pname = args[0]
		b, err = cli.KillPname()
	default:
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	// fmt.Println(string(node.crud("stop", args...)))
}

func (node *Node) Env(args string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	// var b []byte
	// var err error
	// switch len(args) {
	// case 1:
	cli := node.NewSCSClient()
	cli.Name = args
	b, err := cli.Env()
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

func (node *Node) Install(scripts []*client.Script, env map[string]string) {
	// 先读取配置文件
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	for _, script := range scripts {
		cli := node.NewSCSClient()
		b, err := cli.AddScript(script)
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
	cli := node.NewSCSClient()
	cli.Name = args
	b, err := cli.Log()
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
	cli := node.NewSCSClient()
	var b []byte
	var err error
	switch len(args) {
	case 0:
		b, err = cli.StopAll()
	case 1:
		cli.Pname = args[0]
		b, err = cli.StopPname()
	default:
		cli.Pname = args[0]
		cli.Name = args[1]
		b, err = cli.StopName()
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	// fmt.Println(string(node.crud("stop", args...)))
}

func (node *Node) Remove(args ...string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	cli := node.NewSCSClient()
	var b []byte
	var err error
	switch len(args) {
	case 0:
		b, err = cli.RemoveAllScrip()
	case 1:
		cli.Pname = args[0]
		b, err = cli.RemovePnameScrip()
	default:
		cli.Pname = args[0]
		cli.Name = args[1]
		b, err = cli.RemoveNameScrip()
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	// fmt.Println(string(node.crud("stop", args...)))
}

func (node *Node) Enable(pname string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	cli := node.NewSCSClient()
	cli.Pname = pname
	b, err := cli.Enable()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	// fmt.Println(string(node.crud("stop", args...)))
}

func (node *Node) Disable(pname string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	cli := node.NewSCSClient()
	cli.Pname = pname
	b, err := cli.Disable()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	// fmt.Println(string(node.crud("stop", args...)))
}

func (node *Node) Update(args ...string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	cli := node.NewSCSClient()
	var b []byte
	var err error
	switch len(args) {
	case 0:
		b, err = cli.UpdateAll()
	case 1:
		cli.Pname = args[0]
		b, err = cli.UpdatePname()
	default:
		cli.Pname = args[0]
		cli.Name = args[1]
		b, err = cli.UpdateName()
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	// return string(b)
	// return node.Update(args...)
	// fmt.Println(string(node.crud("update", args...)))
}
