大家查看状态的时候不知道有没有看到最后一个`command` (3.6.7版本移除了Disable)
```
scsctl.exe status
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName    Name      Status    Pid    UpTime    Version    CanNotStop  Failed      CPU       MEM(kb)   Command
test     test_0    Stop      0      0s                   false       0           0.00      0         $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
--------------------------------------------------
```
这里很清楚的显示了这个服务的位置，并且显示了执行的命令  

前面我们知道，如果在`scsctl config reload`的时候，会自动启动脚本  

有这么一个场景， 我们想要保留这个配置项，又不想让他运行， 后面说不定需要运行

那么我们可以将此脚本禁止掉， 那么`scsctl config reload`和`scsctl start`就不会运行这个脚本

案例如下，我将上面的例子`test`恢复正常
```
scripts:
- name: test
  command: $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
```
> 重载配置文件查看状态
```
PS E:\code\scs> scsctl.exe config reload
{"code": 200, "msg": "config file reloaded"}

PS E:\code\scs> scsctl.exe status       
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName    Name      Status     Pid     UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   Command
test     test_0    Running    5264    4s                   false       0         false     3.72      63424     $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
--------------------------------------------------
```

> 现在我暂时不想让他运行了
```
PS E:\code\scs> scsctl.exe disable test
{"code": 200, "msg": "waiting stop"}

PS E:\code\scs> scsctl.exe status
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName    Name      Status    Pid    UpTime    Version    CanNotStop  Failed   CPU       MEM(kb)   Command
--------------------------------------------------
```

> 配置文件变化(多了一个 `disable: true`)
```
scripts:
- name: test
  command: $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
  disable: true
  liveness: null
```

**状态里面已经看不到了**

> 后面我又想让它跑起来了
```
PS E:\code\scs> scsctl.exe enable test
{"code": 200, "msg": "waiting start"}

PS E:\code\scs> scsctl.exe status     
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName    Name      Status     Pid      UpTime    Version    CanNotStop  Failed    CPU       MEM(kb)   Command
test     test_0    Running    10000    1s                   false       0         11.78     63356  
   $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
--------------------------------------------------
PS E:\code\scs>
```

> 配置文件变化还原了
```
scripts:
- name: test
  command: $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
  liveness: null
```