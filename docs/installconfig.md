# 服务启动配置
go 模板语法
https://golang.google.cn/pkg/text/template/

## 环境变量PNAME=u5, NAME=u5_1， 参考来自配置文件
```
  - name: u5
    # 查看是否存在文件或命令, 不存在就执行install的命令, 执行时存在env定义的环境变量, 服务启动前执行
    preStart
    - path: /home/git  # 这个条件表示目录或文件是否存在(支持{{ .KEY }}语法， 其中 KEY是env里面的key或环境变量)
      install: mkdir /home/git # (支持{{ .KEY }}语法， 其中 KEY是env里面的key或环境变量)
    - command: git     # 这个条件是命令是否存在 二选一(支持{{ .KEY }}语法， 其中 KEY是env里面的key或环境变量)
      install: yum -y install git   # (支持{{ .KEY }}语法， 其中 KEY是env里面的key或环境变量)
    # 特殊用法， 主要用来设置配置文件，原模板文件格式化成最终的配置文件
    - path: dstFile
      template: srcFile  # 这个条件是命令是否存在 二选一(支持{{ .KEY }}语法， 其中 KEY是env里面的key或环境变量)
    cron:
      # 此行含义， 每个月的25号10:10:10 执行一次
      start: "2020-12-25 10:10:10"
      loop: 1
      isMonth: true  # 如果这里是false， 那么没隔1秒执行一次
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
```
文件必须写在可运行项目的根目录， 以另外一个开源项目 ITflow 的实例如下
通过scsctl 命令直接启动服务，并自动挂载到scs中
```bash
scsctl install -f <install_config_file_path>
```
> install.yaml
```yaml
- name: itflow_api
  env:
    # 线上mysql配置
    MYSQL_USER: root
    MYSQL_PASSWORD: 123456
    MYSQL_HOST: 127.0.0.1
    MYSQL_PORT: 3306
    MYSQL_DB: itflow
    LOG_PATH: log/itflow.log
  preStart:
    # 这一项配置的含义是如果不存在bug.ini文件， 那么就将template的文件格式化后写入到bug.ini文件中
    - path: bug.ini
      # 配置文件模板， 支持text/template语法参考上面的 模板语法 文档地址
      template: bug.ini.tpl
    # 如果没打包，就打包文件
    - path: main
      install: go build main.go
  command: ./main
  # 这项是如果git有新的更新，可以直接 scsctl update iflow_api 可以更新的最新版
  update: go build main.go
```
== 细心的同学会发现，少些了dir 家目录的配置， 通过 scsctl install -f <install_config_file_path>，所有目录都是相对家目录操作  
== 同时自动生成 PROJECT_HOME变量，可以在配置文件模板直接使用  
== 表面看生成的配置文件并没有简化，不如直接改配置文件， 但是如果项目是分布式，N多服务都要启动这个服务，那么是不是有用了  