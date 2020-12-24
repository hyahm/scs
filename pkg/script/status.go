package script

const (
	STOP        string = "Stop"
	RUNNING     string = "Running"
	WAITSTOP    string = "Waiting Stop"
	WAITRESTART string = "Waiting Restart"
)

type ServiceStatus struct {
	Name         string `json:"name"`
	Pid          int    `json:"ppid"`
	Status       string `json:"status"`
	Command      string `json:"command"`
	PName        string `json:"pname"`
	Path         string `json:"path"`
	CanNotStop   bool   `json:"cannotStop"` //
	Stoping      bool   `json:"stoping,omitempty"`
	Up           int64
	Start        string `json:"start"` // 启动的时间
	Version      string `json:"version"`
	Always       bool
	RestartCount int `json:"restartCount"` // 记录失败重启的次数
}
