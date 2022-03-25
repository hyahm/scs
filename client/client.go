package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/api/handle"
	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/config/alert"
	"github.com/hyahm/scs/internal/config/scripts"
	"github.com/sacOO7/gowebsocket"
	"gopkg.in/yaml.v2"
)

var CCfg *ClientConfig

func ReadClientConfig(configfile string) {

	_, err := os.Stat(configfile)
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
		timeout = 5 * time.Second
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout,
	}

}

func (sc *SCSClient) requests(url string, body io.Reader, method ...string) (interface{}, error) {
	httpMethod := http.MethodPost
	if len(method) > 0 {
		httpMethod = method[0]
	}
	req, err := http.NewRequest(httpMethod, sc.Domain+url, body)
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
	err = checkCode(resp.StatusCode)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res := &handle.Response{}
	err = json.Unmarshal(b, res)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = checkCode(res.Code)
	if err != nil {
		return nil, errors.New(res.Msg)
	}
	return res.Data, nil
}

func (sc *SCSClient) requestStatuss(url string, body io.Reader, method ...string) (*controller.StatusList, error) {
	httpMethod := http.MethodPost
	if len(method) > 0 {
		httpMethod = method[0]
	}
	req, err := http.NewRequest(httpMethod, sc.Domain+url, body)
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
	err = checkCode(resp.StatusCode)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res := &controller.StatusList{}
	err = json.Unmarshal(b, res)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = checkCode(res.Code)
	if err != nil {
		return nil, errors.New(res.Msg)
	}
	return res, nil
}

func checkCode(code int) error {
	switch code {
	case 203:
		return ErrToken
	case 500:
		return ErrResponseData
	case 511:
		return ErrStatusNetworkAuthenticationRequired
	case 404:
		return ErrFoundPnameOrName
	case 201:
		return ErrWaitReload
	case 400:
		return ErrHttps
	default:
		return nil
	}
}

func (sc *SCSClient) WebSocket(url string, body io.Reader) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGHUP, syscall.SIGINT)

	wsdomain := sc.Domain + url
	wsdomain = strings.Replace(wsdomain, "https", "wss", 1)
	wsdomain = strings.Replace(wsdomain, "http", "ws", 1)
	socket := gowebsocket.New(wsdomain)

	socket.RequestHeader.Set("Token", sc.Token)
	close := make(chan bool, 1)

	socket.OnConnected = func(socket gowebsocket.Socket) {
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		// log.Println("Recieved message " + message)
		fmt.Println(message)
	}

	// socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
	// 	fmt.Println("Recieved binary data ", data)
	// }
	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		fmt.Println(err)
		close <- true
	}
	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		// log.Println("Recieved ping " + data)
	}

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		// log.Println("Recieved pong " + data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		close <- true
	}
	socket.Connect()
	for {
		select {
		case <-interrupt:
			return
		case <-close:
			return
		}
	}
	// wsdomain := sc.Domain
	// wsdomain = strings.Replace(wsdomain, "http", "ws", 1)
	// req, err := http.NewRequest(http.MethodGet, sc.Domain+url, body)
	// if err != nil {
	// 	return nil, err
	// }
	// key := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().Unix())))
	// req.Header.Set("Upgrade", "websocket")
	// req.Header.Set("Connection", "Upgrade")
	// req.Header.Set("Sec-WebSocket-Version", "13")
	// req.Header.Set("Sec-WebSocket-Key", string(key))
	// resp, err := client(sc.Timeout).Do(req)
	// if err != nil {
	// 	return nil, err
	// }
	// defer resp.Body.Close()
	// switch resp.StatusCode {
	// case 203:
	// 	return nil, ErrToken
	// case 500:
	// 	return nil, ErrResponseData
	// case 511:
	// 	return nil, ErrStatusNetworkAuthenticationRequired
	// case 404:
	// 	return nil, ErrFoundPnameOrName
	// case 201:
	// 	return nil, ErrWaitReload
	// case 400:
	// 	return nil, ErrHttps
	// default:
	// 	return ioutil.ReadAll(resp.Body)
	// }

}

// 标记当前副本不能停止
func (sc *SCSClient) CanNotStop() (interface{}, error) {
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests("/cannotstop/"+sc.Name, nil)
}

// 标记当前副本可以停止
func (sc *SCSClient) CanStop() (interface{}, error) {

	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests("/canstop/"+sc.Name, nil)
}

// 获取当前副本的错误日志
func (sc *SCSClient) Log(line int) {
	if sc.Name == "" {
		golog.Error(ErrNameIsEmpty)
		return
	}
	sc.WebSocket(fmt.Sprintf("/log/%s/%d", sc.Name, line), nil)
}

//  获取当前副本的环境变量
func (sc *SCSClient) Env() (interface{}, error) {
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests("/env/"+sc.Name, nil)
}

