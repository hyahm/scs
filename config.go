package scs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/hyahm/scs/global"

	"github.com/hyahm/golog"
	"gopkg.in/yaml.v2"
)

type Repo struct {
	Url        []string `yaml:"url"`
	Derivative string   `yaml:"derivative"`
}

type config struct {
	Listen      string    `yaml:"listen"`
	Token       string    `yaml:"token,omitempty"`
	Key         string    `yaml:"key,omitempty"`
	Pem         string    `yaml:"pem,omitempty"`
	DisableTls  bool      `yaml:"disableTls"`
	Log         *Logger   `yaml:"log,omitempty"`
	LogCount    int       `yaml:"logCount,omitempty"`
	IgnoreToken []string  `yaml:"ignoreToken,omitempty"`
	Repo        *Repo     `yaml:"repo,omitempty"`
	Alert       *Alert    `yaml:"alert,omitempty"`
	Probe       *Probe    `yaml:"probe,omitempty"`
	SC          []*Script `yaml:"scripts,omitempty"`
}

// 保存的配置文件路径
var cfgfile string

// 保存的全局的配置
var Cfg *config

//
func Start(filename string) {
	// 保存配置文件
	cfgfile = filename
	if err := Load(); err != nil {
		// 第一次报错直接退出
		golog.Fatal(err)
	}
}

func ReadConfig() error {
	return readConfig()
}

func DeleteName(name string) {
	// 检查时间
	// 配置信息填充至状态
	// 读取配置文件
	temp := make([]*Script, 0)
	for _, s := range Cfg.SC {
		if s.Name != name {
			temp = append(temp, s)
		}
	}
	Cfg.SC = temp
}

func Enable(name string) {
	// 检查时间
	// 配置信息填充至状态
	// 读取配置文件
	for i, s := range Cfg.SC {
		if s.Name == name {
			Cfg.SC[i].Disable = false
			return
		}
	}

}

func Disable(name string) {
	// 检查时间
	// 配置信息填充至状态
	// 读取配置文件
	for i, s := range Cfg.SC {
		if s.Name == name {
			Cfg.SC[i].Disable = true
			return
		}
	}
}

func SetReplicate(name string, count int) {
	// 检查时间
	for _, s := range Cfg.SC {
		if s.Name != name {
			s.Replicate = s.Replicate + count
			return
		}
	}
}

func ReLoad() error {
	// reload: 第一次启动     还是 config reload
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
	// 初始化日志
	ReloadLogger(Cfg.Log)
	// 初始化报警器信息
	RunAlert(Cfg.Alert)
	// 初始化硬件检测
	RunProbe(Cfg.Probe)

	// 这里也要判断修改
	// 拷贝一份当前存在的所有的脚本
	temp := make(map[string]struct{})
	for name := range ss.Scripts {
		temp[name] = struct{}{}
	}
	for index := range Cfg.SC {
		// 将数据填充至 SS, 返回是否存在此脚本
		delete(temp, Cfg.SC[index].Name)
		ReloadScripts(Cfg.SC[index])
	}
	// 删除已删除的
	for name := range temp {
		ss.Scripts[name].RemoveScript()
	}
	return nil
}

func ReloadScripts(script *Script) {
	// 对对碰， 处理存在的
	if _, ok := ss.Scripts[script.Name]; ok {
		// 对比
		// 需要重启的
		oldReplicate := ss.Scripts[script.Name].Replicate
		if oldReplicate == 0 {
			oldReplicate = 1
		}
		//
		if script.Replicate == 1 {
			script.Replicate = 0
		}
		newReplicate := script.Replicate

		if newReplicate == 0 {
			newReplicate = 1
		}
		if !CompareScript(script, ss.Scripts[script.Name]) {
			// 如果不一样， 那么 就需要重新启动服务
			golog.Info("restart server")
			ss.Scripts[script.Name] = script
			// 先停止脚本， 更新 server
			err := script.StopScript()
			if err != nil {
				golog.Error()
			}
			script.MakeServer()
			replicate := script.Replicate
			if replicate == 0 {
				replicate = 1
			}
			for i := 0; i < replicate; i++ {
				subname := NewSubname(script.Name, i)
				ss.Infos[subname].Start()
			}

			// 之前有多的副本就需要删除了
			for i := newReplicate; i < oldReplicate; i++ {
				golog.Info("remove " + script.Name + fmt.Sprintf("_%d", i))
				ss.Infos[NewSubname(script.Name, i)].Remove()
			}
			return
		}

		if oldReplicate == newReplicate {
			// 如果一样的名字， 副本数一样的就直接跳过
			return
		}
		if oldReplicate > newReplicate {
			// 如果大于的话， 那么就删除多余的
			for i := newReplicate; i < oldReplicate; i++ {
				golog.Info("remove " + script.Name + fmt.Sprintf("_%d", i))
				ss.Infos[NewSubname(script.Name, i)].Remove()
			}
		} else {
			script.MakeEnv()

			portIndex := 0
			for i := oldReplicate; i < newReplicate; i++ {
				// 根据副本数提取子名称

				subname := NewSubname(script.Name, i)
				if script.Port > 0 {
					portIndex += probePort(script.Port)
					script.TempEnv["PORT"] = strconv.Itoa(script.Port + i + portIndex)
					ss.Infos[subname] = script.add(script.Port+i+portIndex, subname)
				} else {
					script.TempEnv["PORT"] = "0"
					ss.Infos[subname] = script.add(0, subname)
				}
				script.TempEnv["NAME"] = subname.String()

				ss.Infos[subname].Start()

			}
		}
		ss.Scripts[script.Name].Replicate = script.Replicate

	} else {
		golog.Info(script.Name)

		ss.Scripts[script.Name] = script
		script.MakeServer()
		replicate := script.Replicate
		if replicate == 0 {
			replicate = 1
		}
		for i := 0; i < replicate; i++ {
			subname := NewSubname(script.Name, i)
			ss.Infos[subname].Start()
		}
		golog.Info("load complete: ", script.Name)
	}
}

