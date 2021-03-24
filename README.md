# scs

service control service or script

# 主要功能
1.  为了保护数据不会在执行中被人工手动中断丢失， 可以让服务在某段时间内才能停止
2.  监控硬件信息， 主要是磁盘， cpu， 内存, scs服务  
3.  服务之间可以相互控制增删改查  
4.  报警功能api  
5.  支持定时器功能执行命令或脚本
6.  客户端控制多台服务器
7.  通过配置文件安装服务
8.  可以将一些执行耗时的脚本托管给scs处理， 以便快速返回结果  

# 文档
[安装](document/install.md)  
[客户端scsctl使用](document/scsctl.md)  
[报警](document/alert.md)  
[服务添加删除接口](document/script.md)  
[硬件监控配置说明](document/hardware.md)  

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

# 配置文件说明  
```
# 监听端口
listen: :11111
# 脚本日志最大保留行数, 默认100， 一般只会打印异常日志
logCount: 100
# scs 的日志文件
log:
  # scs 的日志文件目录
  path: log
  day: true
  size: 0
# 请求头认证， 脚本与服务器之间交互需要 Token： xxxx,  环境变量TOKEN的值为此token的值
token: 
# 报警方式
alert:
  email:
    # 别名
    nickname: "web administrator"
    host: smtp.qq.com
    port: 465
    username:  1654640g46@qq.com
    password: 123456
    # 收件人， 所有报警， 此收件人都会收到
    to:  
      - 727023460@qq.com
  rocket:
    server: https://chat.hyahm.com
    username:  test
    password: 123456
    to:  
      - "#general"
  telegram:
    server: https://chat.hyahm.com
    username: "test"
    password: "123456"
    to: 
      - "-575533567"
  # https://work.weixin.qq.com/help?person_id=1&doc_id=13376#markdown%E7%B1%BB%E5%9E%8B 固定mark格式
  weixin:
    server: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=dd065367-b753-48fb-a974-bbfff0284c1c
# 本地磁盘， cpu， 内存监控项, 就算没写， 也会默认监控, v2版以后此key更改为probe
hardware:
  # 主动监控点 域名： https://127.0.0.1:11111
  monitor: 
    - https://127.0.0.1:11111
  # 被监控ip, 填写后， 这些监控点机器监控此机器的服务无需token 验证
  monitored:
    - 127.0.0.1
  # mem使用率, 默认90, 小于0不检测
  mem: 60 
  # cpu使用率, 默认90, 小于0不检测
  cpu: 90
  # 硬盘使用率， 默认85, 小于0不检测
  disk: 90        
  # 不需要监控的挂载点， 启动的时候日志会打印监控的磁盘，如果类型很常用，请提交issues
  excludeDisk: 
  # 检测间隔， 默认10秒
  interval: 10s
  # 下次报警时间间隔， 如果恢复了就重置
  continuityInterval: 1h
scripts:
  # 环境变量PNAME=u5, NAME=u5_1
  - name: u5
    # 查看是否存在文件或命令, 不存在就执行install的命令, 执行时存在env定义的环境变量
    lookPath
    - path: git
      install: yum -y install git
    cron:
      # 此行含义， 每个月的25号10:10:10 执行一次
      start: "2020-12-25 10:10:10"
      loop: 1
      isMonth: true  # 如果这里是false， 那么没隔1秒执行一次
    dir: D:\\work\\u5
    # 设置环境变量
    env:
      key: value
    # 执行完成后是否删除,  如果想执行的脚本完后自动删掉，可以启用， 多使用于挂载在后台执行
    deleteWhenExit: false
    # 环境变量PORT, 支持变量$PORT, 当replicate大于1时， 命令的$PORT会递增1
    port: 8080
    # 版本号， 此处是一个命令的结果
    version: "scsd -v"
    command: "python .\\test.py signal"
    # replicate， 开启副本数， 默认 1, 如果大于1并且需要特殊条件才能停止， 请在脚本参数后添加 $NAME   
    # 此参数是传递请求需要的name
    replicate: 10
    # 不写默认10分钟
    continuityInterval: 1h
    # 报警收件人， 此脚本额外的收件人
    alert:
      email: 
        - 727023885460@qq.com
      rocket:
        - "@alert"
```

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

# 重启脚本的api
POST /restart/{pname}/{name}
POST /restart/{pname}
POST /restart

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
