package script

import (
	"time"

	"github.com/hyahm/scs/server/alert"
	"github.com/hyahm/scs/server/at"
	"github.com/hyahm/scs/server/cron"
	"github.com/hyahm/scs/server/lookpath"
)

// 配置文件的数据
type Script struct {
	Name               string               `yaml:"name,omitempty" json:"name"`
	Dir                string               `yaml:"dir,omitempty" json:"dir"`
	Command            string               `yaml:"command,omitempty" json:"command"`
	Replicate          int                  `yaml:"replicate,omitempty" json:"replicate,omitempty"`
	Always             bool                 `yaml:"always,omitempty" json:"always,omitempty"`
	DisableAlert       bool                 `yaml:"disableAlert,omitempty" json:"disableAlert,omitempty"`
	Env                map[string]string    `yaml:"env,omitempty" json:"env,omitempty"`
	ContinuityInterval time.Duration        `yaml:"continuityInterval,omitempty" json:"continuityInterval,omitempty"`
	Port               int                  `yaml:"port,omitempty" json:"port,omitempty"`
	AT                 *at.AlertTo          `yaml:"alert,omitempty" json:"alert,omitempty"`
	Version            string               `yaml:"version,omitempty" json:"version,omitempty"`
	LookPath           []*lookpath.LoopPath `yaml:"lookPath,omitempty" json:"lookPath,omitempty"`
	Disable            bool                 `yaml:"disable,omitempty" json:"disable,omitempty"`
	Cron               *cron.Cron           `yaml:"cron,omitempty" json:"cron,omitempty"`
	Update             string               `yaml:"update,omitempty" json:"update,omitempty"`
	DeleteWhenExit     bool                 `yaml:"deleteWhenExit,omitempty" json:"deleteWhenExit,omitempty"`
}

// 优先执行的代码

func (s *Server) successAlert() {
	// 启动成功后恢复的通知
	if !s.AI.Broken {
		return
	}
	for {
		select {
		// 每3秒一次操作
		case <-time.After(time.Second * 3):
			am := &alert.Message{
				Title:      "service recover",
				Pname:      s.Name,
				Name:       s.SubName,
				BrokenTime: s.AI.Start.String(),
				FixTime:    time.Now().String(),
			}
			alert.AlertMessage(am, s.AT)
			s.AI.Broken = false
			return
		case <-s.Ctx.Done():
			return
		}
	}

}
