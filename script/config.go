package script

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hyahm/scs/alert"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/logger"
	"github.com/hyahm/scs/probe"

	"github.com/hyahm/golog"
	"gopkg.in/yaml.v2"
)

type Repo struct {
	Url        []string `yaml:"url"`
	Derivative string   `yaml:"derivative"`
}

type config struct {
	Listen      string            `yaml:"listen"`
	Token       string            `yaml:"token"`
	Log         logger.Logger     `yaml:"log"`
	LogCount    int               `yaml:"logCount"`
	IgnoreToken []string          `yaml:"ignoreToken"`
	Repo        *Repo             `yaml:"repo"`
	Alert       alert.Alert       `yaml:"alert"`
	Probe       probe.Probe       `yaml:"probe"`
	SC          []internal.Script `yaml:"scripts"`
}

// 保存的配置文件路径
var cfgfile string

// 保存的全局的配置
var Cfg *config

// 保存配置文件
func saveConfig(filename string) {
	// 第一次启动， 保存配置文件路径
	cfgfile = filename
}

//
func Start(filename string) {
	// 保存配置文件
	saveConfig(filename)
	if err := Load(false); err != nil {
		// 第一次报错直接退出
		log.Fatal(err)
	}
}

func Load(reload bool) error {

	// 读取配置文件, 配置文件有问题的话，不做后面的处理， 但是会提示错误信息
	if err := readConfig(); err != nil {
		golog.Error(err)
		return err
	}
	// 检测配置文件的name是否重复
	if err := Cfg.check(); err != nil {
		golog.Error(err)
		return err
	}
	// 装载全局配置
	global.Token = Cfg.Token
	global.Listen = Cfg.Listen
	global.IgnoreToken = Cfg.IgnoreToken

	// 初始化报警信息
	Cfg.Alert.InitAlert()
	// 初始化硬件检测
	Cfg.Probe.InitHWAlert()
	// 初始化日志
	golog.InitLogger(Cfg.Log.Path, Cfg.Log.Size, Cfg.Log.Day)
	// 设置所有级别的日志都显示
	golog.Level = golog.All
	golog.Name = "scs.log"
	for index := range Cfg.SC {
		if Cfg.SC[index].Replicate < 1 {
			Cfg.SC[index].Replicate = 1
		}
		if _, ok := SS.Infos[Cfg.SC[index].Name]; !ok {
			SS.Infos[Cfg.SC[index].Name] = make(map[string]*Script)
		}

		if Cfg.SC[index].ContinuityInterval == 0 {
			Cfg.SC[index].ContinuityInterval = time.Minute * 10
		}

		// 第一次启动的时候
		Cfg.fill(index, reload)

	}
	if reload {
		// 删除多余的
		StopUnUseScript()
		b, _ := yaml.Marshal(Cfg)
		// 跟新配置文件
		return ioutil.WriteFile(cfgfile, b, 0644)
	}
	return nil
}

// 读取配置文件
func readConfig() error {
	b, err := ioutil.ReadFile(cfgfile)
	if err != nil {
		return err
	}
	Cfg = &config{}
	err = yaml.Unmarshal(b, Cfg)
	if err != nil {
		return err
	}
	if Cfg.LogCount == 0 {
		global.LogCount = 100
	} else {
		global.LogCount = Cfg.LogCount
	}
	return nil
}

// 运行的时候， 返回状态 Service, 主要验证服务的有效性
func (c *config) Run() {
	c.check()
	SS.Start()
}

// 检测配置脚本是否
func (c *config) check() error {
	// 检查时间
	// 配置信息填充至状态
	checkrepeat := make(map[string]bool)
	for index := range c.SC {
		// 检查名字是否有重复的
		if _, ok := checkrepeat[c.SC[index].Name]; ok {
			return errors.New("配置文件的脚本名重复：" + c.SC[index].Name)
		}
		checkrepeat[c.SC[index].Name] = true

		// 命令行是空的或者name是空的就忽略
		if strings.Trim(c.SC[index].Command, " ") == "" || strings.Trim(c.SC[index].Name, " ") == "" || strings.Trim(c.SC[index].Dir, " ") == "" {
			continue
		}
		if c.SC[index].Replicate < 1 {
			c.SC[index].Replicate = 1
		}
	}
	return nil
}

