# Api 接口

> 可用于一个服务控制另外一个服务的作用， 请求后，配置文件自动更新配置文件  
> 通过api添加删除挂载的服务(所有接口请带上token请求头)


```
请求头: Token: <token>
pname 是配置文件脚本的name
name  相对于 replicate 来，只有1的个话， pname = name，  否则依次为 pname + "_i" i 是从0开始的索引
# 获取脚本状态
POST /status/{pname}/{name}
POST /status/{pname}
POST /status

# 启动脚本的api 
POST /start/{pname}/{name}
POST /start/{pname}
POST /start
# 停止脚本的api 
POST /stop/{pname}/{name}
POST /stop/{pname}
POST /stop

# 删除脚本的api 
POST /remove/{pname}/{name}
POST /remove/{pname}

# 重启脚本的api
POST /restart/{pname}/{name}
POST /restart/{pname}
POST /restart

# 启用禁用脚本
POST /enable/{pname}
POST /disable/{pname}

# 脚本日志的api
POST /log/{pname}/{name}


POST /script
// 部分参考， 所有配置文件的参数都可以配置
{
  "name": "addscript0",             
	"command" : "pwd",          
	"dir": "/home/cander",
	"version": "v0.0.1"        
}

# 修改可以停止的状态，   
true 不可以停止脚本，   
false 可以停止脚本， 就算收到停止信号也不行,   
 开发常用脚本， 客户端不支持    
POST	/canstop/<name>  
POST	/cannotstop/<name> 

# 外部报警接口
POST  /set/alert  
{
    "title": "haitenda",
    "pname": "asdgasdgasdgfasdf",
    "name": "01358072011",
    "reason": " asdfasdg  \n ashdfljsdf",
    "continuityInterval": 100,  // 这里是间隔多少秒才会发生第二次报警， 默认1小时
    "to": {  // 新增的发件人
        "rocket": ["#general"] 
    }
}
```
> 返回码
```
200:   成功
201：  警告，状态已经被修改
404：  没有找到pname  name
其他错误请参考msg
```
