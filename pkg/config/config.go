package config

import (
	"errors"
	"os"
	"time"

	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/pkg/config/alert"
	"github.com/hyahm/scs/pkg/config/logger"
	"github.com/hyahm/scs/pkg/config/probe"
	"github.com/hyahm/scs/pkg/config/scripts"

	"github.com/hyahm/golog"
	"gopkg.in/yaml.v2"
)

var ConfigFile string

type Config struct {
	Listen      string         `yaml:"listen,omitempty"`
	Token       string         `yaml:"token,omitempty"`
	ProxyHeader string         `yaml:"proxyHeader,omitempty"`
	Key         string         `yaml:"key,omitempty"`
	Cert        string         `yaml:"cert,omitempty"`
	EnableTLS   bool           `yaml:"enableTLS,omitempty"`
	Packet      bool           `yaml:"packet,omitempty"`
	Log         *logger.Logger `yaml:"log,omitempty"`
	IgnoreToken []string       `yaml:"ignoreToken,omitempty"`
	ReadTimeout time.Duration  `yaml:"readTimeout,omitempty"`
	// Repo        *Repo          `yaml:"repo,omitempty"`
	Alert *alert.Alert      `yaml:"alert,omitempty"`
	Probe *probe.Probe      `yaml:"probe,omitempty"`
	SC    []*scripts.Script `yaml:"scripts,omitempty"`
}

// 保存的配置文件路径

// 读文件
func ReadConfig() (*Config, error) {
	// 依次启动
	return reLoad()
}

func (c *Config) Store() {
	global.CS.Cert = c.Cert
	global.CS.EnableTLS = c.EnableTLS
	global.CS.IgNoreToken = c.IgnoreToken
	global.CS.Key = c.Key
	if c.Listen == "" {
		c.Listen = ":11111"
	}
	global.CS.Listen = c.Listen
	global.CS.Token = c.Token
	// global.CS.LogCount = c.log
	global.CS.ReadTime = c.ReadTimeout
}

// 写入配置文件
func (c *Config) WriteConfig(update bool) error {
	if !update {
		return nil
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFile, b, 0644)
}

func reLoad() (*Config, error) {
	// reload: 第一次启动     还是 config reload
	// 读取配置文件, 配置文件有问题的话，不做后面的处理， 但是会提示错误信息
	return readConfig()
}

// 读取配置文件， 找不到就创建一个空文件
func readConfig() (*Config, error) {
	cfg := &Config{}
	b, err := os.ReadFile(ConfigFile)
	if err != nil {
		f, err := os.Create(ConfigFile)
		if err != nil {
			golog.Error(err)
		}
		f.Close()
		return cfg, nil

	}

	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		golog.Error(err)
		return nil, err
	}
	// 检测配置文件的name是否重复
	err = cfg.check()
	if err != nil {
		return nil, err
	}
	// 装载全局配置
	cfg.Store()
	global.ProxyHeader = cfg.ProxyHeader
	// 初始化日志
	logger.ReloadLogger(cfg.Log)
	// 初始化报警器信息
	alert.RunAlert(cfg.Alert)
	// 初始化硬件检测
	probe.RunProbe(cfg.Probe)

	return cfg, nil
}

func (c *Config) check() error {
	// 配置信息填充至状态
	checkrepeat := make(map[string]bool)
	for index := range c.SC {
		if c.SC[index].Name == "" || c.SC[index].Command == "" {
			golog.Fatal("name or commond is empty")
		}
		if !CheckScriptNameRule(c.SC[index].Name) {
			return errors.New("脚本名不符合命名规则：" + c.SC[index].Name)
		}
		// 检查名字是否有重复的
		if _, ok := checkrepeat[c.SC[index].Name]; ok {
			return errors.New("配置文件的脚本名重复：" + c.SC[index].Name)
		}
		checkrepeat[c.SC[index].Name] = true
	}
	return nil
}

// 更新单个script到配置文件
func UpdateScriptToConfigFile(s *scripts.Script, update bool) error {
	// 添加
	if !update {
		return nil
	}
	// 默认配置
	f, err := os.ReadFile(ConfigFile)
	if err != nil {
		return err
	}

	tmp := &Config{}
	err = yaml.Unmarshal(f, tmp)
	if err != nil {
		return err
	}
	for i := range tmp.SC {
		if tmp.SC[i].Name == s.Name {
			if s.Replicate < 0 {
				tmp.SC = append(tmp.SC[:i], tmp.SC[i+1:]...)
			} else {
				tmp.SC[i] = s
			}

		}
	}
	return tmp.WriteConfig(true)

}

// 删除配置文件的所有scripts
func DeleteAllScriptToConfigFile(update bool) error {
	// 添加
	// 默认配置
	f, err := os.ReadFile(ConfigFile)
	if err != nil {
		return err
	}

	tmp := &Config{}
	err = yaml.Unmarshal(f, tmp)
	if err != nil {
		return err
	}
	tmp.SC = nil
	return tmp.WriteConfig(update)
}

// 更新script到配置文件
func RemoveAllScriptToConfigFile(update bool) error {
	// 添加
	// 默认配置
	f, err := os.ReadFile(ConfigFile)
	if err != nil {
		return err
	}

	tmp := &Config{}
	err = yaml.Unmarshal(f, tmp)
	if err != nil {
		return err
	}

	return tmp.WriteConfig(update)
}

// func RemoveAllScripts() {
// 	// 删除所有脚本
// 	RemoveAllScriptToConfigFile()
// }

// 从配置文件删除
func DeleteScriptToConfigFile(s *scripts.Script, update bool) error {
	if !update {
		return nil
	}
	// 删除默认配置
	f, err := os.ReadFile(ConfigFile)
	if err != nil {
		return err
	}

	tmp := &Config{}
	err = yaml.Unmarshal(f, tmp)
	if err != nil {
		return err
	}
	for i := range tmp.SC {
		if tmp.SC[i].Name == s.Name {
			tmp.SC = append(tmp.SC[:i], tmp.SC[i+1:]...)
			break
		}
	}
	return tmp.WriteConfig(update)
}

func AddScriptToConfigFile(s *scripts.Script) error {
	// 默认配置
	if !CheckScriptNameRule(s.Name) {
		return errors.New("script name must be a word, " + s.Name)
	}
	f, err := os.ReadFile(ConfigFile)
	if err != nil {
		return err
	}

	tmp := &Config{}
	err = yaml.Unmarshal(f, tmp)
	if err != nil {
		return err
	}
	tmp.SC = append(tmp.SC, s)
	return tmp.WriteConfig(true)
}
