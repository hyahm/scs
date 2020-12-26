package probe

import (
	"time"

	"github.com/shirou/gopsutil/disk"
)

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
	dp       []disk.PartitionStat
	// 下次报警时间间隔， 如果恢复了就重置
	ContinuityInterval time.Duration `yaml:"continuityInterval"`
}

func (probe *Probe) InitHWAlert() {

	if probe.Interval == 0 {
		probe.Interval = time.Second * 10
	}
	if probe.ContinuityInterval == 0 {
		probe.ContinuityInterval = time.Hour * 1
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
	GlobalProbe.Probe = probe
	if probe.Cpu > 0 || probe.Mem > 0 || probe.Disk > 0 || len(probe.Monitor) > 0 {
		go CheckHardWare()
	}

}
