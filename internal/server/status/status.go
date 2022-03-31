package status

const (
	STOP        string = "Stop"            // 停止
	RUNNING     string = "Running"         // 运行中
	WAITSTOP    string = "Waiting Stop"    // 等待停止
	WAITRESTART string = "Waiting Restart" // 等待重启
	INSTALL     string = "Installing"      // 正在安装
	STARTING    string = "Starting"        // 正在启动中
)

type Status struct {
	Pid          int    `json:"pid"`
	Status       string `json:"status"`
	CanNotStop   bool   `json:"cannotStop"` //
	Start        int64  `json:"start"`      // 启动的时间
	Version      string `json:"version"`
	IsCron       bool   `json:"isCron"`
	RestartCount int    `json:"restartCount"` // 记录失败重启的次数
}
