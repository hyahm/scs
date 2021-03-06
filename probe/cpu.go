package probe

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/alert"
	"github.com/hyahm/scs/message"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
)

type Cpu struct {
	Percent  float64
	AI       *alert.AlertInfo
	Interval time.Duration
}

func GetProcessInfo(pid int32) (float64, uint64, error) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return 0, 0, err
	}
	ci, err := p.CPUPercent()
	if err != nil {
		return 0, 0, err
	}
	m, err := p.MemoryInfo()
	if err != nil {
		return ci, 0, err
	}
	return ci, m.RSS, err

}

func NewCpu() *Cpu {
	return &Cpu{
		Percent: healthDetector.Config.Cpu,
		AI: &alert.AlertInfo{
			AM:                 &message.Message{},
			ContinuityInterval: healthDetector.Config.ContinuityInterval,
		},
		Interval: healthDetector.Config.Interval,
	}
}

func (c *Cpu) Update() {
	c.Percent = healthDetector.Config.Cpu
	c.Interval = healthDetector.Config.Interval
	c.AI.ContinuityInterval = healthDetector.Config.ContinuityInterval
}

func (c *Cpu) Check() {
	percents, err := cpu.Percent(time.Second*1, true)
	if err != nil {
		golog.Error(err)
		return
	}
	var totalPercents float64
	for _, percent := range percents {
		totalPercents += percent
	}
	if totalPercents >= c.Percent*(float64)(len(percents)) {
		c.AI.AM.UsePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", totalPercents), 64)
		c.AI.AM.Top = TopCpu(1)[0].ToString()
		c.AI.BreakDown(fmt.Sprintf("cpu使用率超过 %.2f%%", c.Percent))
		return
	}
	c.AI.AM.UsePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", totalPercents), 64)
	c.AI.Recover("cpu恢复")
}
