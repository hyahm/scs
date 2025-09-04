package pkg

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

var ErrBugMsg = "严重错误， 请提交问题到https://github.com/hyahm/scs"

type Response struct {
	Code    int         `json:"code,omitempty"`
	Msg     string      `json:"msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Version string      `json:"version,omitempty"`
}

func (res *Response) Marshal() []byte {
	b, err := json.Marshal(res)
	if err != nil {
		golog.Error(err)
	}
	return b
}

func Sucess(r *http.Request, data interface{}) {
	xmux.GetInstance(r).Data.(*Response).Data = data
}

func Error(r *http.Request, msg string) {
	xmux.GetInstance(r).Data.(*Response).Code = 500
	xmux.GetInstance(r).Data.(*Response).Msg = msg
}

var ErrNotFound = errors.New("not found pname or name")

// var ResponseMsg map[int]string

func init() {
	// ResponseMsg = make(map[int]string)
	// ResponseMsg[200] = "ok"
	// ResponseMsg[201] = "config file is reloading, waiting completed first"
	// ResponseMsg[404] = "not found pname or name"
	// ResponseMsg[405] = "auth failed"
	// ResponseMsg[406] = "没有找到对应运行的信号参数"
}

// 这是返回给前端的数据结构
type ServiceStatus struct {
	PName        string  `json:"pname"`
	Name         string  `json:"name"`
	Pid          int     `json:"pid"`
	Status       string  `json:"status"`
	Command      string  `json:"command"`
	Path         string  `json:"path"`
	CanNotStop   bool    `json:"cannotStop"` //
	Start        int64   `json:"start"`      // 启动的时间
	Version      string  `json:"version"`
	IsCron       bool    `json:"isCron"`
	RestartCount int     `json:"restartCount"` // 记录失败重启的次数
	Disable      bool    `json:"disable"`      // 是否禁用
	Cpu          float64 `json:"cpu"`
	Mem          uint64  `json:"mem"`
	OS           string  `json:"os"`
}

func (ss *ServiceStatus) Bytes() []byte {
	b, err := json.Marshal(ss)
	if err != nil {
		golog.Error(err)
	}
	return b
}
