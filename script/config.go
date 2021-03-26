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

//
func Start(filename string) {
	// 保存配置文件
	cfgfile = filename
	if err := Load(false); err != nil {
		// 第一次报错直接退出
		log.Fatal(err)
	}
	golog.Info("complate")
}

func Load(reload bool) error {
	// reload: 第一次启动     还是 config reload
	// 读取配置文件, 配置文件有问题的话，不做后面的处理， 但是会提示错误信息
	if err := readConfig(); err != nil {
		golog.Error(err)
		return err
	}

	// 初始化日志
	golog.InitLogger(Cfg.Log.Path, Cfg.Log.Size, Cfg.Log.Day)
	// 设置所有级别的日志都显示
	golog.Level = golog.All
	// 设置 日志名， 如果Cfg.Log.Path为空， 那么输出到控制台
	golog.Name = "scs.log"
	// 检测配置文件的name是否重复
	if err := Cfg.check(); err != nil {
		golog.Error(err)
		return err
	}
	// 装载全局配置
	global.Token = Cfg.Token
	global.Listen = Cfg.Listen
	global.IgnoreToken = Cfg.IgnoreToken

	// 初始化报警器信息
	Cfg.Alert.InitAlert()
	// 初始化硬件检测
	Cfg.Probe.InitHWAlert()

	for index := range Cfg.SC {
		// 如果名字为空， 或者 command 为空， 或者 disable=true 那么就跳过
		if Cfg.SC[index].Name == "" || Cfg.SC[index].Command == "" || Cfg.SC[index].Disable {
			continue
		}

		// 如果ss 的key
		if _, ok := SS.Infos[Cfg.SC[index].Name]; !ok {
			SS.Infos[Cfg.SC[index].Name] = make(map[string]*Script)
		}
		// 第一次启动的时候
		Cfg.fill(index, reload)

	}
	if reload {
		// 删除多余的
		RemoveUnUseScript()
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
	}
	return nil
}

func (c *config) fill(index int, reload bool) {

	// 加载环境变量
	baseEnv := make(map[string]string)
	for k, v := range c.SC[index].Env {
		baseEnv[k] = v
	}

	// 填充系统环境变量到
	for _, v := range os.Environ() {
		kv := strings.Split(v, "=")
		baseEnv[kv[0]] = kv[1]
	}

	for k, v := range c.SC[index].Env {
		// path 环境单独处理， 可以多个值， 其他环境变量多个值请以此写完
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

	baseEnv["TOKEN"] = c.Token
	baseEnv["PNAME"] = c.SC[index].Name

	replica := c.SC[index].Replicate
	if replica < 1 {
		replica = 1
	}
	for i := 0; i < replica; i++ {
		// 根据副本数提取子名称
		subname := fmt.Sprintf("%s_%d", c.SC[index].Name, i)
		if reload {
			// 如果是加载配置文件， 那么删除已经有的
			DelDelScript(subname)
		}

		baseEnv["NAME"] = subname
		baseEnv["PORT"] = strconv.Itoa(c.SC[index].Port + i)
		// 需要单独抽出去<<
		// env := make([]string, 0, len(baseEnv))
		// for k, v := range baseEnv {
		// 	env = append(env, k+"="+v)
		// }

		if SS.HasKey(c.SC[index].Name, subname) {
			// 如果存在键值就修改
			golog.Info("update")
			c.update(index, subname, c.SC[index].Command, baseEnv)
		} else {
			golog.Info("add")
			// 新增
			SS.MakeSubStruct(c.SC[index].Name)
			c.add(index, c.SC[index].Port+i, subname, c.SC[index].Command, baseEnv)
		}

	}
	// 删除多余的副本
	go func() {
		pname := c.SC[index].Name

		replicate := c.SC[index].Replicate
		if replicate < 1 {
			replicate = 1
		}
		l := SS.Len()
		if l > 0 && l > replicate {
			for i := l - 1; i >= replicate; i-- {
				subname := fmt.Sprintf("%s_%d", pname, i)
				if reload {
					// 如果是加载配置文件， 那么删除已经有的
					DelDelScript(subname)
				}
				SS.GetScriptFromPnameAndSubname(pname, subname).Remove()
				// SS.Infos[pname][subname].Stop()
				// delete(SS.Infos[pname], subname)
			}
		}
	}()

}

func (c *config) add(index, port int, subname, command string, baseEnv map[string]string) {

	s := &Script{
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
			Version: Command(c.SC[index].Version),
		},
		DeleteWhenExit:     c.SC[index].DeleteWhenExit,
		Update:             c.SC[index].Update,
		DisableAlert:       c.SC[index].DisableAlert,
		ContinuityInterval: Cfg.SC[index].ContinuityInterval,
		Always:             c.SC[index].Always,
		Disable:            c.SC[index].Disable,
		AI:                 &alert.AlertInfo{},
		Port:               port,

		AT: c.SC[index].AT,
	}
	// 生成对应的文件类型
	s.Log["log"] = make([]string, 0, global.LogCount)
	s.Log["lookPath"] = make([]string, 0, global.LogCount)
	s.Log["update"] = make([]string, 0, global.LogCount)
	if c.SC[index].Cron != nil {
		start, err := time.ParseInLocation("2006-01-02 15:04:05", c.SC[index].Cron.Start, time.Local)
		if err != nil {
			start = time.Time{}
		}
		s.Cron = &Cron{
			Start:   start,
			IsMonth: c.SC[index].Cron.IsMonth,
			Loop:    c.SC[index].Cron.Loop,
		}
	}
	SS.AddScript(subname, s)
	s.Start()

}

