package status

import (
	"encoding/json"

	"github.com/hyahm/golog"
)

const (
	STOP        string = "Stop"            // 停止
	RUNNING     string = "Running"         // 运行中
	WAITSTOP    string = "Waiting Stop"    // 等待停止
	WAITRESTART string = "Waiting Restart" // 等待重启
	INSTALL     string = "Installing"      // 正在安装
	STARTING    string = "Starting"        // 正在启动中
)

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
	Always       bool    `json:"-"`
	IsCron       bool    `json:"isCron"`
	RestartCount int     `json:"restartCount"` // 记录失败重启的次数
	Disable      bool    `json:"disable"`      // 是否禁用
	Cpu          float64 `json:"cpu"`
	Mem          uint64  `json:"mem"`
	OS           string  `json:"os"`
	SCSVerion    string  `json:"scs_version"`
}

func (ss *ServiceStatus) Bytes() []byte {
	b, err := json.Marshal(ss)
	if err != nil {
		golog.Error(err)
	}
	return b
}