func (c *config) fill(index int, reload bool) {
	baseEnv := make(map[string]string)
	for k, v := range c.SC[index].Env {
		baseEnv[k] = v
	}
	for i := 0; i < c.SC[index].Replicate; i++ {
		// 根据副本数提取子名称
		subname := fmt.Sprintf("%s_%d", c.SC[index].Name, i)
		if reload {
			// 如果是加载配置文件， 那么删除已经有的
			DelDelScript(subname)
		}

		baseEnv["TOKEN"] = c.Token
		baseEnv["PNAME"] = c.SC[index].Name
		baseEnv["NAME"] = subname
		baseEnv["PORT"] = strconv.Itoa(c.SC[index].Port + i)
		for _, v := range os.Environ() {
			kv := strings.Split(v, "=")
			baseEnv[kv[0]] = kv[1]
		}
		for k, v := range c.SC[index].Env {
			if k == "PATH" {
				if runtime.GOOS == "windows" {
					baseEnv[k] = baseEnv[k] + ";" + v
				} else {
					baseEnv[k] = baseEnv[k] + ":" + v
				}

			} else {
				baseEnv[k] = v
			}
		}
		// 需要单独抽出去<<
		// env := make([]string, 0, len(baseEnv))
		// for k, v := range baseEnv {
		// 	env = append(env, k+"="+v)
		// }

		if _, ok := SS.Infos[c.SC[index].Name][subname]; ok {
			// 修改
			c.update(index, subname, c.SC[index].Command, baseEnv)
			continue
		}
		// 新增
		c.add(index, c.SC[index].Port+i, subname, c.SC[index].Command, baseEnv)
	}
	go func() {
		pname := c.SC[index].Name
		SS.mu.Lock()
		defer SS.mu.Unlock()
		if len(SS.Infos[pname]) > c.SC[index].Replicate {
			for i := len(SS.Infos[pname]) - 1; i >= c.SC[index].Replicate; i-- {
				ne := fmt.Sprintf("%s_%d", pname, i)
				SS.Infos[pname][ne].Stop()
				delete(SS.Infos[pname], ne)
			}
		}
	}()

}

func (c *config) add(index, port int, subname, command string, baseEnv map[string]string) {

	SS.Infos[c.SC[index].Name][subname] = &Script{
		Name:      c.SC[index].Name,
		LookPath:  c.SC[index].LookPath,
		Command:   command,
		Env:       baseEnv,
		Dir:       c.SC[index].Dir,
		Replicate: c.SC[index].Replicate,
		Log:       make(map[string][]string),
		LogLocker: &sync.RWMutex{},
		SubName:   subname,
		Status: &ServiceStatus{
			Name:    subname,
			PName:   c.SC[index].Name,
			Status:  STOP,
			Path:    c.SC[index].Dir,
			Version: c.SC[index].Version,
		},
		DeleteWithExit:     c.SC[index].DeleteWithExit,
		Update:             c.SC[index].Update,
		DisableAlert:       c.SC[index].DisableAlert,
		ContinuityInterval: c.SC[index].ContinuityInterval,
		Always:             c.SC[index].Always,
		Disable:            c.SC[index].Disable,
		AI:                 &alert.AlertInfo{},
		Port:               port,

		AT: c.SC[index].AT,
	}
	SS.Infos[c.SC[index].Name][subname].Log["log"] = make([]string, 0, global.LogCount)
	SS.Infos[c.SC[index].Name][subname].Log["lookPath"] = make([]string, 0, global.LogCount)
	SS.Infos[c.SC[index].Name][subname].Log["update"] = make([]string, 0, global.LogCount)
	if c.SC[index].Cron != nil {
		start, err := time.ParseInLocation("2006-01-02 15:04:05", c.SC[index].Cron.Start, time.Local)
		if err != nil {
			start = time.Time{}
		}
		SS.Infos[c.SC[index].Name][subname].Cron = &Cron{
			Start:   start,
			IsMonth: c.SC[index].Cron.IsMonth,
			Loop:    c.SC[index].Cron.Loop,
		}
	}

	if strings.Trim(c.SC[index].Command, " ") != "" && strings.Trim(c.SC[index].Name, " ") != "" && !c.SC[index].Disable {
		SS.Infos[c.SC[index].Name][subname].Start()
	}

}

