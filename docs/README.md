# scs

service control service or script

# 主要功能
1.  为了保护数据不会在执行中被人工手动中断丢失， 可以让服务在某段时间内才能停止
2.  监控硬件信息， 主要是磁盘， cpu， 内存,  scs服务  
3.  服务之间可以相互控制增删改查  
4.  报警功能api  
5.  支持定时器功能执行命令或脚本
6.  客户端控制多台服务器
7.  通过配置文件安装服务
8.  可以将一些执行耗时的脚本托管给scs处理， 以便快速返回结果  


# 服务控制脚本
类似supervisor,但是更高级，支持所有系统 
自带监控及通知   
服务控制脚本能否停止 最大程度防止脚本数据丢失  
码云地址: https://gitee.com/cander/scs
# 功能
- [x] 跨平台  
- [x] 支持启动所有脚本  
- [x] 自带硬件检测  
- [x] 一键可启动(linux)  
- [x] 支持的报警渠道- 邮箱, rocket.chat报警, telegram, 企业微信

[版本更新](update.md)



# Api 接口
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
POST /remove

# 重启脚本的api
POST /restart/{pname}/{name}
POST /restart/{pname}
POST /restart

# 启用禁用脚本
POST /restart/{pname}
POST /restart/{pname}

# 脚本日志的api
POST /log/{pname}/{name}

POST /delete/{pname}
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

# 开发包  只有go和python 版， 其他语言请参考上面的api自行封装  
### go版本， 本身就自带

`go get github.com/hyahm/scs/client`
```
package main

import (
	"fmt"
	"log"

	"github.com/hyahm/scs/client"
)

func main() {
	cli := scs.NewClient()
	// 获取https://127.0.0.1:11111 的 脚本状态
	b, err := cli.Status()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

```
> 输出
```vim
{
        "data": [
                {
                        "name": "test_0",
                        "ppid": 0,
                        "status": "Stop",
                        "command": "python test.py",
                        "pname": "test",
                        "path": "F:\\scs",
                        "cannotStop": false,
                        "start": 0,
                        "version": "",
                        "Always": false,
                        "restartCount": 0
                }
        ],
        "code": 200
}
```

### python 版本

https://pypi.org/project/pyscs/
```
pip install pyscs
```
# 客户端
1.1.0 版开始， 配置文件必须要配置值， 不然什么也不会出来
```yaml
nodes:
  localhost: 
    url: https://127.0.0.1:11111
    token:
group:
  local:
    - localhost
```

```

scsctl status 
scsctl status pname
scsctl status pname name
scsctl start 
scsctl start pname
scsctl start pname name
scsctl restart --all
scsctl restart pname 
scsctl restart pname name
scsctl kill --all
scsctl kill pname 
scsctl kill pname name
scsctl stop --all
scsctl stop pname 
scsctl stop pname name
scsctl update --all
scsctl update pname 
scsctl update pname name
scsctl remove --all
scsctl remove pname 
scsctl remove pname name
scsctl enable pname
scsctl disable pname
scsctl log  name[:update|log|lookPath] # 不区分大小写
# 加载配置文件
scsctl config reload
```
# 编译二进制文件（go>=1.13）
```
 git clone https://github.com/hyahm/scs.git
 cd scs
 go env -w GOPROXY=https://goproxy.cn,direct # 国外机器不需要这个
 go build -o scsd cmd/scs/main.go
 go build -o /usr/local/bin/scsctl cmd/scsctl/main.go
 cp default.yaml scs.yaml
 ./scsd
 ```

linux(需要 git tar命令， 关闭selinux),mac, windows 请按照上面自行编译安装
```
/bin/bash -c "$(curl -fsSL http://download.hyahm.com/scs.sh)"
```
