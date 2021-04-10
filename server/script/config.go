package script

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/server/alert"
	"github.com/hyahm/scs/server/logger"
	"github.com/hyahm/scs/server/probe"

	"github.com/hyahm/golog"
	"gopkg.in/yaml.v2"
)

// 读取配置文件， 将值传给script

type Repo struct {
	Url        []string `yaml:"url"`
	Derivative string   `yaml:"derivative"`
}

type config struct {
	Listen      string         `yaml:"listen"`
	Token       string         `yaml:"token"`
	Key         string         `yaml:"key"`
	Pem         string         `yaml:"pem"`
	DisableTls  bool           `yaml:"disableTls"`
	Log         *logger.Logger `yaml:"log"`
	LogCount    int            `yaml:"logCount"`
	IgnoreToken []string       `yaml:"ignoreToken"`
	// Repo        *Repo           `yaml:"repo"`
	Alert *alert.Alert `yaml:"alert"`
	Probe *probe.Probe `yaml:"probe"`
	SC    []*Script    `yaml:"scripts"`
}

// 保存的配置文件路径

// 保存的全局的配置
var cfg *config

//
func Run(filename string) {
	// 保存配置文件
	global.Cfgfile = filename
	// 加载配置文件到各变量中
	if err := Load(false); err != nil {
		// 第一次报错直接退出
		log.Fatal(err)
	}
	golog.Info("complate")
}

func ReadConfig() error {
	return readConfig()
}

func DeleteName(name string) {
	// 检查时间
	// 配置信息填充至状态
	// 读取配置文件
	temp := make([]*Script, 0)
	for _, s := range cfg.SC {
		if s.Name != name {
			temp = append(temp, s)
		}
	}
	cfg.SC = temp
}

func Enable(name string) {
	// 检查时间
	// 配置信息填充至状态
	// 读取配置文件
	for i, s := range cfg.SC {
		if s.Name == name {
			cfg.SC[i].Disable = false
			return
		}
	}

}

func Disable(name string) {
	// 检查时间
	// 配置信息填充至状态
	// 读取配置文件
	for i, s := range cfg.SC {
		if s.Name == name {
			cfg.SC[i].Disable = true
			return
		}
	}
}

func SetReplicate(name string, count int) {
	// 检查时间
	// 配置信息填充至状态
	// 读取配置文件
	for _, s := range cfg.SC {
		if s.Name != name {
			s.Replicate = s.Replicate + count
			return
		}
	}
}

func WriteConfig() error {
	// 检查时间
	// 配置信息填充至状态
	// 读取配置文件
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	// 跟新配置文件
	return ioutil.WriteFile(global.Cfgfile, b, 0644)

}

// config reload
func Reload() error {
	if err := readConfig(); err != nil {
		golog.Error(err)
		return err
	}
	// 运行日志，
	logger.Run(cfg.Log)
	// 检测配置文件的name是否重复， 重复了就返回错误
	if err := cfg.check(); err != nil {
		golog.Error(err)
		return err
	}
	global.Token = cfg.Token
	// 配置文件交给alerter
	alert.Run(cfg.Alert)

	// 初始化硬件检测
	probe.Run(cfg.Probe)

	// cfg.Probe.InitHWAlert()
	RunServer(cfg.SC)
	// if reload {
	// 	// 删除多余的
	// 	RemoveUnUseScript()
	// 	b, _ := yaml.Marshal(cfg)
	// 	// 跟新配置文件
	// 	return ioutil.WriteFile(global.cfgfile, b, 0644)
	// }
	return nil
}

// 第一次启动加载配置文件
func Load(reload bool) error {
	// reload: 第一次启动     还是 config reload
	// 读取配置文件, 配置文件有问题的话，不做后面的处理， 但是会提示错误信息
	if err := readConfig(); err != nil {
		golog.Error(err)
		return err
	}

	// 运行日志，
	logger.Run(cfg.Log)
	// 检测配置文件的name是否重复， 重复了就返回错误
	if err := cfg.check(); err != nil {
		golog.Error(err)
		return err
	}
	// 装载全局配置， 这些配置文件reload 无效
	global.Token = cfg.Token
	global.Listen = cfg.Listen
	global.IgnoreToken = cfg.IgnoreToken
	global.DisableTls = cfg.DisableTls
	global.Key = cfg.Key
	global.Pem = cfg.Pem

	// 配置文件交给alerter
	alert.Run(cfg.Alert)

	// 初始化硬件检测
	probe.Run(cfg.Probe)

	// cfg.Probe.InitHWAlert()
	return RunServer(cfg.SC)

}

// 读取配置文件
func readConfig() error {
	b, err := ioutil.ReadFile(global.Cfgfile)
	if err != nil {
		return err
	}

	cfg = &config{}
	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return err
	}
	if cfg.LogCount == 0 {
		global.LogCount = 100
	} else {
		global.LogCount = cfg.LogCount
	}
	return nil
}

// // 运行的时候， 返回状态 Service, 主要验证服务的有效性
// func (c *config) Run() {
// 	c.check()
// 	ss.Start()
// }

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
	}
	return nil
}

// 更新配置文件
func (c *config) updateConfig(s Script, index int) {
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
