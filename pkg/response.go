package pkg

import (
	"encoding/json"
	"fmt"

	"github.com/hyahm/golog"
)

type Response struct {
	Code int         `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
	Role string      `json:"role,omitempty"`
}

func (res *Response) Marshal() []byte {
	b, err := json.Marshal(res)
	if err != nil {
		golog.Error(err)
	}
	return b
}

func (res *Response) Sucess(msg string) []byte {
	res.Code = 200
	res.Msg = msg
	return res.Marshal()
}

func (res *Response) ErrorE(err error) []byte {
	res.Code = 200
	res.Msg = err.Error()
	return res.Marshal()
}

// func (res *Response) NotFound(role string) []byte {
// 	res.Code = 404
// 	res.Role = role
// 	return res.Marshal()
// }

func NotFoundScript() []byte {
	return []byte(fmt.Sprintf(`{"code": 404, "msg": "not found pname or name" }`))
}

func Waiting(step string) []byte {
	return []byte(fmt.Sprintf(`{"code": 200, "msg": "waiting %s"}`, step))
}

func WaitingConfigChanged() []byte {
	return []byte(fmt.Sprintf(`{"code": 200, "msg": "config file is reloading, waiting completed first" }`))
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
