package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/scripts"
)

// 已经支持多服务器操作， 每台服务器相当于一个node
type Node struct {
	Name    string        `yaml:"-"`
	Url     string        `yaml:"url"`
	Token   string        `yaml:"token"`
	Role    string        `json:""`
	Timeout time.Duration `json:"timeout"`
}

func (node *Node) NewSCSClient() *SCSClient {
	return &SCSClient{
		Domain:  node.Url,
		Token:   node.Token,
		Timeout: node.Timeout,
	}
}

func (node *Node) Reload() {
	res, err := node.NewSCSClient().Reload()
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	fmt.Println(string(res.Marshal()))
}

func (node *Node) Fmt() {
	res, err := node.NewSCSClient().Fmt()
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	fmt.Println(string(res.Marshal()))
}

func (node *Node) Restart(args ...string) {
	cli := node.NewSCSClient()
	var err error
	var res *pkg.Response
	switch len(args) {
	case 0:
		res, err = cli.RestartAll()
	case 1:
		cli.Pname = args[0]
		res, err = cli.RestartPname()
	default:
		cli.Pname = args[0]
		cli.Name = args[1]
		res, err = cli.RestartName()
	}
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	fmt.Println(string(res.Marshal()))
	// fmt.Println(string(node.crud("stop", args...)))
}

type SearchInfo struct {
	Name string `json:"name"`
	Info string `json:"info"`
}

func (node *Node) Search(args string) {

	// b, err := node.NewSCSClient().Repo()
	// // b, err := Requests("POST", fmt.Sprintf("%s/get/repo", node.Url), node.Token, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// resp := &struct {
	// 	Url        []string `json:"url"`
	// 	Derivative string   `json:"derivative"`
	// }{}
	// err = json.Unmarshal(b, resp)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// sl := make([]*SearchInfo, 0)
	// for _, url := range resp.Url {
	// 	node.NewSCSClient().Domain = url
	// 	b, err := node.NewSCSClient().Search(resp.Derivative, args)
	// 	// b, err := Requests("POST", fmt.Sprintf("%s/search/%s/%s", url, resp.Derivative, args), node.Token, nil)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	si := make([]*SearchInfo, 0)
	// 	err = json.Unmarshal(b, &si)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	sl = append(sl, si...)
	// }
	// for _, searchinfo := range sl {
	// 	fmt.Printf("name: %s \t info: %s \n", searchinfo.Name, searchinfo.Info)
	// }

}

func (node *Node) Start(args ...string) {

	cli := node.NewSCSClient()
	var err error
	var res *pkg.Response
	switch len(args) {
	case 0:
		res, err = cli.StartAll()
	case 1:
		cli.Pname = args[0]
		res, err = cli.StartPname()
	default:
		cli.Pname = args[0]
		cli.Name = args[1]
		res, err = cli.StartName()
	}
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}

	fmt.Println(string(res.Marshal()))
	// fmt.Println(string(node.crud("stop", args...)))
}

func (node *Node) Status(args ...string) (*ScriptStatusNode, error) {

	cli := node.NewSCSClient()
	var ssn *pkg.StatusList
	var err error
	switch len(args) {
	case 0:
		ssn, err = cli.StatusAll()
	case 1:
		cli.Pname = args[0]
		ssn, err = cli.StatusPname()
	default:
		cli.Pname = args[0]
		cli.Name = args[1]
		ssn, err = cli.StatusName()
	}
	if err != nil {
		return nil, fmt.Errorf("url: %s, name: %s, msg: %v", node.Url, node.Name, err)
	}

	result := &ScriptStatusNode{
		Nodes:   ssn.Data,
		Version: ssn.Version,
		Name:    node.Name,
		Url:     node.Url,
	}

	if len(ssn.Data) > 0 {
		result.Version = ssn.Version
	}

	return result, nil
}

