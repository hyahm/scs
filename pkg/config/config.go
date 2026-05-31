package config

import (
	"errors"
	"os"
	"time"

	"github.com/hyahm/golog"
	"gopkg.in/yaml.v3"
)

var ConfigFile string

var Cfg *Config

func init() {
	Cfg = &Config{}
}

type Config struct {
	Listen      string        `yaml:"listen,omitempty"`
	Token       string        `yaml:"token,omitempty"`
	ProxyHeader string        `yaml:"proxyHeader,omitempty"`
	Key         string        `yaml:"key,omitempty"`
	Cert        string        `yaml:"cert,omitempty"`
	EnableTLS   bool          `yaml:"enableTLS,omitempty"`
	Debug       bool          `yaml:"debug,omitempty"`
	Packet      bool          `yaml:"packet,omitempty"`
	Log         Logger        `yaml:"log,omitempty"`
	IgnoreToken []string      `yaml:"ignoreToken,omitempty"`
	ReadTimeout time.Duration `yaml:"readTimeout,omitempty"`
	// Repo        *Repo          `yaml:"repo,omitempty"`
	Alert   Alert    `yaml:"alert,omitempty"`
	Probe   Probe    `yaml:"probe,omitempty"`
	Scripts []Script `yaml:"scripts,omitempty"`
}

// 保存的配置文件路径

// 读文件

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

// 读取配置文件， 找不到就创建一个空文件
func ReadConfig() error {

	b, err := os.ReadFile(ConfigFile)
	if err != nil {
		f, err := os.Create(ConfigFile)
		if err != nil {
			golog.Error(err)
		}
		f.Close()
		return nil

	}

	err = yaml.Unmarshal(b, Cfg)
	if err != nil {
		golog.Error(err)
		return err
	}
	// 检测配置文件的name是否重复
	err = Cfg.check()
	if err != nil {
		return err
	}
	// 装载全局配置
	// Cfg.Store()
	// global.ProxyHeader = Cfg.ProxyHeader
	// 初始化日志
	// ReloadLogger(Cfg.Log)
	// 初始化报警器信息
	// RunAlert(Cfg.Alert)
	Cfg.Log.initLogger()
	golog.Debug(Cfg.Log.Path)
	// 初始化硬件检测
	Cfg.Probe.initProbe()
	return nil
}

func (c *Config) check() error {
	// 配置信息填充至状态
	checkrepeat := make(map[string]bool)
	for index := range c.Scripts {
		if c.Scripts[index].Name == "" || c.Scripts[index].Command == "" {
			golog.Fatal("name or commond is empty")
		}
		if !CheckScriptNameRule(c.Scripts[index].Name) {
			return errors.New("脚本名不符合命名规则：" + c.Scripts[index].Name)
		}
		// 检查名字是否有重复的
		if _, ok := checkrepeat[c.Scripts[index].Name]; ok {
			return errors.New("配置文件的脚本名重复：" + c.Scripts[index].Name)
		}
		checkrepeat[c.Scripts[index].Name] = true
	}
	return nil
}

// 更新单个script到配置文件
func UpdateScriptToConfigFile(s Script, update bool) error {
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
	for i := range tmp.Scripts {
		if tmp.Scripts[i].Name == s.Name {
			if s.Replicate < 0 {
				tmp.Scripts = append(tmp.Scripts[:i], tmp.Scripts[i+1:]...)
			} else {
				tmp.Scripts[i] = s
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
	tmp.Scripts = nil
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

// 从配置文件删除
func DeleteScriptToConfigFile(s Script, update bool) error {
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
	for i := range tmp.Scripts {
		if tmp.Scripts[i].Name == s.Name {
			tmp.Scripts = append(tmp.Scripts[:i], tmp.Scripts[i+1:]...)
			break
		}
	}
	return tmp.WriteConfig(update)
}

func AddScriptToConfigFile(s *Script) error {
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
	tmp.Scripts = append(tmp.Scripts, *s)
	return tmp.WriteConfig(true)
}
