# scsctl 命令使用

下面标记的 \<pname\>\<name\>  
`pname`:  是配置文件的name字段  
`name`:   副本名， 也就是scsctl status的 name  
```powershell
PS E:\code\scs> scsctl.exe status
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName    Name      Status     Pid      UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   Command
test     test_0    Running    19416    10m59s               false       0         false     0.08      69504  
   $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
```
**如上面的状态  `pname`就是test ---- `name`是test_0**

> 重载配置文件（重新加载配置文件， 热更新到scs， scs不需要重启）前部分经常提到的命令
```
scsctl config reload
```

> 查看基础配置

```
# servers(这是主要是开发用来调试错误的)
# scripts(类似读取scs.yaml配置文件的scripts的字段)
scsctl get servers | scripts
```

> 禁用获取启用脚本， 这个的主要功能其实保存配置， 解决临时禁用脚本用的

```
 scsctl disable|enable <pname>
```

> 打印副本的环境变量

```
 scsctl env <name> [flag]
```

> 打印副本的日志(n是个整数，显示最后多少行数据)

```
 scsctl log <name> [n][flag]
```

> 杀死某个脚本或副本

```
scsctl kill <pname> [name]
```


> 重启某个脚本或副本

```
scsctl restart (<pname> [name]) | --all | <pname>[,<pname>]
```


> 停止某个脚本或副本

```
scsctl stop (<pname> [name]) | --all | <pname>[,<pname>]
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

> 更多操作说明
```
PS E:\code\scs> scsctl.exe
scs service help, version: v3.6.0

Usage:
  scsctl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      scs server config
  disable     disable script
  enable      enable script
  env         show env
  get         get info
  help        Help about any command
  install     install package
  kill        kill script
  log         script log
  remove      remove script
  restart     restart assign script
  search      search package(future)
  start       start assign script
  status      Print assign script status
  stop        stop script
  tar         archive tar package(future)
  update      update server
  upload      upload package(future)

Flags:
  -g, --group string   show which groupname
  -h, --help           help for scsctl
  -n, --node string    show which nodes
  -v, --version        version for scsctl

Use "scsctl [command] --help" for more information about a command.
```

future: 未来可能会添加的命令
unStable: 不稳定命令