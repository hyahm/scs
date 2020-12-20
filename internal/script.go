package internal

import (
	"time"
)

// 用户操作的数据
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
	AT                 AlertTo           `yaml:"alert,omitempty" json:"alert"`
	KillTime           time.Duration     `yaml:"killTime,omitempty" json:"killTime"`
	Version            string            `yaml:"version,omitempty" json:"version"`
	LookPath           []*LoopPath       `yaml:"lookPath,omitempty" json:"loopPath"`
}

type LoopPath struct {
	Command string `yaml:"command,omitempty" json:"command"`
	Install string `yaml:"install,omitempty" json:"install"`
}

// 优先执行的代码
