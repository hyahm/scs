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
	"scs/probe"
	"scs/script"
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
	Listen   string            `yaml:"listen"`
	Token    string            `yaml:"token"`
	Log      logger.Logger     `yaml:"log"`
	LogCount int               `yaml:"logCount"`
	Repo     *Repo             `yaml:"repo"`
	Alert    alert.Alert       `yaml:"alert"`
	Probe    probe.HWAlert     `yaml:"probe"`
	SC       []internal.Script `yaml:"scripts"`
}

// 保存的全局的配置

// 保存的配置文件路径
var cfgfile string
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
	if err := Load(); err != nil {
		// 第一次报错直接退出
		log.Fatal(err)
	}
}

func Load() error {
	// 读取配置文件, 配置文件有问题的话，不做后面的处理， 但是会提示错误信息
	if err := readConfig(); err != nil {
		golog.Error(err)
		return err
	}
	// 检测配置文件
	if err := Cfg.checkName(); err != nil {
		golog.Error(err)
		return err
	}
	// 装载配置
	global.Token = Cfg.Token
	global.Listen = Cfg.Listen
	// 初始化报警信息
	Cfg.initAlert()
	Cfg.initHWAlert()
	golog.InitLogger(Cfg.Log.Path, Cfg.Log.Size, Cfg.Log.Day)
	// 设置所有级别的日志都显示
	golog.Level = golog.All
	for index := range Cfg.SC {
		if Cfg.SC[index].Replicate < 1 {
			Cfg.SC[index].Replicate = 1
		}
		// baseEnv := os.Environ()

		if _, ok := script.SS.Infos[Cfg.SC[index].Name]; !ok {
			script.SS.Infos[Cfg.SC[index].Name] = make(map[string]*script.Script)
		}

		if Cfg.SC[index].ContinuityInterval == 0 {
			Cfg.SC[index].ContinuityInterval = time.Minute * 10
		}
		if Cfg.SC[index].KillTime == 0 {
			Cfg.SC[index].KillTime = time.Second * 1
		}
		Cfg.fill(index)
	}
	return nil
	// 启动多余的服务
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
	c.checkName()
	script.SS.Start()
}

func (c *config) initAlert() {
	// 报警配置转移到了 alert.Alerts
	if alert.Alerts == nil {
		alert.Alerts = make(map[string]alert.SendAlerter)
	}
	if c.Alert.Email != nil {
		if c.Alert.Email.Host != "" && c.Alert.Email.UserName != "" &&
			c.Alert.Email.Password != "" {
			if c.Alert.Email.Port == 0 {
				c.Alert.Email.Port = 465
			}
			alert.Alerts["email"] = c.Alert.Email

		}
	}
	if c.Alert.Rocket != nil {
		if c.Alert.Rocket.Server != "" && c.Alert.Rocket.Username != "" &&
			c.Alert.Rocket.Password != "" {
			alert.Alerts["rocket"] = c.Alert.Rocket
		}

	}

	if c.Alert.Telegram != nil {
		if c.Alert.Telegram.Server != "" && c.Alert.Telegram.Username != "" &&
			c.Alert.Telegram.Password != "" {
			alert.Alerts["telegram"] = c.Alert.Telegram
		}

	}
}

func (c *config) initHWAlert() {

	if c.Probe.Interval == 0 {
		c.Probe.Interval = time.Second * 10
	}
	if c.Probe.ContinuityInterval == 0 {
		c.Probe.ContinuityInterval = time.Hour * 1
	}

	if c.Probe.Cpu == 0 {
		c.Probe.Cpu = 90
	}
	if c.Probe.Mem == 0 {
		c.Probe.Cpu = 90
	}
	if c.Probe.Disk == 0 {
		c.Probe.Cpu = 85
	}
	probe.VarAT.HWA = c.Probe
	if c.Probe.Cpu > 0 || c.Probe.Mem > 0 || c.Probe.Disk > 0 || len(c.Probe.Monitor) > 0 {
		go probe.CheckHardWare()
	}

}

