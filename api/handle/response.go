package handle

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

func NotFoundScript(role string) []byte {
	return []byte(fmt.Sprintf(`{"code": 404, "msg": "not found this script", "role": "%s"}`, role))
}

func Waiting(step, role string) []byte {
	return []byte(fmt.Sprintf(`{"code": 200, "msg": "waiting %s", "role": "%s"}`, step, role))
}

func WaitingConfigChanged(role string) []byte {
	return []byte(fmt.Sprintf(`{"code": 200, "msg": "config file is reloading, waiting completed first", "role": "%s"}`, role))
}
