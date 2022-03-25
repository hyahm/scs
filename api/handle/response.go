package handle

import (
	"encoding/json"

	"github.com/hyahm/golog"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (res *Response) Marshal() []byte {
	b, err := json.Marshal(res)
	if err != nil {
		golog.Error(err)
	}
	return b
}

func (res *Response) Sucess() []byte {
	res.Code = 200
	return res.Marshal()
}

func (res *Response) NotFound() []byte {
	res.Code = 404
	return res.Marshal()
}