//  获取当前副本的环境变量
func (sc *SCSClient) Info() (interface{}, error) {
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests("/server/info/"+sc.Name, nil)
}

// 重新加载配置文件
func (sc *SCSClient) Reload() (interface{}, error) {
	return sc.requests("/-/reload", nil)
}

// 杀掉此脚本及其对应的所有副本
func (sc *SCSClient) KillPname() (interface{}, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/kill/"+sc.Pname, nil)
}

// 杀掉此副本
func (sc *SCSClient) KillName() (interface{}, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests(fmt.Sprintf("/kill/%s/%s", sc.Pname, sc.Name), nil)
}

// 更新所有脚本
func (sc *SCSClient) UpdateAll() (interface{}, error) {
	return sc.requests("/update", nil)
}

// 更新此脚本
func (sc *SCSClient) UpdatePname() (interface{}, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/update/"+sc.Pname, nil)
}

// 更新此副本
func (sc *SCSClient) UpdateName() (interface{}, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests(fmt.Sprintf("/update/%s/%s", sc.Pname, sc.Name), nil)
}

// 重启所有脚本
func (sc *SCSClient) RestartAll() (interface{}, error) {
	return sc.requests("/restart", nil)
}

// 重启此脚本
func (sc *SCSClient) RestartPname() (interface{}, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/restart/"+sc.Pname, nil)
}

// 重启当前副本
func (sc *SCSClient) RestartName() (interface{}, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests(fmt.Sprintf("/restart/%s/%s", sc.Pname, sc.Name), nil)
}

// 启动所有脚本
func (sc *SCSClient) StartAll() (interface{}, error) {
	return sc.requests("/start", nil)
}

// 启动当前脚本
func (sc *SCSClient) StartPname() (interface{}, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/start/"+sc.Pname, nil)
}

// 启动当前副本
func (sc *SCSClient) StartName() (interface{}, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests(fmt.Sprintf("/start/%s/%s", sc.Pname, sc.Name), nil)
}

// 停止所有脚本
func (sc *SCSClient) StopAll() (interface{}, error) {
	return sc.requests("/stop", nil)
}

// 停止当前脚本
func (sc *SCSClient) StopPname() (interface{}, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/stop/"+sc.Pname, nil)
}

// 停止当前副本
func (sc *SCSClient) StopName() (interface{}, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests(fmt.Sprintf("/stop/%s/%s", sc.Pname, sc.Name), nil)
}

// 移除所有脚本
// func (sc *SCSClient) RemoveAllScrip() ([]byte, error) {
// 	return sc.requests("/remove", nil)
// }

// 移除当前脚本
func (sc *SCSClient) RemovePnameScrip() (interface{}, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/remove/"+sc.Pname, nil)
}

// 移除当前副本
func (sc *SCSClient) RemoveNameScrip() (interface{}, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requests(fmt.Sprintf("/remove/%s/%s", sc.Pname, sc.Name), nil)
}

// 启用脚本
func (sc *SCSClient) Enable() (interface{}, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/enable/"+sc.Pname, nil)
}

func (sc *SCSClient) GetServers() (interface{}, error) {
	return sc.requests("/get/servers", nil)
}

func (sc *SCSClient) GetAlarms() (interface{}, error) {
	return sc.requests("/get/alarms", nil)
}

func (sc *SCSClient) GetScripts() (interface{}, error) {
	return sc.requests("/get/scripts", nil)
}

// 禁用脚本
func (sc *SCSClient) Disable() (interface{}, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requests("/disable/"+sc.Pname, nil)
}

func (sc *SCSClient) Repo() (interface{}, error) {
	return sc.requests("/get/repo", nil)
}

func (sc *SCSClient) Search(derivative, serviceName string) (interface{}, error) {
	return sc.requests(fmt.Sprintf("/search/%s/%s", derivative, serviceName), nil)
}

// 添加或修改脚本
func (sc *SCSClient) AddScript(s *scripts.Script) (interface{}, error) {
	send, _ := json.Marshal(s)
	return sc.requests("/script", bytes.NewReader(send))
}

// 获取此所有脚本的状态
func (sc *SCSClient) StatusAll() (*controller.StatusList, error) {
	return sc.requestStatuss("/status", nil)
}

// 获取此脚本的状态
func (sc *SCSClient) StatusPname() (*controller.StatusList, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.requestStatuss("/status/"+sc.Pname, nil)
}

// 获取此副本的状态
func (sc *SCSClient) StatusName() (*controller.StatusList, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.requestStatuss(fmt.Sprintf("/status/%s/%s", sc.Pname, sc.Name), nil)
}

// 检测远程机器的健康探针
func (sc *SCSClient) Probe() (interface{}, error) {
	return sc.requests("/probe", nil)
}

// 发送报警
func (sc *SCSClient) Alert(alert *alert.RespAlert) (interface{}, error) {
	send, _ := json.Marshal(alert)
	return sc.requests("/set/alert", bytes.NewReader(send))
}