// 检测配置脚本是否
func (c *config) checkName() error {
	// 配置信息填充至状态
	checkrepeat := make(map[string]bool)
	for index := range c.SC {

		// 检查名字是否有重复的
		if _, ok := checkrepeat[c.SC[index].Name]; ok {
			return errors.New("配置文件的脚本名重复：" + c.SC[index].Name)
		}
		checkrepeat[c.SC[index].Name] = true

		// 命令行是空的或者name是空的就忽略
		if strings.Trim(c.SC[index].Command, " ") == "" || strings.Trim(c.SC[index].Name, " ") == "" {
			continue
		}
		if c.SC[index].Replicate < 1 {
			c.SC[index].Replicate = 1
		}
	}
	return nil
}

func (c *config) fill(index int) {
	baseEnv := make(map[string]string)
	for k, v := range c.SC[index].Env {
		baseEnv[k] = v
	}
	for i := 0; i < c.SC[index].Replicate; i++ {
		// 根据副本数提取子名称
		subname := fmt.Sprintf("%s_%d", c.SC[index].Name, i)
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
	// 添加到
	script.SS.Infos[c.SC[index].Name][subname] = &script.Script{
		Name:      c.SC[index].Name,
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
		AI:                 &alert.AlertInfo{},
		Port:               port,
		AT:                 c.SC[index].AT,
		KillTime:           c.SC[index].KillTime,
	}
	// 新增的时候
	script.SetUseScript(subname, c.SC[index].Name)
	script.SS.Infos[c.SC[index].Name][subname].Start()

}

func (c *config) update(index int, subname, command string, baseEnv map[string]string) {
	// 添加到
	for k, v := range baseEnv {
		script.SS.Infos[c.SC[index].Name][subname].Env[k] = v
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
	script.SS.Infos[c.SC[index].Name][subname].Status.Version = c.SC[index].Version
	script.SS.Infos[c.SC[index].Name][subname].KillTime = c.SC[index].KillTime
	golog.Info(c.SC[index].Version)
	if script.SS.Infos[c.SC[index].Name][subname].Status.Status == script.STOP {
		// 如果是停止的name就启动
		script.SS.Infos[c.SC[index].Name][subname].Start()
	}
	// 删除需要删除的服务
	script.DelDelScript(subname)

}

func (c *config) AddScript(s internal.Script) error {
	// 启动脚本
	golog.Info(s.KillTime)
	if s.Replicate < 1 {
		s.Replicate = 1
	}

	if _, ok := script.SS.Infos[s.Name]; !ok {
		script.SS.Infos[s.Name] = make(map[string]*script.Script)
	}
	// 保存上次的副本数，
	if s.ContinuityInterval == 0 {
		s.ContinuityInterval = time.Minute * 10
	}
	if s.KillTime == 0 {
		s.KillTime = time.Second * 1
	}
	// 添加到
	for i, v := range c.SC {
		if v.Name == s.Name {
			c.SC[i] = s
			c.fill(i)

			b, err := yaml.Marshal(c)
			if err != nil {
				return err
			}

			return ioutil.WriteFile(cfgfile, b, 0644)
		}
	}
	golog.Info(s)
	c.SC = append(c.SC, s)
	index := len(c.SC) - 1
	c.fill(index)

	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cfgfile, b, 0644)
}

func (c *config) DelScript(pname string) error {
	if _, ok := script.SS.Infos[pname]; ok {
		for name := range script.SS.Infos[pname] {
			if script.SS.Infos[pname][name].Status.Status == script.RUNNING {
				go script.SS.Infos[pname][name].Stop()
			}
		}

	} else {
		return errors.New("not found this pname:" + pname)
	}
	for i, s := range c.SC {
		if s.Name == pname {
			c.SC = append(c.SC[:i], c.SC[i+1:]...)
		}
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cfgfile, b, 0644)
}
