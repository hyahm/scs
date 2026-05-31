package config

import (
	"github.com/hyahm/golog"
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

func (log *Logger) initLogger() {
	if log == nil {
		log = defaultLogger()
	}

	if log.Path == "" {
		log.Path = "log"
	}
	if Cfg.Debug {
		golog.SetLevel(golog.DEBUG)
		return
	}
	golog.SetDir(log.Path)
	golog.InitLogger("scs.log", 0, true)

}
