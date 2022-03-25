<!--
 * @Author: your name
 * @Date: 2022-02-27 10:46:37
 * @LastEditTime: 2022-02-27 10:58:05
 * @LastEditors: Please set LastEditors
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /scs/docs/preStart.md
-->
## 启动前执行`preStart`

> 属于`scripts.preStart` 下
```yaml

#判断路径或文件是否存在， 不存在就执行Install的命令或 Template 模板
path: 
#判断命令是否存在， 不存在就执行 Install
command:
#执行命令，判断返回值， 也可搭配  EQ NE GT LT GE LE 根据执行的结果做比较， 注意是根据字符左到右的比较
execCommand:

# 如果要比较版本的话，需要指定分隔符， v1.4.1 (分隔符为.)  v1-43-1 (分隔符为-)
separation: (分隔符， 比较大小才有效)， 与ExecCommand+ EQ | NE | GT | LT | GE | LE搭配使用
eq:  等于
ne:  不等于
gt:  大于
lt:  小于
ge:  大于等于
le:  小于等于
# 如果不满足上面的条件需要执行的脚本
install:
# 仅仅用来判断配置文件是否存在来生成配置文件使用， 
#通过 go text/template 模板渲染, 里面的通过环境变量嵌入
# 比如要使用当前目录及端口
template:

```

## 通过`path`和 `template` 生成配置文件
> 准备一个模板文件   config.py.tpl
```
# encoding=utf-8
aa = {{ .PROJECT_HOME }}
user = {{ .USER }}
port = {{ .PORT }}
level = {{ .TEST_ENV }}
```
> 修改配置文件
```yaml
scripts:
  - name: test
    always: true
    env:
      TEST_ENV: test
    preStart:
      # 这里的意思是， 寻找config.py 文件，如果不存在就根据config.py.tpl生成， 里面的变量可以通过环境变量注入
      - path: config.py 
        template: config.py.tpl
    command: $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
```
> 好的，我们重载配置文件并重启服务, 生成 config.py的内容如下
```
# encoding=utf-8
# 因为我们没有添加 dir， 也就是这里是空的
aa = 
# 这里的no value 是因为没有找到USER这个环境变量
user = <no value>
# PORT 是必定存在的，如果没指定默认是 0
port = 0
# 这个是上面脚本env里面设置的值
level = test
```
> 可以通过下面的命令查看`env`, USER 的变量确实是找不到
```powershell
PS E:\code\scs> scsctl.exe env test_0 | findstr "USER"
USERDOMAIN: DESKTOP-6SSBCHN
USERDOMAIN_ROAMINGPROFILE: DESKTOP-6SSBCHN
USERPROFILE: C:\Users\cande
ALLUSERSPROFILE: C:\ProgramData
USERNAME: cande
-----------------------------------------------
PS E:\code\scs> scsctl.exe env test_0 | findstr "PROJECT_HOME"
PROJECT_HOME:
```

## 判断是否有ruby命令, 
> 因为是命令所以要使用 `command: ruby`, 配置如下

```yaml
scripts:
  # 这个是脚本的名字， 名字必须是字母数字或下划线组成，以后都是根据此名字来操作
  - name: test
    always: true
    env:
      TEST_ENV: test
    priStart:
      command: ruby
    # 这个是执行的基础命令
    command: $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
```

> 重载配置文件并重启服务
```yaml
PS E:\code\scs> scsctl.exe status
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName    Name      Status     Pid      UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   Command
test     test_0    Running    16616    3m9s                 false       0         false     0.13      65296  
   $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
--------------------------------------------------
```

> 我本地是没有ruby命令的， 但是程序还是正常运行的， 因为这个是不会影响`command`的, 而是搭配`install`|`template`而做一些操作， 如果这些操作失败，那么脚本将会启动失败


> 我们添加配置如下, `aaaaa` 这个肯定是执行失败的

```yaml
scripts:
  # 这个是脚本的名字， 名字必须是字母数字或下划线组成，以后都是根据此名字来操作
  - name: test
    always: true
    env:
      TEST_ENV: test
    preStart:
      - command: ruby
        install: aaaaa
    # 这个是执行的基础命令
    command: $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
```

> 再次重载配置文件并重启服务
```yaml
PS E:\code\scs> scsctl.exe status
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName    Name      Status    Pid    UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   
Command
test     test_0    Stop      0      0s                   false       0         false     0.00      0
$n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
```

## 版本是否符合要求 scsctl 的版本是否>=3.7（3.6.1版本以前重启不生效，3.6.1修复）

> 当前写文档的版本
```
PS E:\code\scs> scsctl.exe -v
scsctl version v3.6.0
```

> 因为执行了命令，所以我们需要使用到 `execCommand`, 版本号用`.`来分割的， 所有配置文件的修改如下

```yaml
- execCommand: scsctl.exe -v
  separation: .
  ge: v3.7.0
  # 这里是如果版本<v3.7.0那么就输出  [ERROR] need version >= v3.7.0， 并异常退出, 
  # 当然也可以自己写脚本安装对应版本，这样的话脚本就是正常的running
  install: echo "[ERROR] need version >= v3.7.0"; exit 1;
```
> 日志步骤如下

```
PS E:\code\scs> scsctl.exe config reload
{"code": 200, "msg": "config file reloaded"}
PS E:\code\scs> scsctl.exe restart test 
{"code": 200, "msg": "waiting restart"}
PS E:\code\scs> scsctl.exe status       
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName    Name      Status    Pid    UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   Command    
test     test_0    Stop      0      0s                   false       0         false     0.00      0         $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
--------------------------------------------------
PS E:\code\scs> scsctl.exe log test_0   
2022-03-14 15:40:41 - [INFO] -  - 2022年3月14日 15:40:41
2022-03-14 15:42:01 - [INFO] -  - [ERROR] need version >= v3.7.0
```