package logger

import (
	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
)

type Logger struct {
	Path  string `yaml:"path"`
	Size  int64  `yaml:"size"`
	Day   bool   `yaml:"day"`
	Clear int    `json:"clear"`
}

func defaultLogger() *Logger {
	return &Logger{
		Path:  "log",
		Size:  0,
		Day:   true,
		Clear: 7,
	}
}

func ReloadLogger(log *Logger) {
	if log == nil {
		log = defaultLogger()
	}

	global.CS.LogDir = log.Path
	global.CS.CleanLog = log.Clear
	golog.InitLogger(global.CS.LogDir, log.Size, log.Day, global.CS.CleanLog)
	// 设置所有级别的日志都显示
	// golog.Level = golog.ALL
	// 设置 日志名， 如果Cfg.Log.Path为空， 那么输出到控制台
	// golog.Name = "scs.log"
}
