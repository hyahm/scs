package config

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/hyahm/golog"
	"github.com/shirou/gopsutil/disk"
)

var healthDetector *Detector

type Detector struct {
	Probe  *Probe
	Ctx    context.Context
	Cancel context.CancelFunc
	Config *Config
	Cps    []CheckPointer
}

// 保存配置文件信息
type Probe struct {
	Monitor   []string `yaml:"monitor,omitempty"`
	Monitored []string `yaml:"monitored,omitempty"`
	// 内存使用率 默认90
	Mem float64 `yaml:"mem,omitempty"`
	IO  float64 `yaml:"io,omitempty"`
	// cpu使用率, 默认90
	Cpu float64 `yaml:"cpu,omitempty"`
	// 硬盘使用率， 默认90
	Disk          float64  `yaml:"disk,omitempty"`
	ExcludeDisk   []string `yaml:"excludeDisk,omitempty"`
	DiskPartition []string `yaml:"diskPartition,omitempty"`
	// 检测间隔， 默认10秒
	Interval time.Duration `yaml:"interval,omitempty"`
	// 下次报警时间间隔， 如果恢复了就重置
	ContinuityInterval time.Duration `yaml:"continuityInterval,omitempty"`
}

func (p *Probe) InitProbe() {
	if p.Interval == 0 {
		p.Interval = time.Second * 10
	}
	if p.ContinuityInterval == 0 {
		p.ContinuityInterval = time.Hour * 1
	}
	if len(p.DiskPartition) == 0 {
		p.DiskPartition = getDisk(p.ExcludeDisk)
	}
	if p.Mem == 0 {
		p.Mem = 90
	}
	if p.Cpu == 0 {
		p.Cpu = 90
	}
	if p.Disk == 0 {
		p.Disk = 90
	}
}

// InitDetector 初始化全局硬件检测器，绑定配置并创建 context。
func InitDetector(cfg *Config) {
	ctx, cancel := context.WithCancel(context.Background())
	healthDetector = &Detector{
		Probe:  &cfg.Probe,
		Ctx:    ctx,
		Cancel: cancel,
		Config: cfg,
		Cps:    make([]CheckPointer, 4),
	}
}

// StopDetector 停止硬件检测循环。
func StopDetector() {
	if healthDetector != nil && healthDetector.Cancel != nil {
		healthDetector.Cancel()
	}
}

// CheckHardWare 主循环：按 Interval 周期调用各检测器的 Check。
func CheckHardWare() {
	if healthDetector == nil {
		golog.Error("healthDetector 未初始化，无法启动硬件检测")
		return
	}
	// 初始化检测器
	if healthDetector.Config.Probe.Cpu > 0 {
		healthDetector.Cps[0] = NewCpu()
	}
	if healthDetector.Config.Probe.Mem > 0 {
		healthDetector.Cps[1] = NewMem()
	}
	if healthDetector.Config.Probe.Disk > 0 {
		healthDetector.Cps[2] = NewDisk()
	}
	if len(healthDetector.Config.Probe.Monitor) > 0 {
		healthDetector.Cps[3] = NewMonitor()
	}

	golog.Info("硬件检测启动，间隔: ", healthDetector.Config.Probe.Interval)
	for {
		select {
		case <-healthDetector.Ctx.Done():
			golog.Info("硬件检测退出")
			return
		case <-time.After(healthDetector.Config.Probe.Interval):
			for _, check := range healthDetector.Cps {
				if IsNil(check) {
					continue
				}
				check.Check()
			}
		}
	}
}

// ReloadDetector 重载时更新各检测器配置。
func ReloadDetector(cfg *Config) {
	if healthDetector == nil {
		InitDetector(cfg)
		return
	}
	healthDetector.Config = cfg
	healthDetector.Probe = &cfg.Probe
	// 更新已有检测器
	for _, check := range healthDetector.Cps {
		if !IsNil(check) {
			check.Update()
		}
	}
	// 按新配置启用/禁用检测器
	if cfg.Probe.Cpu > 0 && IsNil(healthDetector.Cps[0]) {
		healthDetector.Cps[0] = NewCpu()
	} else if cfg.Probe.Cpu == 0 {
		healthDetector.Cps[0] = nil
	}
	if cfg.Probe.Mem > 0 && IsNil(healthDetector.Cps[1]) {
		healthDetector.Cps[1] = NewMem()
	} else if cfg.Probe.Mem == 0 {
		healthDetector.Cps[1] = nil
	}
	if cfg.Probe.Disk > 0 && IsNil(healthDetector.Cps[2]) {
		healthDetector.Cps[2] = NewDisk()
	} else if cfg.Probe.Disk == 0 {
		healthDetector.Cps[2] = nil
	}
	if len(cfg.Probe.Monitor) > 0 && IsNil(healthDetector.Cps[3]) {
		healthDetector.Cps[3] = NewMonitor()
	} else if len(cfg.Probe.Monitor) == 0 {
		healthDetector.Cps[3] = nil
	}
}

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	return !vi.IsValid() || vi.IsNil()
}

func getDisk(excludeDisk []string) []string {
	dp := make([]disk.PartitionStat, 0)
	parts, err := disk.Partitions(true)
	if err != nil {
		golog.Error(err)
		return []string{}
	}
	excludePath := make(map[string]int)
	for _, he := range excludeDisk {
		excludePath[strings.ToUpper(he)] = 0
	}

	for _, part := range parts {
		if _, ok := excludePath[strings.ToUpper(part.Mountpoint)]; ok {
			continue
		}
		if _, ok := cludeType[strings.ToUpper(part.Fstype)]; ok {
			dp = append(dp, part)
		}
	}
	list := make([]string, 0, len(dp))
	for _, part := range dp {
		list = append(list, part.Mountpoint)
		golog.Infof("alert disk: --%s--, type: %s", part.Mountpoint, part.Fstype)
	}
	return list
}
