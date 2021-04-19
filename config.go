package scs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

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
	Repo        *Repo     `yaml:"repo"`
	Alert       *Alert    `yaml:"alert,omitempty"`
	Probe       *Probe    `yaml:"probe,omitempty"`
	SC          []*Script `yaml:"scripts,omitempty"`
}

// 判断2个map 的值是否相等
func EqualMap(m1, m2 map[string]string) bool {
	if len(m1) == len(m2) {
		if len(m1) == 0 {
			return true
		}

		for k, v := range m1 {
			if m2[k] != v {
				return false
			}
		}
		return true
	} else {
		return false
	}

}

// 判断2个[]string 的值是否相等
func EqualStringArray(s1, s2 []string) bool {

	if len(s1) == len(s2) {
		if len(s1) == 0 {
			return true
		}
		// 先转map
		sm1 := make(map[string]struct{})
		for _, value := range s1 {
			sm1[value] = struct{}{}
		}

		sm2 := make(map[string]struct{})
		for _, value := range s2 {
			sm2[value] = struct{}{}
		}

		for k := range sm1 {
			if _, ok := sm2[k]; !ok {
				return false
			}
		}
		return true
	} else {
		return false
	}

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
	// 配置信息填充至状态
	// 读取配置文件
	for _, s := range Cfg.SC {
		if s.Name != name {
			s.Replicate = s.Replicate + count
			return
		}
	}
}

func DeleteAll() {
	// 检查时间
	// 配置信息填充至状态
	// 读取配置文件
	Cfg.SC = nil
}

func WriteConfig() error {
	// 检查时间
	// 配置信息填充至状态
	// 读取配置文件

	b, _ := yaml.Marshal(Cfg)
	// 跟新配置文件
	return ioutil.WriteFile(cfgfile, b, 0644)

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
	// global.Listen = Cfg.Listen
	// global.IgnoreToken = Cfg.IgnoreToken
	// global.DisableTls = Cfg.DisableTls
	// global.Key = Cfg.Key
	// global.Pem = Cfg.Pem
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
		ok := ReloadScripts(Cfg.SC[index])
		if !ok {
			if ss.Infos[Cfg.SC[index].Name] == nil {
				ss.Infos[Cfg.SC[index].Name] = make(map[string]*Server)
			}
			Cfg.SC[index].MakeServer()
			for _, svc := range ss.Infos[Cfg.SC[index].Name] {
				svc.Start()
			}
			// ss.Scripts[Cfg.SC[index].Name].StartServer()
		}
	}
	// 删除已删除的
	for name := range temp {

		ss.Scripts[name].RemoveScript()
	}
	return nil
}

