package internal

import (
	"errors"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/pkg/config"
	"gopkg.in/yaml.v2"
)

var ConfigFile string

// 配置文件的操作

type InConfig struct {
	sync.RWMutex
	config config.Config
}

var cfg = &InConfig{}

// 写入配置文件
func WriteConfig(c config.Config) error {
	cfg.Lock()
	defer cfg.Unlock()

	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFile, b, 0644)
}

func GetConfig() config.Config {
	cfg.RLock()
	defer cfg.RUnlock()
	return cfg.config
}

func GetContinuityInterval() time.Duration {
	cfg.RLock()
	defer cfg.RUnlock()
	return cfg.config.Probe.ContinuityInterval
}

func GetLogPath() string {
	cfg.RLock()
	defer cfg.RUnlock()
	return cfg.config.Log.Path
}

func GetDebug() bool {
	cfg.RLock()
	defer cfg.RUnlock()
	return cfg.config.Debug
}

func GetEnableTLS() bool {
	cfg.RLock()
	defer cfg.RUnlock()
	return cfg.config.EnableTLS
}

func GetKey() string {
	cfg.RLock()
	defer cfg.RUnlock()
	return cfg.config.Key
}

func GetCert() string {
	cfg.RLock()
	defer cfg.RUnlock()
	return cfg.config.Cert
}

func GetToken() string {
	cfg.RLock()
	defer cfg.RUnlock()
	return cfg.config.Token
}

func GetListen() string {
	cfg.RLock()
	defer cfg.RUnlock()
	return cfg.config.Listen
}

func GetAllScript() []config.Script {
	cfg.RLock()
	defer cfg.RUnlock()
	return cfg.config.Scripts
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

	err = yaml.Unmarshal(b, &cfg.config)
	if err != nil {
		golog.Error(err)
		return err
	}
	cfg.Lock()
	defer cfg.Unlock()
	// 检测配置文件的name是否重复
	err = cfg.checkLocked()
	if err != nil {
		return err
	}
	// 初始化警告线
	cfg.config.Probe.InitProbe()
	return nil
}

// 检查脚本名称
func CheckScriptNameRule(name string) bool {
	reg, _ := regexp.Compile(`^\w+$`)
	return reg.MatchString(name)
}

// 注意：调用此函数前，调用方必须持有 c.mu 锁。
func (c *InConfig) checkLocked() error {
	// 外层有锁了，这里不用加锁
	checkrepeat := make(map[string]bool)
	for index := range c.config.Scripts {
		if c.config.Scripts[index].Name == "" || c.config.Scripts[index].Command == "" {
			golog.Fatal("name or commond is empty")
		}
		if !CheckScriptNameRule(c.config.Scripts[index].Name) {
			return errors.New("脚本名不符合命名规则：" + c.config.Scripts[index].Name)
		}
		// 检查名字是否有重复的
		if _, ok := checkrepeat[c.config.Scripts[index].Name]; ok {
			return errors.New("配置文件的脚本名重复：" + c.config.Scripts[index].Name)
		}
		checkrepeat[c.config.Scripts[index].Name] = true
	}
	return nil
}

// 更新单个script到配置文件
func UpdateScriptToConfigFile(s config.Script, update bool) error {
	// 添加
	if !update {
		return nil
	}
	// 默认配置
	// f, err := os.ReadFile(ConfigFile)
	// if err != nil {
	// 	return err
	// }
	cfg.RLock()
	tmp := cfg.config
	cfg.RUnlock()
	tmp.Scripts = make([]config.Script, 0)
	// err = yaml.Unmarshal(f, tmp)
	// if err != nil {
	// 	return err
	// }
	for i := range GetAllScript() {
		if tmp.Scripts[i].Name == s.Name {
			if s.Replicate < 0 {
				tmp.Scripts = append(tmp.Scripts[:i], tmp.Scripts[i+1:]...)
			} else {
				tmp.Scripts[i] = s
			}

		}
	}
	return WriteConfig(tmp)

}

// // 删除配置文件的所有scripts
// func DeleteAllScriptToConfigFile(update bool) error {
// 	// 添加
// 	// 默认配置
// 	f, err := os.ReadFile(ConfigFile)
// 	if err != nil {
// 		return err
// 	}

// 	tmp := &InConfig{}
// 	err = yaml.Unmarshal(f, tmp)
// 	if err != nil {
// 		return err
// 	}
// 	tmp.Scripts = nil
// 	return tmp.WriteConfig(update)
// }

// // 更新script到配置文件
// func RemoveAllScriptToConfigFile(update bool) error {
// 	// 添加
// 	// 默认配置
// 	f, err := os.ReadFile(ConfigFile)
// 	if err != nil {
// 		return err
// 	}

// 	tmp := &Config{}
// 	err = yaml.Unmarshal(f, tmp)
// 	if err != nil {
// 		return err
// 	}

// 	return tmp.WriteConfig(update)
// }

// // 从配置文件删除
// func DeleteScriptToConfigFile(s Script, update bool) error {
// 	if !update {
// 		return nil
// 	}
// 	// 删除默认配置
// 	f, err := os.ReadFile(ConfigFile)
// 	if err != nil {
// 		return err
// 	}

// 	tmp := &InConfig{}
// 	err = yaml.Unmarshal(f, tmp)
// 	if err != nil {
// 		return err
// 	}
// 	for i := range tmp.Scripts {
// 		if tmp.Scripts[i].Name == s.Name {
// 			tmp.Scripts = append(tmp.Scripts[:i], tmp.Scripts[i+1:]...)
// 			break
// 		}
// 	}
// 	return tmp.WriteConfig(update)
// }

// func AddScriptToConfigFile(s *Script) error {
// 	// 默认配置
// 	if !CheckScriptNameRule(s.Name) {
// 		return errors.New("script name must be a word, " + s.Name)
// 	}
// 	f, err := os.ReadFile(ConfigFile)
// 	if err != nil {
// 		return err
// 	}

// 	tmp := &Config{}
// 	err = yaml.Unmarshal(f, tmp)
// 	if err != nil {
// 		return err
// 	}
// 	tmp.Scripts = append(tmp.Scripts, *s)
// 	return tmp.WriteConfig(true)
// }
