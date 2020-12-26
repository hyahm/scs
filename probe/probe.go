package probe

import (
	"scs/global"
	"strings"
	"time"

	"github.com/hyahm/golog"
	"github.com/shirou/gopsutil/disk"
)

var Exit chan struct{}

//  保存配置文件信息
type Probe struct {
	Monitor   []string `yaml:"monitor"`
	Monitored []string `yaml:"monitored"`
	// 内存使用率 默认90
	Mem float64 `yaml:"mem"`
	// cpu使用率, 默认90
	Cpu float64 `yaml:"cpu"`
	// 硬盘使用率， 默认90
	Disk        float64  `yaml:"disk"`
	ExcludeDisk []string `yaml:"excludeDisk"`
	// 检测间隔， 默认10秒
	Interval time.Duration `yaml:"interval"`
	// 下次报警时间间隔， 如果恢复了就重置
	ContinuityInterval time.Duration `yaml:"continuityInterval"`
}

func (probe *Probe) InitHWAlert() {
	Exit = make(chan struct{}, 2)
	if probe.Interval == 0 {
		probe.Interval = time.Second * 10
	}
	if probe.ContinuityInterval == 0 {
		probe.ContinuityInterval = time.Hour * 1
	}
	if len(probe.Monitored) > 0 {
		global.Monitored = probe.Monitored
	} else {
		global.Monitored = make([]string, 0)
	}
	if probe.Cpu == 0 {
		probe.Cpu = 90
	}
	if probe.Mem == 0 {
		probe.Mem = 90
	}
	if probe.Disk == 0 {
		probe.Disk = 85
	}
	if probe.Cpu > 0 || probe.Mem > 0 || probe.Disk > 0 || len(probe.Monitor) > 0 {
		go probe.CheckHardWare()
	}

}

func (probe *Probe) CheckHardWare() {
	cps := make([]CheckPointer, 0)
	if probe.Cpu > 0 {
		cps = append(cps, NewCpu(probe.Cpu, probe.Interval, probe.ContinuityInterval))
	}
	if probe.Mem > 0 {
		cps = append(cps, NewMem(probe.Mem, probe.Interval, probe.ContinuityInterval))
	}
	if probe.Disk > 0 {
		cps = append(cps, NewDisk(probe.Disk, probe.getDisk(), probe.Interval, probe.ContinuityInterval))
	}

	if len(probe.Monitor) > 0 {
		cps = append(cps, NewMonitor(probe.Monitor, probe.Interval, probe.ContinuityInterval))
	}
	for {
		select {
		case <-Exit:
			golog.Info("exit check")
			return
		case <-time.After(probe.Interval):
			for _, check := range cps {
				check.Check()
			}
		}
	}

}

func (probe *Probe) getDisk() []disk.PartitionStat {
	dp := make([]disk.PartitionStat, 0)
	parts, err := disk.Partitions(true)
	if err != nil {
		golog.Error(err)
		return dp
	}
	excludePath := make(map[string]int)
	for _, he := range probe.ExcludeDisk {
		excludePath[strings.ToUpper(he)] = 0
	}

	mountNames := make(map[string]string)
	for _, part := range parts {
		if _, ok := excludePath[strings.ToUpper(part.Mountpoint)]; ok {
			continue
		}

		if _, ok := cludeType[strings.ToUpper(part.Fstype)]; ok {
			mountNames[part.Mountpoint] = part.Fstype
			dp = append(dp, part)
			continue
		}

	}
	for _, part := range dp {
		golog.Infof("alert dist: --%s--, type: %s", part.Mountpoint, part.Fstype)
	}
	return dp
}
