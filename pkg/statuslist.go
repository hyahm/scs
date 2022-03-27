package pkg

import (
	"encoding/json"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/status"
)

type StatusList struct {
	Data    []*status.ServiceStatus `json:"data"`
	Code    int                     `json:"code"`
	Msg     string                  `json:"msg"`
	Version string                  `json:"version"`
	Role    string                  `json:"role"`
}

func (sl *StatusList) Marshal() []byte {
	sl.Code = 200
	b, err := json.Marshal(sl)
	if err != nil {
		golog.Error(err)
	}
	return b
}