func ReloadScripts(script *Script) bool {
	// 对对碰， 处理存在的
	if _, ok := ss.Scripts[script.Name]; ok {
		// 对比
		go func() {
			// 需要重启的
			if !CompareScript(script, ss.Scripts[script.Name]) {
				ss.Scripts[script.Name] = script
				err := script.RestartScript()
				if err != nil {
					golog.Error()
				}
			}
			oldReplicate := ss.Scripts[script.Name].Replicate
			if oldReplicate == 0 {
				oldReplicate = 1
			}

			newReplicate := script.Replicate
			if newReplicate == 0 {
				newReplicate = 1
			}
			if oldReplicate == newReplicate {
				return
			}
			if oldReplicate > newReplicate {
				// 如果大于的话， 那么就删除多余的
				for i := newReplicate; i < oldReplicate; i++ {
					golog.Info("remove " + script.Name + fmt.Sprintf("_%d", i))
					ss.Infos[script.Name][script.Name+fmt.Sprintf("_%d", i)].Remove()
				}
			} else {
				script.MakeEnv()
				start := time.Now()
				portIndex := 0
				for i := oldReplicate; i < newReplicate; i++ {
					// 根据副本数提取子名称

					subname := fmt.Sprintf("%s_%d", script.Name, i)
					if script.Port > 0 {
						portIndex += probePort(script.Port)
						script.TempEnv["PORT"] = strconv.Itoa(script.Port + i + portIndex)
						ss.Infos[script.Name][subname] = script.add(script.Port+i+portIndex, i, subname)
					} else {
						script.TempEnv["PORT"] = "0"
						ss.Infos[script.Name][subname] = script.add(0, i, subname)
					}
					script.TempEnv["NAME"] = subname

					ss.Infos[script.Name][subname].Start()

				}
				golog.Info(time.Since(start).Seconds())
			}
			ss.Scripts[script.Name] = script
		}()
		return true
	} else {
		return false
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
		// 将数据填充至 SS
		ss.Scripts[Cfg.SC[index].Name] = Cfg.SC[index]
		// 启动服务
		if ss.Infos[Cfg.SC[index].Name] == nil {
			ss.Infos[Cfg.SC[index].Name] = make(map[string]*Server)
		}
		ss.Scripts[Cfg.SC[index].Name].MakeServer()
		for _, svc := range ss.Infos[Cfg.SC[index].Name] {
			svc.Start()
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

// func (c *config) fill(index int, reload bool) {

// 	// 加载环境变量
// 	baseEnv := make(map[string]string)

// 	// 填充系统环境变量到
// 	pathEnvName := "PATH"
// 	for _, v := range os.Environ() {
// 		kv := strings.Split(v, "=")
// 		if strings.ToUpper(kv[0]) == pathEnvName {
// 			pathEnvName = kv[0]
// 		}
// 		baseEnv[kv[0]] = kv[1]
// 	}
// 	golog.Info(pathEnvName)
// 	for k, v := range c.SC[index].Env {
// 		// path 环境单独处理， 可以多个值， 其他环境变量多个值请以此写完
// 		if strings.ToLower(k) == strings.ToLower(pathEnvName) {
// 			if runtime.GOOS == "windows" {
// 				baseEnv[pathEnvName] = baseEnv[pathEnvName] + ";" + v
// 			} else {
// 				golog.Info(pathEnvName)
// 				baseEnv[pathEnvName] = baseEnv[pathEnvName] + ":" + v
// 			}
// 		} else {
// 			baseEnv[k] = v
// 		}
// 	}

// 	baseEnv["TOKEN"] = c.Token
// 	baseEnv["PNAME"] = c.SC[index].Name

// 	replica := c.SC[index].Replicate
// 	if replica < 1 {
// 		replica = 1
// 	}
// 	for i := 0; i < replica; i++ {
// 		// 根据副本数提取子名称
// 		subname := fmt.Sprintf("%s_%d", c.SC[index].Name, i)
// 		if reload {
// 			// 如果是加载配置文件， 那么删除已经有的
// 			golog.Info("delete subname")
// 			DelDelScript(subname)
// 		}
// 		baseEnv["NAME"] = subname
// 		baseEnv["PORT"] = strconv.Itoa(c.SC[index].Port + i)
// 		// 需要单独抽出去<<
// 		// env := make([]string, 0, len(baseEnv))
// 		// for k, v := range baseEnv {
// 		// 	env = append(env, k+"="+v)
// 		// }

// 		if SS.HasSubName(c.SC[index].Name, subname) {
// 			// 如果存在键值就修改

// 			golog.Info("update")
// 			c.update(index, subname, c.SC[index].Command, baseEnv)
// 		} else {
// 			golog.Info("add")
// 			// 新增
// 			SS.MakeSubStruct(c.SC[index].Name)
// 			c.add(index, c.SC[index].Port+i, subname, c.SC[index].Command, baseEnv)
// 		}

// 	}
// 	// 删除多余的副本
// 	go func() {
// 		pname := c.SC[index].Name

// 		replicate := c.SC[index].Replicate
// 		if replicate < 1 {
// 			replicate = 1
// 		}
// 		l := SS.Len()
// 		if l > 0 && l > replicate {
// 			for i := l - 1; i >= replicate; i-- {
// 				subname := fmt.Sprintf("%s_%d", pname, i)
// 				if reload {
// 					// 如果是加载配置文件， 那么删除已经有的
// 					DelDelScript(subname)
// 				}
// 				SS.GetScriptFromPnameAndSubname(pname, subname).Remove()
// 				// SS.Infos[pname][subname].Stop()
// 				// delete(SS.Infos[pname], subname)
// 			}
// 		}
// 	}()

// }

// func (c *config) update(index int, subname, command string, baseEnv map[string]string) {
// 	// 修改

// 	scriptInfo := ss.GetScriptFromPnameAndSubname(c.SC[index].Name, subname)

// 	scriptInfo.Env = baseEnv
// 	scriptInfo.LookPath = c.SC[index].LookPath
// 	if c.SC[index].Cron != nil {
// 		start, err := time.ParseInLocation("2006-01-02 15:04:05", c.SC[index].Cron.Start, time.Local)
// 		if err != nil {
// 			start = time.Time{}
// 		}
// 		scriptInfo.Cron = &Cron{
// 			StartTime: start,
// 			IsMonth:   c.SC[index].Cron.IsMonth,
// 			Loop:      c.SC[index].Cron.Loop,
// 		}
// 	}

// 	scriptInfo.Command = command
// 	scriptInfo.Update = c.SC[index].Update
// 	scriptInfo.Log = make(map[string][]string)
// 	ss.Infos[c.SC[index].Name][subname].LogLocker = &sync.RWMutex{}
// 	ss.Infos[c.SC[index].Name][subname].Log["log"] = make([]string, 0, global.LogCount)
// 	ss.Infos[c.SC[index].Name][subname].Log["lookPath"] = make([]string, 0, global.LogCount)
// 	ss.Infos[c.SC[index].Name][subname].Log["update"] = make([]string, 0, global.LogCount)
// 	ss.Infos[c.SC[index].Name][subname].Script = c.SC[index]
// 	ss.Infos[c.SC[index].Name][subname].ContinuityInterval = Cfg.SC[index].ContinuityInterval
// 	ss.Infos[c.SC[index].Name][subname].Port = c.SC[index].Port + index
// 	ss.Infos[c.SC[index].Name][subname].AT = c.SC[index].AT
// 	ss.Infos[c.SC[index].Name][subname].Disable = c.SC[index].Disable
// 	ss.Infos[c.SC[index].Name][subname].Status.Version = getVersion(c.SC[index].Version)
// 	ss.Infos[c.SC[index].Name][subname].Status.Disable = c.SC[index].Disable
// 	// 更新的时候

// 	if ss.Infos[c.SC[index].Name][subname].Status.Status == STOP {
// 		// 如果是停止的name就启动
// 		if strings.Trim(c.SC[index].Command, " ") != "" && strings.Trim(c.SC[index].Name, " ") != "" && !c.SC[index].Disable {
// 			ss.Infos[c.SC[index].Name][subname].Start()
// 		}
// 	}
// }

// // 添加script到配置文件
// func (c *config) updateConfig(s Script, index int) {
// 	if s.Dir != "" {
// 		c.SC[index].Dir = s.Dir
// 	}
// 	if s.Command != "" {
// 		c.SC[index].Command = s.Command
// 	}
// 	if s.Env != nil {
// 		for k, v := range s.Env {
// 			c.SC[index].Env[k] = v
// 		}
// 	}
// 	if s.Replicate != 0 {
// 		c.SC[index].Replicate = s.Replicate
// 	}

// 	c.SC[index].Always = s.Always
// 	c.SC[index].DisableAlert = s.DisableAlert
// 	if s.Port != 0 {
// 		c.SC[index].Port = s.Port
// 	}
// 	if s.AT != nil {
// 		c.SC[index].AT = s.AT
// 	}
// 	if s.Version != "" {
// 		c.SC[index].Version = s.Version
// 	}
// 	c.SC[index].Cron = s.Cron

// 	if len(s.LookPath) > 0 {
// 		c.SC[index].LookPath = s.LookPath
// 	}
// }

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
			golog.Info(1111)
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
	ss.ScriptLocker.RLock()
	defer ss.ScriptLocker.RUnlock()
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
	ss.ServerLocker.Lock()
	defer ss.ServerLocker.Unlock()
	if _, ok := ss.Infos[pname]; ok {
		// go func() {
		// wg := &sync.WaitGroup{}
		for name := range ss.Infos[pname] {
			ss.Infos[pname][name].Remove()
		}

	} else {
		return errors.New("not found this pname:" + pname)
	}
	for i, s := range c.SC {
		if s.Name == pname {
			c.SC = append(c.SC[:i], c.SC[i+1:]...)
			delete(ss.Infos, pname)
			break
		}
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cfgfile, b, 0644)
}
