## 多副本配置

**隶属于 script**

> 启动10个脚本， 没有监听端口的
```
replicate: 10
```

> 启动10个服务， 带了监听端口的必须指定`port`, 不然有9个起不来
```
# 设置port环境变量， 每个副本端口号自动+1， 如果遇到被使用的端口自动跳过+1
port: 8000
replicate: 10
# 设置可以通过${PORT}直接传递给启动脚本， 也可以代码通过获取环境变量 PORT来设置端口号
command:  python3 test.py --port {{ .PORT }}
```

> 为了测试，我们增加一个脚本`scripts里面应该是有2个脚本`
```yaml
scripts:
  # 这个是脚本的名字， 名字必须是字母数字或下划线组成，以后都是根据此名字来操作
  - name: test
    always: true
    env:
      TEST_ENV: test
    preStart:
      - path: config.py 
        template: config.py.tpl
      - execCommand: scsctl.exe -v
        separation: .
        ge: v3.7.0
        install: echo "[ERROR] need version >= v3.7.0"; exit 1;
    # 这个是执行的基础命令
    command: $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
  - name: replicate
    # 这个是执行的基础命令
    replicate: 10
    port: 18000
    command: echo "{{ .PORT }}"
```

> 重载查看状态（这次我们增加了脚本，重载的时候自动全部启动，不需要`scsctl start`）
```
PS E:\code\scs> scsctl.exe config reload
{"code": 200, "msg": "config file reloaded"}
PS E:\code\scs> scsctl.exe status
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName        Name           Status    Pid    UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   Command
replicate    replicate_0    Stop      0      0s                   false       0         false     0.00      0      
   echo "18000"
replicate    replicate_1    Stop      0      0s                   false       0         false     0.00      0      
   echo "18001"
replicate    replicate_2    Stop      0      0s                   false       0         false     0.00      0      
   echo "18002"
replicate    replicate_3    Stop      0      0s                   false       0         false     0.00      0      
   echo "18003"
replicate    replicate_4    Stop      0      0s                   false       0         false     0.00      0      
   echo "18004"
replicate    replicate_5    Stop      0      0s                   false       0         false     0.00      0      
   echo "18005"
replicate    replicate_6    Stop      0      0s                   false       0         false     0.00      0      
   echo "18006"
replicate    replicate_7    Stop      0      0s                   false       0         false     0.00      0      
   echo "18007"
replicate    replicate_8    Stop      0      0s                   false       0         false     0.00      0      
   echo "18008"
replicate    replicate_9    Stop      0      0s                   false       0         false     0.00      0      
   echo "18009"
test         test_0         Stop      0      0s                   false       0         false     0.00      0      
   $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
--------------------------------------------------
```

**注意看，一共有10个副本分别从0-9, 端口号依次+1**

## 副本扩缩容

> 有的时候副本数太多，我们需要减少数量, 比如我们只要5个
```
# 我们将配置文件改为5
replicate: 5
```

> 重载配置,再查看
```powershell
PS E:\code\scs> scsctl.exe config reload
{"code": 200, "msg": "config file reloaded"}
PS E:\code\scs> scsctl.exe status
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName        Name           Status    Pid    UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   Command
replicate    replicate_0    Stop      0      0s                   false       0         false     0.00      0      
   echo "18000"
replicate    replicate_1    Stop      0      0s                   false       0         false     0.00      0      
   echo "18001"
replicate    replicate_2    Stop      0      0s                   false       0         false     0.00      0      
   echo "18002"
replicate    replicate_3    Stop      0      0s                   false       0         false     0.00      0      
   echo "18003"
replicate    replicate_4    Stop      0      0s                   false       0         false     0.00      0      
   echo "18004"
test         test_0         Stop      0      0s                   false       0         false     0.00      0      
   $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
--------------------------------------------------
```

**已经停止了5个，是按照`index`删除最后生成的副本**


## 删除某个副本(假如我要删除replicate_2的副本数)

```powershell
PS E:\code\scs> scsctl.exe remove replicate replicate_2
{"code": 200, "msg": "waiting stop"}
PS E:\code\scs> scsctl.exe status
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName        Name           Status    Pid    UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   Command
replicate    replicate_0    Stop      0      0s                   false       0         false     0.00      0      
   echo "18000"
replicate    replicate_1    Stop      0      0s                   false       0         false     0.00      0      
   echo "18001"
replicate    replicate_3    Stop      0      0s                   false       0         false     0.00      0      
   echo "18003"
replicate    replicate_4    Stop      0      0s                   false       0         false     0.00      0      
   echo "18004"
test         test_0         Stop      0      0s                   false       0         false     0.00      0      
   $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
--------------------------------------------------
```

**注意看，已经删除了replicate_2的副本**

> 我们看看配置文件有没有变化(3.6.1修复了这个问题， 不然这里显示的是9)
```
- name: replicate
  command: echo "{{ .PORT }}"
  replicate: 4
  port: 18000
  liveness: null
```