package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/pkg/message"
	"github.com/shirou/gopsutil/cpu"
)

// IO 磁盘 IO 等待检测器，通过 iowait 占比判断 IO 负载。
type IO struct {
	Percent  float64
	AI       *AlertInfo
	Interval time.Duration
}

// NewIO 创建 IO 检测器，阈值从健康检测配置中读取。
func NewIO() *IO {
	return &IO{
		Percent: healthDetector.Config.Probe.IO,
		AI: &AlertInfo{
			AM:                 &message.Message{},
			ContinuityInterval: healthDetector.Config.Probe.ContinuityInterval,
		},
		Interval: healthDetector.Config.Probe.Interval,
	}
}

// Update 热更新 IO 检测器配置。
func (io *IO) Update() {
	io.Percent = healthDetector.Config.Probe.IO
	io.Interval = healthDetector.Config.Probe.Interval
	io.AI.ContinuityInterval = healthDetector.Config.Probe.ContinuityInterval
}

// Check 检测所有 CPU 的 iowait 总占比，超过阈值则报警。
func (io *IO) Check() {
	times, err := cpu.Times(true)
	if err != nil {
		golog.Error(err)
		return
	}
	var totalIowait float64
	for _, t := range times {
		totalIowait += t.Iowait
	}
	if totalIowait >= io.Percent {
		io.AI.AM.UsePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", totalIowait), 64)
		io.AI.BreakDown(fmt.Sprintf("IO等待超过 %.2f%%", io.Percent))
		return
	}
	io.AI.AM.UsePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", totalIowait), 64)
	io.AI.Recover("IO等待恢复")
}