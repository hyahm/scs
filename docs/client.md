<!--
 * @Author: your name
 * @Date: 2022-01-21 23:33:10
 * @LastEditTime: 2022-02-27 11:04:55
 * @LastEditors: your name
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /scs/docs/client.md
-->


## 客户端配置文件

第一次启动的时候，是没有做任何操作的， 默认会生成一个配置文件`<$HOME>/.scsctl.yaml`， 如下
```yaml
nodes:
  localhost: 
    url: https://127.0.0.1:11111
    token:
group:
  local:
    - localhost
```

!> 除了通过脚本一键安装外，默认都是不会生成token的， 也是不安全的

## scsctl 命令参数说明
```yaml
- nodes: # 节点
  localhost: # 节点名
    url: https://127.0.0.1:11111  # 节点的完整域名地址
    token:   # token
group:  # 节点组
  local:  # 节点组名
    - localhost # 节点名
```

客户端可以管理多个节点， 而且可以根据节点名进行分组

## 我要通过我本地的scsctl管理我本地和虚拟机的linux

还记得我本机和虚拟机的ip吗， 修改本地配置文件`<$HOME>/.scsctl.yaml`如下
```
nodes: # 节点
  - localhost: # 节点名
    url: http://127.0.0.1:11111  # 节点的完整域名地址
    token: "123456" # token
  - linux:
    url: https://192.168.101.12:11111  # 节点的完整域名地址
    # 虚拟机是通过脚本一键安装的，所以会随机生成一个token
    token:  'gSstfcoK&3,5LQbOCR|8k+8Ftp#TJ}' # token
group:  # 节点组
  local:  # 节点组名
    - localhost # 节点名
```
> 查看一下状态, 2个node，名字分别是`localhost`, `linux`

```powershell
PS E:\code\scs> scsctl.exe status
<node: localhost, url: https://127.0.0.1:11111>
--------------------------------------------------
PName    Name      Status     Pid      UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   Command
test     test_0    Running    10944    18m22s               false       0         false     0.14      70624     $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
--------------------------------------------------
<node: linux, url: https://192.168.101.12:11111>
--------------------------------------------------
PName    Name    Status    Pid    UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   Command
--------------------------------------------------
```

## 配置了多个node 还是想要显示部分的

> 我们配置不动， 注意看, 
```
group:  # 节点组
  local:  # 节点组名
    - localhost # 节点名
```
> 我们只要看我们本地的服务, 因为只有一个节点，我们可以直接指定显示的节点
```
PS E:\code\scs> scsctl.exe status -n localhost
<node: localhost, url: https://127.0.0.1:11111>
--------------------------------------------------
PName    Name      Status     Pid      UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   Command
test     test_0    Running    10944    21m59s               false       0         false     0.12      70624     $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
--------------------------------------------------
```
> 如果要查看多个节点必须要用`-g`来指定， 并在对应组名下列出要显示的节点名
```
PS E:\code\scs> scsctl.exe status -g local
<node: localhost, url: https://127.0.0.1:11111>
--------------------------------------------------
PName    Name      Status     Pid      UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   Command
test     test_0    Running    10944    23m13s               false       0         false     0.12      70628     $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
--------------------------------------------------
```

!> 操作的所有命令与操作本地的都是一致的， 因为请求是异步的，所以可以上千台操作， 但是你应该不会认真去看1000台的返回结果的

!> 如果你看到这里来了， 恭喜你，运维操作的事已经毕业了， 后面基本是开发的工作了，但是我还是建议你看下去， 因为主要功能的第一条还没讲