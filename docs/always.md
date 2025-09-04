程序异常退出自动拉起配置
`always: true`

> 我们将原来的脚本改成异常的并设置成`always: true`, 现在的配置是这样

```yaml
scripts:
  # 这个是脚本的名字， 名字必须是字母数字或下划线组成，以后都是根据此名字来操作
  - name: test
    # 异常退出的，1秒后自动拉起
    always: true
    # 将之前的$去掉， 在报警案例中我们用过， 会异常报警
    command: n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
```
> 重载配置文件
> 重启服务(不要忘记这一步， 重载配置文件不会影响运行的服务)
```
PS E:\code\scs> scsctl.exe config reload
{"code": 200, "msg": "config file reloaded"}
```
> 查看状态

```powershell
PS E:\code\scs> scsctl.exe status       
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName    Name      Status    Pid    UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   
Command
test     test_0    Stop      0      0s                   false       1       false     0.00      0
n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
--------------------------------------------------
```
> 查看状态(等待几秒继续)

```powershell
PS E:\code\scs> scsctl.exe status       
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName    Name      Status    Pid    UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   
Command
test     test_0    Running      0      0s                   false       6       false     0.00      0
n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
------
```

**注意`Failed`字段， 这个是重启的次数记录**