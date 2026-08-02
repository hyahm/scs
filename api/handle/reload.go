package handle

import (
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/internal/cache"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/xmux"
)

func Reload(w http.ResponseWriter, r *http.Request) {
	// 读取新配置
	if err := internal.ReadConfig(); err != nil {
		golog.Error(err)
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
		return
	}
	newCfg := internal.GetConfig()

	// 1. 日志级别热更新
	if newCfg.Debug {
		golog.SetLevel(golog.DEBUG)
	} else {
		golog.SetLevel(golog.INFO)
	}

	// 2. 报警通道热更新
	config.ReloadAlert(newCfg.Alert)

	// 3. 硬件检测热更新
	config.ReloadDetector(&newCfg)

	// 4. 脚本热更新：对比新旧脚本，处理新增/删除/变更
	oldScripts := cache.GetAllScript()
	newScripts := internal.GetAllScript()
	newMap := make(map[string]config.Script)
	for _, s := range newScripts {
		newMap[s.Name] = s
	}

	// 处理新增和变更的脚本
	for _, newScript := range newScripts {
		oldScript, existed := oldScripts[newScript.Name]
		if !existed {
			// 新增脚本
			golog.Infof("reload: 新增脚本 %s", newScript.Name)
			cache.AddScript(newScript)
			continue
		}
		// 已存在，检查是否变更
		if !config.EqualScript(oldScript, newScript) {
			golog.Infof("reload: 脚本 %s 配置变更，重启中", newScript.Name)
			// 停止并删除旧实例
			for _, svc := range cache.GetGroupServer(newScript.Name) {
				svc.Remove()
			}
			cache.RemoveGroupServer(newScript.Name)
			cache.RemoveScript(newScript.Name)
			// 重新添加（会启动新实例）
			cache.AddScript(newScript)
		}
	}

	// 处理已删除的脚本
	for name := range oldScripts {
		if _, stillExists := newMap[name]; !stillExists {
			golog.Infof("reload: 删除脚本 %s", name)
			removeScriptByName(name)
		}
	}

	// 5. 回写配置
	cfg := internal.GetConfig()
	if err := internal.WriteConfig(cfg); err != nil {
		golog.Error(err)
	}

	golog.Info("reload: 所有配置热更新完成")
}

func Fmt(w http.ResponseWriter, r *http.Request) {
	if err := internal.ReadConfig(); err != nil {
		golog.Error(err)
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
		return
	}
	cfg := internal.GetConfig()
	if err := internal.WriteConfig(cfg); err != nil {
		golog.Error(err)
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
		return
	}
}

func removeScriptByName(name string) {
	for _, svc := range cache.GetGroupServer(name) {
		svc.Remove()
	}
	cache.RemoveGroupServer(name)
	cache.RemoveScript(name)
}