func Load() error {
	// reload: 第一次启动     还是 config reload
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
	global.DisableTls = Cfg.DisableTls
	global.Key = Cfg.Key
	global.Pem = Cfg.Pem
	// 初始化日志

	ReloadLogger(Cfg.Log)
	// 初始化报警器信息
	RunAlert(Cfg.Alert)
	// 初始化硬件检测
	RunProbe(Cfg.Probe)
	for index := range Cfg.SC {
		// 将数据填充至 SS的script中
		ss.Scripts[Cfg.SC[index].Name] = Cfg.SC[index]
		ss.Scripts[Cfg.SC[index].Name].MakeServer()
		replicate := Cfg.SC[index].Replicate
		if replicate == 0 {
			replicate = 1
		}
		for i := 0; i < replicate; i++ {
			subname := NewSubname(Cfg.SC[index].Name, i)
			ss.Infos[subname].Start()
		}

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
		if Cfg.SC[index].Name == "" || Cfg.SC[index].Command == "" {
			golog.Fatal("name or commond is empty")
		}
		// 检查名字是否有重复的
		if _, ok := checkrepeat[c.SC[index].Name]; ok {
			return errors.New("配置文件的脚本名重复：" + c.SC[index].Name)
		}
		checkrepeat[c.SC[index].Name] = true
	}
	return nil
}

// 更新script到配置文件
func UpdateScriptToConfigFile(s *Script) error {
	// 添加
	// 默认配置
	f, err := ioutil.ReadFile(cfgfile)
	if err != nil {
		return err
	}

	tmp := &config{}
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
	b, err := yaml.Marshal(tmp)
	if err != nil {
		return err
	}
	// 跟新配置文件
	return ioutil.WriteFile(cfgfile, b, 0644)
}

// 删除配置文件的所有scripts
func DeleteAllScriptToConfigFile() error {
	// 添加
	// 默认配置
	f, err := ioutil.ReadFile(cfgfile)
	if err != nil {
		return err
	}

	tmp := &config{}
	err = yaml.Unmarshal(f, tmp)
	if err != nil {
		return err
	}
	tmp.SC = nil
	b, err := yaml.Marshal(tmp)
	if err != nil {
		return err
	}
	// 跟新配置文件
	return ioutil.WriteFile(cfgfile, b, 0644)
}

// 更新script到配置文件
func RemoveAllScriptToConfigFile() error {
	// 添加
	// 默认配置
	f, err := ioutil.ReadFile(cfgfile)
	if err != nil {
		return err
	}

	tmp := &config{}
	err = yaml.Unmarshal(f, tmp)
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(tmp)
	if err != nil {
		return err
	}
	// 跟新配置文件
	return ioutil.WriteFile(cfgfile, b, 0644)
}

func RemoveAllScripts() {
	// 删除所有脚本
	RemoveAllScriptToConfigFile()
	for _, v := range ss.Scripts {
		v.RemoveScript()
	}
	ss.Scripts = make(map[string]*Script)
}

func DeleteScriptToConfigFile(s *Script) error {
	// 删除默认配置
	f, err := ioutil.ReadFile(cfgfile)
	if err != nil {
		return err
	}

	tmp := &config{}
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
	b, err := yaml.Marshal(tmp)
	if err != nil {
		return err
	}
	// 跟新配置文件
	return ioutil.WriteFile(cfgfile, b, 0644)
}

func HaveScript(pname string) bool {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	_, ok := ss.Scripts[pname]
	return ok
}

func AddScriptToConfigFile(s *Script) error {
	// 添加
	// 默认配置
	f, err := ioutil.ReadFile(cfgfile)
	if err != nil {
		return err
	}

	tmp := &config{}
	err = yaml.Unmarshal(f, tmp)
	if err != nil {
		return err
	}
	tmp.SC = append(tmp.SC, s)
	b, err := yaml.Marshal(tmp)
	if err != nil {
		return err
	}
	// 跟新配置文件
	return ioutil.WriteFile(cfgfile, b, 0644)
}

func (c *config) DelScript(pname string) error {
	// del := make(chan bool)
	ss.Mu.Lock()
	defer ss.Mu.Unlock()

	if _, ok := ss.Scripts[pname]; ok {
		// go func() {
		// wg := &sync.WaitGroup{}
		replicate := ss.Scripts[pname].Replicate
		if replicate == 0 {
			replicate = 1
		}

		for i := 0; i < replicate; i++ {
			subname := NewSubname(pname, i)
			ss.Infos[subname].Remove()
		}

	} else {
		return errors.New("not found this pname:" + pname)
	}
	for i, s := range c.SC {
		if s.Name == pname {
			c.SC = append(c.SC[:i], c.SC[i+1:]...)
			delete(ss.Scripts, pname)
			break
		}
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cfgfile, b, 0644)
}
