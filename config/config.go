package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"scs/alert"
	"scs/global"
	"scs/internal"
	"scs/logger"
	"scs/pkg/script"
	"scs/probe"
	"strconv"
	"strings"
	"time"

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
		if _, ok := script.SS.Infos[Cfg.SC[index].Name]; !ok {
			script.SS.Infos[Cfg.SC[index].Name] = make(map[string]*script.Script)
		}

		if Cfg.SC[index].ContinuityInterval == 0 {
			Cfg.SC[index].ContinuityInterval = time.Minute * 10
		}
		// 重新加载配置文件的时候

		// 第一次启动的时候
		Cfg.fill(index, reload)

	}
	if reload {
		// 删除多余的
		script.StopUnUseScript()
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
	script.SS.Start()
}

// 检测配置脚本是否
func (c *config) check() error {
	// 检查时间
	// 配置信息填充至状态
	checkrepeat := make(map[string]bool)
	for index := range c.SC {
		if c.SC[index].Cron != nil {
			_, err := time.ParseInLocation("2006-01-02 15:04:05", c.SC[index].Cron.Start, time.Local)
			if err != nil {
				return err
			}

		}
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
			script.DelDelScript(subname)
		}

		baseEnv["TOKEN"] = c.Token
		baseEnv["PNAME"] = c.SC[index].Name
		baseEnv["NAME"] = subname
		baseEnv["PORT"] = strconv.Itoa(c.SC[index].Port + i)
		command := strings.ReplaceAll(c.SC[index].Command, "$NAME", subname)
		command = strings.ReplaceAll(command, "$PNAME", c.SC[index].Name)
		command = strings.ReplaceAll(command, "$PORT", strconv.Itoa(c.SC[index].Port+i))

		if _, ok := script.SS.Infos[c.SC[index].Name][subname]; ok {
			// 修改
			c.update(index, subname, command, baseEnv)
			continue
		}
		// 新增
		c.add(index, c.SC[index].Port+i, subname, command, baseEnv)
	}

}

func (c *config) add(index, port int, subname, command string, baseEnv map[string]string) {

	script.SS.Infos[c.SC[index].Name][subname] = &script.Script{
		Name:      c.SC[index].Name,
		LookPath:  c.SC[index].LookPath,
		Command:   command,
		Env:       baseEnv,
		Dir:       c.SC[index].Dir,
		Replicate: c.SC[index].Replicate,
		Log:       make([]string, 0, c.LogCount),
		SubName:   subname,
		Status: &script.ServiceStatus{
			Name:    subname,
			PName:   c.SC[index].Name,
			Status:  script.STOP,
			Path:    c.SC[index].Dir,
			Version: c.SC[index].Version,
		},
		DisableAlert:       c.SC[index].DisableAlert,
		ContinuityInterval: c.SC[index].ContinuityInterval,
		Always:             c.SC[index].Always,
		Disable:            c.SC[index].Disable,
		AI:                 &alert.AlertInfo{},
		Port:               port,

		AT: c.SC[index].AT,
	}
	if c.SC[index].Cron != nil {
		// 前面已经验证过了，不需要再验证
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", c.SC[index].Cron.Start, time.Local)
		script.SS.Infos[c.SC[index].Name][subname].Cron = &script.Cron{
			Start:   start,
			IsMonth: c.SC[index].Cron.IsMonth,
			Loop:    c.SC[index].Cron.Loop,
		}
	}
	// 新增的时候
	if err := script.SS.Infos[c.SC[index].Name][subname].LookCommandPath(); err != nil {
		golog.Error(err)
		return
	}
	if strings.Trim(c.SC[index].Command, " ") != "" && strings.Trim(c.SC[index].Name, " ") != "" &&
		!c.SC[index].Disable {
		script.SS.Infos[c.SC[index].Name][subname].Start()
	}

}

func (c *config) update(index int, subname, command string, baseEnv map[string]string) {
	// 修改
	for k, v := range baseEnv {
		script.SS.Infos[c.SC[index].Name][subname].Env[k] = v
	}

	script.SS.Infos[c.SC[index].Name][subname].LookPath = c.SC[index].LookPath
	if c.SC[index].Cron != nil {
		// 前面已经验证过了，不需要再验证
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", c.SC[index].Cron.Start, time.Local)
		script.SS.Infos[c.SC[index].Name][subname].Cron = &script.Cron{
			Start:   start,
			IsMonth: c.SC[index].Cron.IsMonth,
			Loop:    c.SC[index].Cron.Loop,
		}
	}

	script.SS.Infos[c.SC[index].Name][subname].Command = command
	script.SS.Infos[c.SC[index].Name][subname].Dir = c.SC[index].Dir
	script.SS.Infos[c.SC[index].Name][subname].Replicate = c.SC[index].Replicate
	script.SS.Infos[c.SC[index].Name][subname].Log = make([]string, 0, c.LogCount)
	script.SS.Infos[c.SC[index].Name][subname].DisableAlert = c.SC[index].DisableAlert
	script.SS.Infos[c.SC[index].Name][subname].Always = c.SC[index].Always
	script.SS.Infos[c.SC[index].Name][subname].ContinuityInterval = c.SC[index].ContinuityInterval
	script.SS.Infos[c.SC[index].Name][subname].Port = c.SC[index].Port + index
	script.SS.Infos[c.SC[index].Name][subname].AT = c.SC[index].AT
	script.SS.Infos[c.SC[index].Name][subname].Disable = c.SC[index].Disable
	script.SS.Infos[c.SC[index].Name][subname].Status.Version = c.SC[index].Version
	// 更新的时候
	if err := script.SS.Infos[c.SC[index].Name][subname].LookCommandPath(); err != nil {
		golog.Error(err)
		return
	}
	if script.SS.Infos[c.SC[index].Name][subname].Status.Status == script.STOP {
		// 如果是停止的name就启动
		if strings.Trim(c.SC[index].Command, " ") != "" && strings.Trim(c.SC[index].Name, " ") != "" && !c.SC[index].Disable {
			script.SS.Infos[c.SC[index].Name][subname].Start()
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
	if _, ok := script.SS.Infos[s.Name]; !ok {
		script.SS.Infos[s.Name] = make(map[string]*script.Script)
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
	if _, ok := script.SS.Infos[pname]; ok {
		// go func() {
		// wg := &sync.WaitGroup{}
		for name := range script.SS.Infos[pname] {
			script.SS.Infos[pname][name].Remove()
		}

	} else {
		return errors.New("not found this pname:" + pname)
	}
	for i, s := range c.SC {
		if s.Name == pname {
			c.SC = append(c.SC[:i], c.SC[i+1:]...)
			delete(script.SS.Infos, pname)
			break
		}
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cfgfile, b, 0644)
}
