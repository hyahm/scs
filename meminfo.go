package scs

import (
	"fmt"

	"github.com/shirou/gopsutil/process"
)

// 计算 内存占用前1的进程
type memInfo struct {
	name    string
	Pid     int32
	percent float32
}

func (mi *memInfo) ToString() string {
	return fmt.Sprintf("Process: %s, Pid: %d, Percent: %.2f%%", mi.name, mi.Pid, mi.percent)
}

func TopMem(top int) []*memInfo {
	procs := make([]*memInfo, 0)
	cs, err := process.Processes()
	if err != nil {
		return procs
	}
	for _, c := range cs {
		p, err := c.MemoryPercent()
		if err != nil {
			continue
		}
		name, _ := c.Name()
		pi := &memInfo{
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
