package handle

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"scs/client/node"
	"scs/config"
	"scs/install"
	"scs/internal"
	"scs/script"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
	"github.com/sacOO7/gowebsocket"
	"gopkg.in/yaml.v2"
)

type PackageInfo struct {
	// 包信息
	Name string `json:"name"`
	Info string `json:"info"`
}

func InstallPackage(w http.ResponseWriter, r *http.Request) {
	e := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// 接受额外参数，覆盖本身的变量
	name := xmux.Var(r)["name"]
	// 读取到配置文件后， 第一步，先读取depend， 获取depend里面的环境变量
	// 第二部， 添加自己的环境变量
	// 第三步，执行install.sh 脚本， websocket， 失败就返回， 否则进行第四步
	// 第四步， 生成script， 如果name是空的就简单了， 不用生成script
	// 如果需要生成script， 那么需要把第一步和第二步的环境变量都添加进script的env中
	// 替换 dir, command, env的值
	golog.Info(name)
	// 先获取可以得到包的url
	packageUrl := ""
	ps := make([]*PackageInfo, 0)
	for _, url := range config.Cfg.Repo.Url {
		b, err := node.Requests(http.MethodPost, fmt.Sprintf("%s/search/%s/%s", url, config.Cfg.Repo.Derivative, name), "", nil)
		if err != nil {
			golog.Info(err)
			w.Write([]byte(err.Error()))
			return
		}

		golog.Info(string(b))
		err = json.Unmarshal(b, &ps)
		if err != nil {
			golog.Info(err)
			w.Write([]byte(name + ".yaml config error"))
			return
		}
		if len(ps) > 0 {
			packageUrl = url
			break
		}
	}
	if packageUrl != "" {
		// 获取配置文件信息
		b, err := node.Requests("GET", fmt.Sprintf("%s/install/%s/%s/%s", packageUrl, config.Cfg.Repo.Derivative, name, name+".yaml"), "", nil)
		if err != nil {
			golog.Info(err)
			w.Write([]byte(err.Error()))
			return
		}
		golog.Info(string(b))
		ic := &install.InstallConfig{}
		err = yaml.Unmarshal(b, ic)
		if err != nil {
			golog.Info(err)
			w.Write([]byte(err.Error()))
			return
		}
		for k, v := range e {
			ic.Env[k] = v
		}
		// 读取到配置文件后， 第一步，先读取depend， 获取depend里面的环境变量
		// 自己的环境变量
		ic.GetDependEnv()
		ss := &script.Script{
			Env: ic.Env,
			// exitCode chan int // 如果推出信号是9
		}
		// 执行安装的脚本
		wgetUrl := ""
		if runtime.GOOS != "windows" {
			wgetUrl = packageUrl + "/install/" + config.Cfg.Repo.Derivative + "/" + name + "/" + "install.sh"
		} else {
			wgetUrl = packageUrl + "/install/" + config.Cfg.Repo.Derivative + "/" + name + "/" + "install.bat"
		}
		scriptByte, err := node.Requests(http.MethodGet, wgetUrl, "", nil)
		if err != nil {
			golog.Info(err)
			w.Write([]byte(err.Error()))
			return
		}
		ss.Install(string(scriptByte))
		// 挂载脚本
		if ic.Script != nil {
			// 替换name
			// 替换dir
			// 替换command

			config.Cfg.AddScript(*ic.Script)
		}
		// 执行shell /bin/bash -c "$(curl -fsSL http://download.hyahm.com/scs.sh)"
		// 生成script
	} else {
		w.Write([]byte("repo url error"))
		return
	}

}

func InstallScript(w http.ResponseWriter, r *http.Request) {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New("ws://echo.websocket.org/")

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Recieved message " + message)
	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		log.Println("Recieved binary data ", data)
	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved ping " + data)
	}

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved pong " + data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		return
	}

	socket.Connect()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return
		}
	}

	s := xmux.GetData(r).Data.(*internal.Script)
	s.ContinuityInterval = s.ContinuityInterval * 1000000000
	if s.KillTime == 0 {
		s.KillTime = 1 * time.Second
	} else {
		golog.Info(s.KillTime)
		s.KillTime = s.KillTime * 1000000000
	}
	golog.Infof("%+v", *s)
	if err := config.Cfg.AddScript(*s); err != nil {
		w.Write([]byte(`{"code": 201, "msg": "already exist script"}`))
		return
	}
	w.Write([]byte(`{"code": 200, "msg": "already add script"}`))
	return
}
