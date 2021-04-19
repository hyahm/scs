

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
# scsctl 命令参数说明

 这里的pname, name 与scsctl status 中的对应  
 尖括号表示必须  
 中括号表示可选  
 |  表示多选一
```
scsctl --help
  config      scs server config  
  env         env
  help        Help about any command
  install     install package(保留命令)
  kill        kill script
  log         script log
  restart     restart assign script
  search      search package(保留命令)
  start       start assign script
  status      Print assign script status
  stop        stop script
  tar         archive tar package(保留命令)
  upload      upload package(保留命令)
```
###### 查看服务状态信息
```
scsctl status
```

###### 重新加载配置文件
```
scsctl config reload
```
###### 显示某服务中的环境变量
```
scsctl env <name>
```
###### 查看某服务中的日志
```
scsctl log <name>
```
###### 强制停止某个服务，不管是否允许停止
```
scsctl kill <pname> [name]
```

###### 启动服务
```
scsctl start [pname] [name]
```

###### 停止服务
```

scsctl stop <pname>|<--all> [name] 
```

###### 重启服务
```
scsctl restart <pname>|<--all> [name] 
```