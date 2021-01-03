package probe

import (
	"fmt"
	"scs/alert"
	"time"

	"github.com/hyahm/golog"
	"github.com/shirou/gopsutil/mem"
)

type Mem struct {
	Percent  float64
	AI       *alert.AlertInfo
	Interval time.Duration
}

func NewMem(percent float64, interval, continuityInterval time.Duration) *Mem {
	return &Mem{
		Percent: percent,
		AI: &alert.AlertInfo{
			AM:                 &alert.Message{},
			ContinuityInterval: continuityInterval,
		},
		Interval: interval,
	}
}
func (m *Mem) Update(probe *Probe) {
	m.Percent = probe.Mem
	m.Interval = probe.Interval
	m.AI.ContinuityInterval = probe.ContinuityInterval
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
	m.AI.Recover("内存恢复")
}
