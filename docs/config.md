# 配置文件


# 配置文件说明  
```
# 监听端口
listen: :11111
# 脚本日志最大保留行数, 默认100， 一般只会打印异常日志

disableTls: true
key: 
pem: 
# scs 的日志文件
log:
  # scs 的日志文件目录
  path: log
  day: true
  size: 0
# scsctl log 的最大长度
logCount: 100
# 请求头认证， 脚本与服务器之间交互需要 Token： xxxx,  环境变量TOKEN的值为此token的值
token: 
# 客户端免token认证
ignoreToken:
- 127.0.0.1
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

    - path: /home/git  # 这个条件表示目录或文件是否存在
      command: git     # 这个条件是命令是否存在 二选一
      install: yum -y install git
    cron:
      # 此行含义， 每个月的25号10:10:10 执行一次
      start: "2020-12-25 10:10:10"
      loop: 1
      isMonth: true  # 如果这里是false， 那么没隔1秒执行一次
    dir: D:\\work\\u5
    # 是够禁用脚本， 为了保留配置又不想运行显示就启用
    disable: true
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