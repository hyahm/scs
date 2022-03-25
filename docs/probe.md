# 检测器

## 磁盘监控
> 根据上小结的日志结果可以看到这一段

```powershell
2022-03-11 23:52:34 - [INFO] - DESKTOP-6SSBCHN - E:/code/scs/probe/probe.go:242 - alert dist: --C:--, type: NTFS
2022-03-11 23:52:34 - [INFO] - DESKTOP-6SSBCHN - E:/code/scs/probe/probe.go:242 - alert dist: --D:--, type: NTFS
2022-03-11 23:52:34 - [INFO] - DESKTOP-6SSBCHN - E:/code/scs/probe/probe.go:242 - alert dist: --E:--, type: NTFS
2022-03-11 23:52:34 - [INFO] - DESKTOP-6SSBCHN - E:/code/scs/probe/probe.go:242 - alert dist: --F:--, type: NTFS
```
上面显示的是监控的磁盘路径， 目前我电脑是4个盘符， 所以全部都检测到了, 


***如果有需要监控的磁盘但是没显示出来，请提交issue, 目前支持的磁盘格式： EXT4, NTFS, NFS4, XFS, APFS***

## 配置监控项测试

目前可以知道， 我们配置文件没有配置任何跟监控有关的内容 ， 但是还是有监控到硬件， 是因为有默认值， 
参考下面, 为了便于后面操作，我们直接将下面的内容添加到配置文件中
```
probe:
  # 监控哪一台的机器的scs服务
  monitor: 
  # 监控我服务的机器可以免token
  monitored: 
  # mem使用率, 默认90, 小于0表示不启用检测
  mem: 90 
  # cpu使用率, 默认90, 小于0表示不启用检测
  cpu: 90
  # 硬盘使用率， 默认85, 小于0表示不启用检测
  disk: 85        
  # 排除的挂载点， 默认已经去掉了swap， 设备, 数组
  excludeDisk: 
  # 检测间隔， 默认10秒
  interval: 10s
  # 下次报警时间间隔， 如果恢复了就重置
  continuityInterval: 1h
```
> 没错，熟悉的命令又来了， 当然现在是看不到任何效果的， 因为相当于配置文件没有修改过
```powershell
scsctl config reload
```

## 排除监控某些盘符

有些磁盘是不想监控的， 比如`F:`, 那么可以选择性的排除检测, 配置项`excludeDisk`增加一个`F:`, 不区分大小写

```
  excludeDisk: 
    - F:
```

> 重载配置文件, 发现报错了

```powershell
PS E:\code\scs> scsctl.exe config reload
{"code": 500, "msg": "yaml: unmarshal errors:
  line 25: cannot unmarshal !!map into string"}
```
> 看看是否对运行的脚本有影响
```
PS E:\code\scs> scsctl.exe status
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName    Name      Status     Pid      UpTime    Version    CanNotStop  Failed    Disable   CPU       MEM(kb)   Command      
test     test_0    Running    11016    25m43s               false       0         false     0.04      69348     $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
```
**通过结果发现没有影响， 这是因为配置文件里面有特殊符号，我们要用双引号引起来**
> 配置文件修改一下
```
  excludeDisk: 
    - "F:"
```

### 如果修改了配置, 重载配置文件即可生效
```
scsctl config reload
```

### 最后查看一下日志, `F` 盘已经看不到了
```
2022-03-13 10:43:54 - [INFO] - DESKTOP-6SSBCHN - E:/code/scs/probe/probe.go:244 - alert dist: --C:--, type: NTFS
2022-03-13 10:43:54 - [INFO] - DESKTOP-6SSBCHN - E:/code/scs/probe/probe.go:244 - alert dist: --D:--, type: NTFS
2022-03-13 10:43:54 - [INFO] - DESKTOP-6SSBCHN - E:/code/scs/probe/probe.go:244 - alert dist: --E:--, type: NTFS
```