package scs

import (
	"fmt"

	"github.com/shirou/gopsutil/process"
)

// 计算cpu占用前1的进程
type cpuInfo struct {
	name    string
	Pid     int32
	percent float64
}

func (ci *cpuInfo) ToString() string {
	return fmt.Sprintf("Process: %s, Pid: %d, Percent: %.2f%%", ci.name, ci.Pid, ci.percent)
}

func TopCpu(top int) []*cpuInfo {
	procs := make([]*cpuInfo, 0)
	cs, err := process.Processes()
	if err != nil {
		return procs
	}

	// top5 := make([]*process.Process, 0, 5)
	for _, c := range cs {
		p, err := c.CPUPercent()
		if err != nil {
			continue
		}
		name, _ := c.Name()
		pi := &cpuInfo{
			name:    name,
			Pid:     c.Pid,
			percent: p,
		}
		procs = append(procs, pi)
	}

	for i := 0; i < top; i++ {
		max := i
		for j := i; j < len(procs); j++ {
			if procs[j].percent > procs[max].percent {
				max = j
			}
		}
		procs[max], procs[i] = procs[i], procs[max]
	}

	return procs[:top]
}
