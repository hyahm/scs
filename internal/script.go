package internal

import (
	"time"
)

// 配置文件的数据
type Script struct {
	Name               string            `yaml:"name,omitempty" json:"name"`
	Dir                string            `yaml:"dir,omitempty" json:"dir"`
	Command            string            `yaml:"command,omitempty" json:"command"`
	Replicate          int               `yaml:"replicate,omitempty" json:"replicate"`
	Always             bool              `yaml:"always,omitempty" json:"always"`
	DisableAlert       bool              `yaml:"disableAlert,omitempty" json:"disableAlert"`
	Env                map[string]string `yaml:"env,omitempty" json:"env"`
	ContinuityInterval time.Duration     `yaml:"continuityInterval,omitempty" json:"continuityInterval"`
	Port               int               `yaml:"port,omitempty" json:"port"`
	AT                 *AlertTo          `yaml:"alert,omitempty" json:"alert"`
	Version            string            `yaml:"version,omitempty" json:"version"`
	Loop               int               `yaml:"loop,omitempty" json:"loop"`
	LookPath           []*LoopPath       `yaml:"lookPath,omitempty" json:"loopPath"`
	Disable            bool              `yaml:"disable,omitempty" json:"disable"`
	Cron               *Cron             `yaml:"cron,omitempty" json:"cron"`
}

type LoopPath struct {
	Path    string `yaml:"path,omitempty" json:"path"`
	Install string `yaml:"install,omitempty" json:"install"`
}

type Cron struct {
	// 开始执行的时间戳
	Start string `yaml:"start,omitempty" json:"start"`
	// 间隔的时间， 如果IsMonth 为true， loop 单位为月， 否则为秒
	IsMonth bool `yaml:"isMonth,omitempty" json:"isMonth"`
	Loop    int  `yaml:"loop,omitempty" json:"loop"`
}

// 优先执行的代码