func (c *config) update(index int, subname, command string, baseEnv map[string]string) {
	// 修改

	SS.Infos[c.SC[index].Name][subname].Env = baseEnv
	SS.Infos[c.SC[index].Name][subname].LookPath = c.SC[index].LookPath
	if c.SC[index].Cron != nil {
		start, err := time.ParseInLocation("2006-01-02 15:04:05", c.SC[index].Cron.Start, time.Local)
		if err != nil {
			start = time.Time{}
		}
		SS.Infos[c.SC[index].Name][subname].Cron = &Cron{
			Start:   start,
			IsMonth: c.SC[index].Cron.IsMonth,
			Loop:    c.SC[index].Cron.Loop,
		}
	}

	SS.Infos[c.SC[index].Name][subname].Command = command
	SS.Infos[c.SC[index].Name][subname].DeleteWithExit = c.SC[index].DeleteWithExit
	SS.Infos[c.SC[index].Name][subname].Update = c.SC[index].Update
	SS.Infos[c.SC[index].Name][subname].Dir = c.SC[index].Dir
	SS.Infos[c.SC[index].Name][subname].Replicate = c.SC[index].Replicate
	SS.Infos[c.SC[index].Name][subname].Log = make(map[string][]string)
	SS.Infos[c.SC[index].Name][subname].LogLocker = &sync.RWMutex{}
	SS.Infos[c.SC[index].Name][subname].Log["log"] = make([]string, 0, global.LogCount)
	SS.Infos[c.SC[index].Name][subname].Log["lookPath"] = make([]string, 0, global.LogCount)
	SS.Infos[c.SC[index].Name][subname].Log["update"] = make([]string, 0, global.LogCount)
	SS.Infos[c.SC[index].Name][subname].DisableAlert = c.SC[index].DisableAlert
	SS.Infos[c.SC[index].Name][subname].Always = c.SC[index].Always
	SS.Infos[c.SC[index].Name][subname].ContinuityInterval = c.SC[index].ContinuityInterval
	SS.Infos[c.SC[index].Name][subname].Port = c.SC[index].Port + index
	SS.Infos[c.SC[index].Name][subname].AT = c.SC[index].AT
	SS.Infos[c.SC[index].Name][subname].Disable = c.SC[index].Disable
	SS.Infos[c.SC[index].Name][subname].Status.Version = c.SC[index].Version
	// 更新的时候

	if SS.Infos[c.SC[index].Name][subname].Status.Status == STOP {
		// 如果是停止的name就启动
		if strings.Trim(c.SC[index].Command, " ") != "" && strings.Trim(c.SC[index].Name, " ") != "" && !c.SC[index].Disable {
			SS.Infos[c.SC[index].Name][subname].Start()
		}
	}
}

// 更新配置文件
func (c *config) updateConfig(s internal.Script, index int) {
	if s.Dir != "" {
		c.SC[index].Dir = s.Dir
	}
	if s.Command != "" {
		c.SC[index].Command = s.Command
	}
	if s.Env != nil {
		for k, v := range s.Env {
			c.SC[index].Env[k] = v
		}
	}
	if s.Replicate != 0 {
		c.SC[index].Replicate = s.Replicate
	}

	c.SC[index].Always = s.Always
	c.SC[index].DisableAlert = s.DisableAlert
	if s.ContinuityInterval != 0 {
		c.SC[index].ContinuityInterval = s.ContinuityInterval
	} else {
		c.SC[index].ContinuityInterval = time.Minute * 10
	}
	if s.Port != 0 {
		c.SC[index].Port = s.Port
	}
	if s.AT != nil {
		c.SC[index].AT = s.AT
	}
	if s.Version != "" {
		c.SC[index].Version = s.Version
	}
	c.SC[index].Cron = s.Cron

	if len(s.LookPath) > 0 {
		c.SC[index].LookPath = s.LookPath
	}
}

func (c *config) AddScript(s internal.Script) error {
	if _, ok := SS.Infos[s.Name]; !ok {
		SS.Infos[s.Name] = make(map[string]*Script)
	}
	golog.Infof("%+v", s)
	// 添加到配置文件
	for i, v := range c.SC {
		if v.Name == s.Name {
			// 修改
			c.updateConfig(s, i)
			c.fill(i, true)
			b, err := yaml.Marshal(c)
			if err != nil {
				return err
			}
			// 跟新配置文件
			return ioutil.WriteFile(cfgfile, b, 0644)
		}
	}
	// 添加
	// 默认配置
	if s.Replicate < 1 {
		s.Replicate = 1
	}

	if s.ContinuityInterval == 0 {
		s.ContinuityInterval = time.Minute * 10
	}
	c.SC = append(c.SC, s)
	index := len(c.SC) - 1
	c.fill(index, true)

	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	// 跟新配置文件
	return ioutil.WriteFile(cfgfile, b, 0644)
}

func (c *config) DelScript(pname string) error {
	// del := make(chan bool)
	if _, ok := SS.Infos[pname]; ok {
		// go func() {
		// wg := &sync.WaitGroup{}
		for name := range SS.Infos[pname] {
			SS.Infos[pname][name].Remove()
		}

	} else {
		return errors.New("not found this pname:" + pname)
	}
	for i, s := range c.SC {
		if s.Name == pname {
			c.SC = append(c.SC[:i], c.SC[i+1:]...)
			delete(SS.Infos, pname)
			break
		}
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cfgfile, b, 0644)
}