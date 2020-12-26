package probe

import (
	"fmt"
	"scs/alert"
	"strconv"
	"time"

	"github.com/hyahm/golog"
	"github.com/shirou/gopsutil/cpu"
)

type Cpu struct {
	Percent  float64
	AI       *alert.AlertInfo
	Interval time.Duration
}

func NewCpu(percent float64, interval, continuityInterval time.Duration) *Cpu {
	return &Cpu{
		Percent: percent,
		AI: &alert.AlertInfo{
			AM:                 &alert.Message{},
			ContinuityInterval: continuityInterval,
		},
		Interval: interval,
	}
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
	c.AI.Recover("cpu恢复")
}
