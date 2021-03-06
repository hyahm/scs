package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/hyahm/scs/alert"
	"github.com/hyahm/scs/server"
	"gopkg.in/yaml.v2"
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
}

type SCSClient struct {
	Domain  string
	Token   string
	Pname   string
	Name    string
	Timeout time.Duration
}

func NewClient(timeout ...time.Duration) *SCSClient {
	var rto time.Duration
	if len(timeout) > 0 {
		rto = timeout[0]
	}

	return &SCSClient{
		Domain:  "https://127.0.0.1:11111",
		Token:   os.Getenv("TOKEN"),
		Pname:   os.Getenv("PNAME"),
		Name:    os.Getenv("NAME"),
		Timeout: rto,
	}
}

func client(timeout time.Duration) *http.Client {
	if timeout == 0 {
		timeout = 3 * time.Second
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout,
	}

}

func (sc *SCSClient) requests(url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, sc.Domain+url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Token", sc.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client(sc.Timeout).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 203:
		return nil, ErrToken
	case 500:
		return nil, ErrResponseData
	case 511:
		return nil, ErrStatusNetworkAuthenticationRequired
	case 404:
		return nil, ErrFoundPnameOrName
	case 201:
		return nil, ErrWaitReload
	case 400:
		return nil, ErrHttps
	default:
		return ioutil.ReadAll(resp.Body)
	}

}

// 标记当前副本不能停止
func (sc *SCSClient) CanNotStop() ([]byte, error) {
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests("/cannotstop/"+sc.Name, nil)
}

// 标记当前副本可以停止
func (sc *SCSClient) CanStop() ([]byte, error) {

	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests("/canstop/"+sc.Name, nil)
}

// 获取当前副本的错误日志
func (sc *SCSClient) Log() ([]byte, error) {
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests("/log/"+sc.Name, nil)
}

//  获取当前副本的环境变量
func (sc *SCSClient) Env() ([]byte, error) {
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests("/env/"+sc.Name, nil)
}

// 重新加载配置文件
func (sc *SCSClient) Reload() ([]byte, error) {
	return sc.requests("/-/reload", nil)
}

// 杀掉此脚本及其对应的所有副本
func (sc *SCSClient) KillPname() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/kill/"+sc.Pname, nil)
}

// 杀掉此副本
func (sc *SCSClient) KillName() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests(fmt.Sprintf("/kill/%s/%s", sc.Pname, sc.Name), nil)
}

// 更新所有脚本
func (sc *SCSClient) UpdateAll() ([]byte, error) {
	return sc.requests("/update", nil)
}

// 更新此脚本
func (sc *SCSClient) UpdatePname() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/update/"+sc.Pname, nil)
}

// 更新此副本
func (sc *SCSClient) UpdateName() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests(fmt.Sprintf("/update/%s/%s", sc.Pname, sc.Name), nil)
}

// 重启所有脚本
func (sc *SCSClient) RestartAll() ([]byte, error) {
	return sc.requests("/restart", nil)
}

// 重启此脚本
func (sc *SCSClient) RestartPname() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/restart/"+sc.Pname, nil)
}

// 重启当前副本
func (sc *SCSClient) RestartName() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests(fmt.Sprintf("/restart/%s/%s", sc.Pname, sc.Name), nil)
}

// 启动所有脚本
func (sc *SCSClient) StartAll() ([]byte, error) {
	return sc.requests("/start", nil)
}

// 启动当前脚本
func (sc *SCSClient) StartPname() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/start/"+sc.Pname, nil)
}

// 启动当前副本
func (sc *SCSClient) StartName() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests(fmt.Sprintf("/start/%s/%s", sc.Pname, sc.Name), nil)
}

// 停止所有脚本
func (sc *SCSClient) StopAll() ([]byte, error) {
	return sc.requests("/stop", nil)
}

// 停止当前脚本
func (sc *SCSClient) StopPname() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/stop/"+sc.Pname, nil)
}

// 停止当前副本
func (sc *SCSClient) StopName() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests(fmt.Sprintf("/stop/%s/%s", sc.Pname, sc.Name), nil)
}

// 移除所有脚本
func (sc *SCSClient) RemoveAllScrip() ([]byte, error) {
	return sc.requests("/remove", nil)
}

// 移除当前脚本
func (sc *SCSClient) RemovePnameScrip() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/remove/"+sc.Pname, nil)
}

// 移除当前副本
func (sc *SCSClient) RemoveNameScrip() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests(fmt.Sprintf("/remove/%s/%s", sc.Pname, sc.Name), nil)
}

// 启用脚本
func (sc *SCSClient) Enable() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/enable/"+sc.Pname, nil)
}

func (sc *SCSClient) GetServers() ([]byte, error) {
	return sc.requests("/debug/servers", nil)
}

func (sc *SCSClient) GetScripts() ([]byte, error) {
	return sc.requests("/debug/scripts", nil)
}

// 禁用脚本
func (sc *SCSClient) Disable() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/disable/"+sc.Pname, nil)
}

func (sc *SCSClient) Repo() ([]byte, error) {
	return sc.requests("/get/repo", nil)
}

func (sc *SCSClient) Search(derivative, serviceName string) ([]byte, error) {
	return sc.requests(fmt.Sprintf("/search/%s/%s", derivative, serviceName), nil)
}

// 添加或修改脚本
func (sc *SCSClient) AddScript(s *server.Script) ([]byte, error) {
	send, _ := json.Marshal(s)
	return sc.requests("/script", bytes.NewReader(send))
}

// 获取此所有脚本的状态
func (sc *SCSClient) StatusAll() ([]byte, error) {
	return sc.requests("/status", nil)
}

// 获取此脚本的状态
func (sc *SCSClient) StatusPname() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/status/"+sc.Pname, nil)
}

// 获取此副本的状态
func (sc *SCSClient) StatusName() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests(fmt.Sprintf("/status/%s/%s", sc.Pname, sc.Name), nil)
}

// 检测远程机器的健康探针
func (sc *SCSClient) Probe() ([]byte, error) {
	return sc.requests("/probe", nil)
}

// 发送报警
func (sc *SCSClient) Alert(alert *alert.RespAlert) ([]byte, error) {
	send, _ := json.Marshal(alert)
	return sc.requests("/set/alert", bytes.NewReader(send))
}
