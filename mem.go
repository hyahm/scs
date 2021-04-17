package scs

import (
	"fmt"
	"time"

	"github.com/hyahm/golog"
	"github.com/shirou/gopsutil/mem"
)

type Mem struct {
	Percent  float64
	AI       *AlertInfo
	Interval time.Duration
}

func NewMem() *Mem {
	return &Mem{
		Percent: healthDetector.Config.Mem,
		AI: &AlertInfo{
			AM:                 &Message{},
			ContinuityInterval: healthDetector.Config.ContinuityInterval,
		},
		Interval: healthDetector.Config.Interval,
	}
}
func (m *Mem) Update() {
	m.Percent = healthDetector.Config.Mem
	m.Interval = healthDetector.Config.Interval
	m.AI.ContinuityInterval = healthDetector.Config.ContinuityInterval
}
func (m *Mem) Check() {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		golog.Error(err)
		return
	}

	if float64(memInfo.Used)/float64(memInfo.Total)*100 >= m.Percent {
		m.AI.AM.Top = TopMem(1)[0].ToString()
		m.AI.AM.Use = memInfo.Used / 1024 / 1024 / 1024
		m.AI.AM.Total = memInfo.Total / 1024 / 1024 / 1024
		m.AI.BreakDown(fmt.Sprintf("内存繁忙超过%.2f%%", m.Percent))
		return
	}
	m.AI.AM.Use = memInfo.Used / 1024 / 1024 / 1024
	m.AI.Recover("内存恢复")
}
