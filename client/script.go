package client

import (
	"time"

	"github.com/hyahm/scs/client/alert"
)

// 配置文件的数据
type Script struct {
	Name               string            `yaml:"name,omitempty" json:"name"`
	Dir                string            `yaml:"dir,omitempty" json:"dir"`
	Command            string            `yaml:"command,omitempty" json:"command"`
	Replicate          int               `yaml:"replicate,omitempty" json:"replicate,omitempty"`
	Always             bool              `yaml:"always,omitempty" json:"always,omitempty"`
	DisableAlert       bool              `yaml:"disableAlert,omitempty" json:"disableAlert,omitempty"`
	Env                map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
	ContinuityInterval time.Duration     `yaml:"continuityInterval,omitempty" json:"continuityInterval,omitempty"`
	Port               int               `yaml:"port,omitempty" json:"port,omitempty"`
	AT                 *alert.AlertTo    `yaml:"alert,omitempty" json:"alert,omitempty"`
	Version            string            `yaml:"version,omitempty" json:"version,omitempty"`
	Loop               int               `yaml:"loop,omitempty" json:"loop,omitempty"`
	LookPath           []*LoopPath       `yaml:"lookPath,omitempty" json:"lookPath,omitempty"`
	Disable            bool              `yaml:"disable,omitempty" json:"disable,omitempty"`
	Cron               *Cron             `yaml:"cron,omitempty" json:"cron,omitempty"`
	Update             string            `yaml:"update,omitempty" json:"update,omitempty"`
	DeleteWhenExit     bool              `yaml:"deleteWhenExit,omitempty" json:"deleteWhenExit,omitempty"`
}

type LoopPath struct {
	Path    string `yaml:"path,omitempty" json:"path,omitempty"`
	Command string `yaml:"command,omitempty" json:"command,omitempty"`
	Install string `yaml:"install,omitempty" json:"install,omitempty"`
}

type Cron struct {
	// 开始执行的时间戳
	Start string `yaml:"start,omitempty" json:"start,omitempty"`
	// 间隔的时间， 如果IsMonth 为true， loop 单位为月， 否则为秒
	IsMonth bool `yaml:"isMonth,omitempty" json:"isMonth,omitempty"`
	Loop    int  `yaml:"loop,omitempty" json:"loop,omitempty"`
}

// 优先执行的代码
