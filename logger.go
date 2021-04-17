package scs

import (
	"time"

	"github.com/hyahm/golog"
)

type Logger struct {
	Path  string        `yaml:"path"`
	Size  int64         `yaml:"size"`
	Day   bool          `yaml:"day"`
	Clear time.Duration `json:"clear"`
}

func ReloadLogger(log *Logger) {
	clear := log.Clear
	if log.Clear == 0 {
		clear = 7 * time.Hour * 24
	}
	golog.InitLogger(log.Path, log.Size, log.Day, clear)
	// 设置所有级别的日志都显示
	golog.Level = golog.All
	// 设置 日志名， 如果Cfg.Log.Path为空， 那么输出到控制台
	golog.Name = "scs.log"
}
