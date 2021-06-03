package cron

import (
	"time"
)

type Cron struct {
	// 开始执行的时间戳
	First     chan bool     `yaml:"-" json:"-"` // 是否是start等于空时的第一次启动
	LoopTime  time.Duration `yaml:"-" json:"-"` // 循环的时间
	Start     string        `yaml:"start,omitempty" json:"start,omitempty"`
	StartTime time.Time     `yaml:"-" json:"-"` // 下次启动的时间
	// 间隔的时间， 如果IsMonth 为true， loop 单位为月， 否则为秒
	IsMonth bool `yaml:"isMonth,omitempty" json:"isMonth,omitempty"`
	Loop    int  `yaml:"loop,omitempty" json:"loop,omitempty"`
}

// 比较cron的配置是否相等， 如果
func (c *Cron) IsEqual(newc *Cron) bool {
	if c == nil && newc == nil {
		return true
	}
	if (c == nil && newc != nil) || (c != nil && newc == nil) {
		return false
	}

	if c.Start != newc.Start ||
		c.IsMonth != newc.IsMonth ||
		c.Loop != newc.Loop {
		return false
	}
	return true
}

// 计算下次启动时间
func (c *Cron) ComputerStartTime() {
	c.LoopTime = time.Duration(c.Loop) * time.Second
	c.StartTime = time.Now().Add(c.LoopTime)
	if c.IsMonth {
		c.StartTime = time.Now().AddDate(0, c.Loop, 0)
	}
}

func CompareCron(c1, c2 *Cron) bool {
	if c1 == nil && c2 != nil || c1 != nil && c2 == nil {
		return false
	}
	if c1 == nil && c2 == nil {
		return true
	}
	if c1.Start != c2.Start ||
		c1.IsMonth != c2.IsMonth ||
		c1.Loop != c2.Loop {
		return false
	}
	return true
}
