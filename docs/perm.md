# 远程权限管理

如下配置, 我想把test分享给开发查看， 我们可以在对应的脚本下加一个token
```
token: aosdhf-90y890Y*(G0(&TG0*G9F)78fg87)
ignoreToken:
- 127.0.0.1

scripts:
  - name: remote
    command: while true;do date; sleep 5; done

  - name: test
    dir: E:\code\scs
    token: 123
    command: python test.py
    update: "git pull"
```

> 将此token告知对应开发的同事, scsctl 的配置
```
nodes:
  test:
    url: https://192.168.101.11:11111
    token: 123
```

> 可以看到日志
```
[root@master scs]# scsctl log test_0
2022-03-27 21:04:06 - [INFO] - DESKTOP-6SSBCHN - this is test token can stop
2022-03-27 21:04:06 - [INFO] - DESKTOP-6SSBCHN - this is test token can not stop
2022-03-27 21:04:11 - [INFO] - DESKTOP-6SSBCHN - this is test token can stop
2022-03-27 21:04:11 - [INFO] - DESKTOP-6SSBCHN - this is test token can not stop
2022-03-27 21:04:14 - [INFO] - DESKTOP-6SSBCHN - this is test token can stop
```


> 更新操作
```shell
[root@master scs]# scsctl update test
{"code":200,"msg":"waiting update","role":"scripts"}


[root@master scs]# scsctl log test_0
2022-03-27 21:08:17 - [INFO] - DESKTOP-6SSBCHN - this is test token can stop-----
2022-03-27 21:08:17 - [INFO] - DESKTOP-6SSBCHN - this is test token can not stop------
2022-03-27 21:08:20 - [INFO] - DESKTOP-6SSBCHN - this is test token can stop-----
```

> look可以操作的命令如下
```
  env         show env
  get         get info
  help        Help about any command
  info        show service
  log         script log
  restart     restart assign script
  start       start assign script
  status      Print assign script status
  stop        stop script
  update
```

> 以下是源码对应的权限
```
func simpleHandle() *xmux.RouteGroup {
	// 只是调试的权限
	simple := xmux.NewRouteGroup().AddPageKeys(scripts.SimpleRole.ToString())
	simple.Post("/status/{pname}/{name}", handle.Status)
	simple.Post("/status/{pname}", handle.StatusPname)
	simple.Post("/start/{pname}", handle.StartPname)
	simple.Post("/start/{pname}/{name}", handle.Start)
	simple.Post("/update/{pname}/{name}", handle.Update)
	simple.Post("/update/{pname}", handle.UpdatePname)
	simple.Post("/update", handle.UpdateAll) // complete
	simple.Post("/start", handle.StartAll)   // complete
	simple.Post("/status", handle.AllStatus) // complete
	simple.Get("/log/{name}/{int:line}", handle.Log).BindResponse(nil)

	simple.Post("/restart/{pname}/{name}", handle.Restart)
	simple.Post("/restart/{pname}", handle.RestartPname)
	simple.Post("/restart", handle.RestartAll) // complete
	return simple
}

func ScriptHandle() *xmux.RouteGroup {
  // 脚本的操作权限，包含调试的权限， 
	script := xmux.NewRouteGroup().AddPageKeys(scripts.ScriptRole.ToString())
	script.Post("/stop/{pname}/{name}", handle.Stop)
	script.Post("/stop/{pname}", handle.StopPname)
	script.Post("/kill/{pname}", handle.KillPname)
	script.Post("/kill/{pname}/{name}", handle.Kill)
	script.Post("/env/{name}", handle.GetEnvName)
	script.Post("/server/info/{name}", handle.ServerInfo) // 获取某个server信息
	script.Post("/get/servers", handle.GetServers)        // 获取所有server信息 complete
	script.Post("/get/index/{pname}", handle.GetIndex)    // 获取某script 对应副本的index
	script.Post("/cannotstop/{name}", handle.CanNotStop).BindJson(pkg.SignalRequest{})
	script.Post("/parameter/{name}", handle.SetParameter).BindJson(pkg.SignalRequest{})
	script.Post("/canstop/{name}", handle.CanStop)
	script.Post("/get/scripts", handle.GetScripts) // complete
	script.Post("/stop", handle.StopAll)           // complete
	script.Post("/remove/{pname}/{name}", handle.Remove)
	script.Post("/remove/{pname}", handle.RemovePname)
	script.Post("/send/alert", handle.Alert)
	script.AddGroup(simpleHandle())
	return script
}

func AdminHandle() *xmux.RouteGroup {
	// 所有接口的权限
	admin := xmux.NewRouteGroup().AddPageKeys(scripts.AdminRole.ToString()).AddModule(module.CheckToken)

	admin.Post("/get/alert", handle.GetAlert)  
	admin.Post("/-/reload", handle.Reload)      
	admin.Post("/-/fmt", handle.Fmt)           
	admin.Post("/get/alarms", handle.GetAlarms) 
	admin.Post("/get/repo", handle.GetRepo)     
	admin.Post("/script", handle.AddScript).BindJson(&scripts.Script{})
	admin.Post("/enable/{pname}", handle.Enable)  
	admin.Post("/disable/{pname}", handle.Disable)
	admin.AddGroup(ScriptHandle())
	return admin
}
```