package logger

import (
	"path/filepath"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
)

type Logger struct {
	Path  string        `yaml:"path"`
	Size  int64         `yaml:"size"`
	Day   bool          `yaml:"day"`
	Clear time.Duration `json:"clear"`
}

func defaultLogger() *Logger {
	return &Logger{
		Path:  "",
		Size:  0,
		Day:   true,
		Clear: 30 * 24 * time.Hour,
	}
}

func ReloadLogger(log *Logger) {
	if log == nil {
		log = defaultLogger()
	}
	global.LogDir = filepath.Dir(log.Path)

	if global.LogDir == "." {
		global.LogDir = "log"
	}
	global.CleanLog = log.Clear

	golog.InitLogger(log.Path, log.Size, log.Day, log.Clear)
	// 设置所有级别的日志都显示
	golog.Level = golog.All
	// 设置 日志名， 如果Cfg.Log.Path为空， 那么输出到控制台
	// golog.Name = "scs.log"
}
