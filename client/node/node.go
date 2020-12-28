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
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	switch len(args) {
	case 0:
		b, err := Requests("POST", fmt.Sprintf("%s/restart", node.Url), node.Token, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
	case 1:
		b, err := Requests("POST", fmt.Sprintf("%s/restart/%s", node.Url, args[0]), node.Token, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(b))
	default:
		b, err := Requests("POST", fmt.Sprintf("%s/restart/%s/%s", node.Url, args[0], args[1]), node.Token, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(b))

	}

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
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	switch len(args) {
	case 0:
		b, err := Requests("POST", fmt.Sprintf("%s/start", node.Url), node.Token, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
	case 1:
		b, err := Requests("POST", fmt.Sprintf("%s/start/%s", node.Url, args[0]), node.Token, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(b))
	default:
		b, err := Requests("POST", fmt.Sprintf("%s/start/%s/%s", node.Url, args[0], args[1]), node.Token, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(b))
	}

}

func (node *Node) Status(args ...string) {
	if node.Wg != nil {
		defer node.Wg.Done()
	}

	var s status = make([]*script.ServiceStatus, 0)
	switch len(args) {
	case 0:
		b, err := Requests("POST", fmt.Sprintf("%s/status", node.Url), node.Token, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(b) == 0 {
			break
		}
		err = json.Unmarshal(b, &s)
		if err != nil {
			fmt.Println(err.Error() + " or token error")
			return
		}
	case 1:
		b, err := Requests("POST", fmt.Sprintf("%s/status/%s", node.Url, args[0]), node.Token, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = json.Unmarshal(b, &s)
		if err != nil {
			fmt.Println(err.Error() + " or token error")
			return
		}
	default:
		b, err := Requests("POST", fmt.Sprintf("%s/status/%s/%s", node.Url, args[0], args[1]), node.Token, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = json.Unmarshal(b, &s)
		if err != nil {
			fmt.Println(err.Error() + " or token error")
			return
		}
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
	// default:
	// 	b, err = Requests("POST", fmt.Sprintf("%s/env/%s/%s", node.Url, args[0], args[1]), node.Token, nil)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// }
	l := make(map[string]string, 0)
	// m := make(map[string]string)
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
	if node.Wg != nil {
		defer node.Wg.Done()
	}
	switch len(args) {
	case 0:
		b, err := Requests("POST", fmt.Sprintf("%s/stop", node.Url), node.Token, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
	case 1:
		b, err := Requests("POST", fmt.Sprintf("%s/stop/%s", node.Url, args[0]), node.Token, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(b))
	default:
		b, err := Requests("POST", fmt.Sprintf("%s/stop/%s/%s", node.Url, args[0], args[1]), node.Token, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(b))
	}
}
