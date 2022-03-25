## 配置文件完整版说明  
```
# 监听端口
listen: :11111
# ssl 配置
disableTls: true
key: 
pem: 
# scs 的日志文件
log:
  # scs 的日志文件路径
  path: log/scs.log
  # 每天切割日志
  day: true
  # 按照文件大小来切割日志
  size: 0
  # 清除超过多长时间的日志， 默认不清除日志
  clean: 
# 请求头认证， 脚本与服务器之间交互需要 Token： xxxx,  环境变量TOKEN的值为此token的值
token: 
# 客户端免token认证
ignoreToken:
- 127.0.0.1
# 报警器
alert:
  email:
    # 别名不一定生效
    nickname: "scs"
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
    to: 
      - "-575533567"
  # https://work.weixin.qq.com/help?person_id=1&doc_id=13376#markdown%E7%B1%BB%E5%9E%8B 固定mark格式
  weixin:
    server: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  callback:
    # 接受请求的url
    urls: 
      - http://127.0.0.1:8080
    # 请求方式
    method:  POST
    # 请求头
    headers:
        Content-Type: 
          - application/json

# 本地磁盘， cpu， 内存监控项, 就算没写， 也会默认监控, v2版以后此key更改为probe
probe:
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
  excludeDisk: []
  # 检测间隔， 默认10秒
  interval: 10s
  # 下次报警时间间隔， 如果恢复了就重置
  continuityInterval: 1h
scripts:
  # 环境变量PNAME=u5, NAME=u5_1
  - name: u5
    # 查看是否存在文件或命令, 不存在就执行install的命令, 执行时存在env定义的环境变量, 服务启动前执行
    preStart
    - path: /home/git  # 这个条件表示目录或文件是否存在(支持text/template语法， 其中 KEY是env里面的key或环境变量)
      install: mkdir /home/git # (支持text/template语法， 其中 KEY是env里面的key或环境变量)
    - command: git     # 这个条件是命令是否存在 二选一(支持text/template语法， 其中 KEY是env里面的key或环境变量)
      install: yum -y install git   # (支持text/template语法， 其中 KEY是env里面的key或环境变量)
    # 特殊用法， 主要用来设置配置文件，原模板文件格式化成最终的配置文件
    - path: dstFile
      template: srcFile  # 这个条件是命令是否存在 二选一(支持text/template语法， 其中 KEY是env里面的key或环境变量)
    cron:
      # 此行含义， 每个月的25号10:10:10 执行一次， 不填就是当前时间循环
      start: "2020-12-25 10:10:10"
      # 循环间隔时间， 必填
      loop: 1  
      isMonth: true  # 如果这里是false， 那么没隔1秒执行一次
      times: 10  # 循环的次数， 0就是无限循环
    dir: D:\\work\\u5
    # 是够禁用脚本， 为了保留配置又不想运行显示就启用
    disable: true
    # 设置环境变量,key全为大写
    env:
      key: value
    # 执行完成后是否删除,  如果想执行的脚本完后自动删掉，可以启用， 多使用于挂载在后台执行
    deleteWhenExit: false
    # 环境变量PORT, 支持变量$PORT, 当replicate大于1时， 副本环境变量PORT会递增1
    port: 8080
    # 版本号， 此处是一个命令的结果
    version: "scsd -v"
    # (支持{{ .KEY }}语法， 其中 KEY是env里面的key或环境变量)
    command: "python .\\test.py signal"
    # replicate， 开启副本数， 默认 1, 如果大于1并且需要特殊条件才能停止， 请在脚本参数后添加 $NAME   
    # 此参数是传递请求需要的name
    replicate: 10
    # 不写默认1小时
    continuityInterval: 1h
    update: "git pull"
    # 报警收件人， 此脚本额外的收件人
    alert:
      email: 
        - 727023885460@qq.com
      rocket:
        - "@alert"
```

## 修改监听地址

> 可以修改为 :22222, 或者只允许本地的 127.0.0.1:22222
```
listen: :22222

```

## https证书配置
```yaml
# ssl 配置, 不配置key， pem 是不安全的, 里面文件的路径是scsd执行目录的相对路径
disableTls: true
key: 
pem: 
```

**上面4项修改必须要重启scsd服务才能生效，也是唯一需要重启服务生效的配置项**