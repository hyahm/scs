# scsctl 命令使用

-- 也就是scsctl status 里面的pname 和 name
pname:  是配置文件的name字段
name:   副本名， 也就是scsctl status的 name

> 重载配置文件（重新加载配置文件， 热更新到scs， scs不需要重启）
```
scsctl config reload
```

> 打印debug配置

```
# servers(这是主要是开发用来调试错误的)
# scripts(类似读取scs.yaml配置文件的scripts的字段)
scsctl debug servers | scripts
```

> 禁用获取启用脚本， 这个的主要功能其实保存配置， 解决临时禁用脚本用的

```
 scsctl disable|enable <pname>
```

> 打印副本的环境变量

```
 scsctl env <name> [flag]
```

> 打印副本的日志

```
 scsctl log <name> [flag]
```

> 杀死某个脚本或副本

```
scsctl kill <pname> [name]
```


> 重启某个脚本或副本

```
scsctl restart (<pname> [name]) | --all
```


> 停止某个脚本或副本

```
scsctl stop (<pname> [name]) | --all
```

> 移除某个脚本或副本（同时也移除了配置文件的里面的配置）

```
scsctl remove <pname> [name]
```


> 杀死某个脚本或副本

```
scsctl start [pname] [name]
```


> 杀死某个脚本或副本

```
scsctl status [pname] [name]
```