func (c *config) update(index int, subname, command string, baseEnv map[string]string) {
	// 修改

	scriptInfo := SS.GetScriptFromPnameAndSubname(c.SC[index].Name, subname)

	scriptInfo.Env = baseEnv
	scriptInfo.LookPath = c.SC[index].LookPath
	if c.SC[index].Cron != nil {
		start, err := time.ParseInLocation("2006-01-02 15:04:05", c.SC[index].Cron.Start, time.Local)
		if err != nil {
			start = time.Time{}
		}
		scriptInfo.Cron = &Cron{
			Start:   start,
			IsMonth: c.SC[index].Cron.IsMonth,
			Loop:    c.SC[index].Cron.Loop,
		}
	}

	scriptInfo.Command = command
	scriptInfo.DeleteWhenExit = c.SC[index].DeleteWhenExit
	scriptInfo.Update = c.SC[index].Update
	scriptInfo.Dir = c.SC[index].Dir
	scriptInfo.Replicate = c.SC[index].Replicate
	scriptInfo.Log = make(map[string][]string)
	SS.Infos[c.SC[index].Name][subname].LogLocker = &sync.RWMutex{}
	SS.Infos[c.SC[index].Name][subname].Log["log"] = make([]string, 0, global.LogCount)
	SS.Infos[c.SC[index].Name][subname].Log["lookPath"] = make([]string, 0, global.LogCount)
	SS.Infos[c.SC[index].Name][subname].Log["update"] = make([]string, 0, global.LogCount)
	SS.Infos[c.SC[index].Name][subname].DisableAlert = c.SC[index].DisableAlert
	SS.Infos[c.SC[index].Name][subname].Always = c.SC[index].Always
	SS.Infos[c.SC[index].Name][subname].ContinuityInterval = Cfg.SC[index].ContinuityInterval
	SS.Infos[c.SC[index].Name][subname].Port = c.SC[index].Port + index
	SS.Infos[c.SC[index].Name][subname].AT = c.SC[index].AT
	SS.Infos[c.SC[index].Name][subname].Disable = c.SC[index].Disable
	SS.Infos[c.SC[index].Name][subname].Status.Version = Command(c.SC[index].Version)
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
	SS.Mu.Lock()
	defer SS.Mu.Unlock()
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