func (node *Node) Kill(args ...string) {

	cli := node.NewSCSClient()
	var res *pkg.Response
	var err error
	switch len(args) {
	case 2:
		cli.Pname = args[0]
		cli.Name = args[1]
		res, err = cli.KillName()
	case 1:
		cli.Pname = args[0]
		res, err = cli.KillPname()
	default:
		return
	}
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	fmt.Println(string(res.Marshal()))
}

func (node *Node) Env(args string) {

	cli := node.NewSCSClient()
	cli.Name = args
	res, err := cli.Env()
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}

	out, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	fmt.Println(string(out))
}

func (node *Node) Info(args string) {

	cli := node.NewSCSClient()
	cli.Name = args
	res, err := cli.Info()
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}

	out, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
	}
	fmt.Println(string(out))
}

func (node *Node) Install(scripts []*scripts.Script, env map[string]string) {
	// 先读取配置文件

	for _, script := range scripts {
		cli := node.NewSCSClient()
		res, err := cli.AddScript(script)
		if err != nil {
			fmt.Printf("name: %s , msg: %v\n", node.Name, err)
			return
		}
		fmt.Println(string(res.Marshal()))
	}

}

func (node *Node) Log(args string, line int) {
	cli := node.NewSCSClient()
	cli.Name = args
	cli.Log(line)
}

func (node *Node) Stop(args ...string) {

	cli := node.NewSCSClient()
	var err error
	var res *pkg.Response
	switch len(args) {
	case 0:
		res, err = cli.StopAll()
	case 1:
		cli.Pname = args[0]
		res, err = cli.StopPname()
	default:
		cli.Pname = args[0]
		cli.Name = args[1]
		res, err = cli.StopName()
	}
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	fmt.Println(string(res.Marshal()))

}

func (node *Node) Remove(args ...string) {

	cli := node.NewSCSClient()
	var err error
	var res *pkg.Response
	switch len(args) {
	case 0:
		fmt.Printf("name: %s , msg: remove all have been removed\n", node.Name)
		return
		// b, err = cli.RemoveAllScrip()
	case 1:
		cli.Pname = args[0]
		res, err = cli.RemovePnameScrip()
	default:
		cli.Pname = args[0]
		cli.Name = args[1]
		res, err = cli.RemoveNameScrip()
	}
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	fmt.Println(string(res.Marshal()))

}

func (node *Node) Enable(pname string) {

	cli := node.NewSCSClient()
	cli.Pname = pname
	res, err := cli.Enable()
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	fmt.Println(string(res.Marshal()))
}

func (node *Node) GetServers() {
	cli := node.NewSCSClient()
	res, err := cli.GetServers()
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	out, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	fmt.Println(string(out))
}

func (node *Node) GetIndexs() {
	cli := node.NewSCSClient()
	res, err := cli.GetIndexs()
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	out, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	fmt.Println(string(out))
}

func (node *Node) GetAlerts() {
	cli := node.NewSCSClient()
	res, err := cli.GetAlarms()
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	out, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	fmt.Println(string(out))
}

func (node *Node) GetScripts() {
	cli := node.NewSCSClient()
	res, err := cli.GetScripts()
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	out, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	fmt.Println(string(out))
}

func (node *Node) Disable(pname string) {

	cli := node.NewSCSClient()
	cli.Pname = pname
	res, err := cli.Disable()
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	fmt.Println(string(res.Marshal()))
	// fmt.Println(string(node.crud("stop", args...)))
}

func (node *Node) Update(args ...string) {

	cli := node.NewSCSClient()
	var err error
	var res *pkg.Response
	switch len(args) {
	case 0:
		res, err = cli.UpdateAll()
	case 1:
		cli.Pname = args[0]
		res, err = cli.UpdatePname()
	default:
		cli.Pname = args[0]
		cli.Name = args[1]
		res, err = cli.UpdateName()
	}
	if err != nil {
		fmt.Printf("name: %s , msg: %v\n", node.Name, err)
		return
	}
	fmt.Println(string(res.Marshal()))
}
