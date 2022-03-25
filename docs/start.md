
# 快速开始

> 先准备一个脚本，请根据自己的系统选择对应的脚本
- windows
先准备一个脚本, 执行结果如下图所示
```powershell
PS E:\code\scs> $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}      
2022年3月11日 22:58:10
2022年3月11日 22:58:20
```

- linux or mac
达到上面类似效果
```bash
[root@localhost scs]# while true;do date; sleep 10; done
Fri Mar 11 09:52:13 EST 2022
Fri Mar 11 09:52:23 EST 2022
```

> 现在我们要将这一小段脚本交给scs来执行(windows为例)
- 默认配置文件是空的， 增加以下配置  
```yaml
scripts:
  # 这个是脚本的名字， 名字必须是字母数字或下划线组成，以后都是根据此名字来操作
  - name: test
    # 这个是执行的基础命令
    command: $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
```

- 我们要让配置生效
  - 方法一: 重启服务（这显然不是一个好方法）
  - 方法一: 直接重新加载配置文件（推荐）

    ```powershell
    # 提示已经加载完成
    PS E:\code\scs> scsctl.exe config reload
    {"code": 200, "msg": "config file reloaded"}
    ```

- 查看结果

```
PS E:\code\scs> scsctl.exe status       
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName    Name      Status     Pid      UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   Command
test     test_0    Running    16128    3s                   false       0         false     7.06      76196     $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}

PName        脚本的唯一标识
Name         副本的唯一标识 默认为脚本名 + "_<index>"
Status       脚本的状态
Pid          pid
UpTime       运行的时间
Version      脚本版本
CanNotStop   是否不能被停止
Failed       失败重启的次数统计
Disable      脚本是否禁用
CPU          此脚本cpu使用率
MEM(kb)      此脚本内存的使用率
Command      启动的命令  文件夹+命令
```

- 根据上面的显示，可以看到是在运行中的， 但是怎么证明在运行呢
  
首先想到的必定是日志(下面的中文乱码先无视)

```powershell
PS E:\code\scs> scsctl.exe log test_0
2022-03-11 23:24:34 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:24:34
2022-03-11 23:24:44 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:24:44
2022-03-11 23:24:54 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:24:54
2022-03-11 23:25:04 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:25:04
2022-03-11 23:25:14 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:25:14
2022-03-11 23:25:24 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:25:24
2022-03-11 23:25:34 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:25:34
2022-03-11 23:25:44 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:25:44
2022-03-11 23:25:54 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:25:54
2022-03-11 23:26:04 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:26:04
2022-03-11 23:26:14 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:26:14
```

