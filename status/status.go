package status

const (
	STOP        string = "Stop"
	RUNNING     string = "Running"
	WAITSTOP    string = "Waiting Stop"
	WAITRESTART string = "Waiting Restart"
	INSTALL     string = "Installing"
)

type ServiceStatus struct {
	PName      string `json:"pname"`
	Name       string `json:"name"`
	Pid        int    `json:"ppid"`
	Status     string `json:"status"`
	Command    string `json:"command"`
	Path       string `json:"path"`
	CanNotStop bool   `json:"cannotStop"` //
	// Stoping      bool   `json:"stoping,omitempty"`
	Start        int64  `json:"start"` // 启动的时间
	Version      string `json:"version"`
	Always       bool
	RestartCount int     `json:"restartCount"` // 记录失败重启的次数
	Disable      bool    `json:"disable"`      // 是否禁用
	Up           string  `json:"-"`
	Cpu          float64 `json:"cpu"`
	Mem          uint64  `json:"mem"`
}
