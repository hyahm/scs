package logger

import (
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
)

type Logger struct {
	Path  string `yaml:"path"`
	Clear int    `json:"clear"`
}

func defaultLogger() *Logger {
	return &Logger{
		Path:  "log",
		Clear: 7,
	}
}

func ReloadLogger(log *Logger) {
	if log == nil {
		log = defaultLogger()
	}
	logdir := "log"
	global.CS.LogDir = log.Path
	global.CS.CleanLog = log.Clear
	if global.CS.LogDir != "" {
		logdir = global.CS.LogDir
	}
	golog.SetDir(logdir)
	if global.CS.CleanLog > 0 {
		golog.SetExpireDuration(time.Duration(global.CS.CleanLog) * time.Hour * 24)
	}

	if log.Path != "" {
		golog.SetDir(log.Path)
	}
	golog.InitLogger("scs.log", 0, true)

	// 设置所有级别的日志都显示
	// golog.Level = golog.ALL
	// 设置 日志名， 如果Cfg.Log.Path为空， 那么输出到控制台
	// golog.Name = "scs.log"
}
