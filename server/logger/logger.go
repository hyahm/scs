package logger

import "github.com/hyahm/golog"

var lg *Logger

type Logger struct {
	Path string `yaml:"path"`
	Size int64  `yaml:"size"`
	Day  bool   `yaml:"day"`
}

// 如果Cfg.Log.Path为空， 那么输出到控制台
func Run(logger *Logger) {
	lg = logger
	golog.InitLogger(lg.Path, lg.Size, lg.Day)
	// 设置所有级别的日志都显示
	golog.Level = golog.All
	// 设置 日志名
	golog.Name = "scs.log"
}
