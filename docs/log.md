# 日志


## 日志查看
由上面的快速启动的第一个脚本引出了最后的日志
```powershell
PS E:\code\scs> .\scsctl.exe log test_0
2022-03-11 23:24:34 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:24:34
2022-03-11 23:24:44 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:24:44
2022-03-11 23:24:54 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:24:54
```

如果你不`ctrl+c`会发现，日志是实时的(默认显示最后10行)


## 显示末尾20行的日志

在你上面的脚本跑了好久后，后面加参数20就是最后20行 

```powershell
PS E:\code\scs> .\scsctl.exe log test_0 20
2022-03-11 23:29:54 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:29:54
...
2022-03-11 23:33:04 - [INFO] -  - E:/code/scs/server/log.go:51 - 2022��3��11�� 23:33:04
```

## 日志文件

直接这样显示文件不方面查找， 那么我们可以直接打开文件查看  

文件位置默认在 `log/test_0.log` 这里， 相对于执行目录的为位置

为了避免日志文件太大导致 `.\scsctl.exe log <name>` 读取慢， 默认10M就切割日志，时间戳为旧文件前缀

## 日志配置
```yaml
log:
  # scs服务 的日志文件路径，脚本日志文件继承此处设置的目录，
  path: log/scs.log
  # 每天切割日志
  day: true
  # 按照文件大小来切割日志
  size: 0
  # 清除超过多长时间的日志， 默认不清除日志
  # 脚本日志文件继承此处配置， 单位  天
  clean: 30
```

> 我们先添加进去看看， 建议将 `scripts` 字段放到最后
```yaml
log:
  # scs 的日志文件路径，脚本日志文件继承此处设置的目录，
  path: log/scs.log
  # 每天切割日志
  day: true
  # 按照文件大小来切割日志
  size: 0
  # 清除超过多长时间的日志， 默认不清除日志
  # 脚本日志文件继承此处配置， 单位支持 h|m|s (时分秒)
  clean: 720h

scripts:
  # 这个是脚本的名字， 名字必须是字母数字或下划线组成，以后都是根据此名字来操作
  - name: test
    # 这个是执行的基础命令
    command: $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}

```
> 加载配置文件
```powershell
PS E:\code\scs> .\scsctl.exe config reload
{"code": 200, "msg": "config file reloaded"
```

> 查看效果
- 查看`log`目录下已经生成了 `scs.log` 文件, 内容如下， 说明修改日志配置通过`scsctl config reload` 同样生效
```
2022-03-11 23:52:34 - [INFO] - DESKTOP-6SSBCHN - E:/code/scs/probe/probe.go:242 - alert dist: --C:--, type: NTFS
2022-03-11 23:52:34 - [INFO] - DESKTOP-6SSBCHN - E:/code/scs/probe/probe.go:242 - alert dist: --D:--, type: NTFS
2022-03-11 23:52:34 - [INFO] - DESKTOP-6SSBCHN - E:/code/scs/probe/probe.go:242 - alert dist: --E:--, type: NTFS
2022-03-11 23:52:34 - [INFO] - DESKTOP-6SSBCHN - E:/code/scs/probe/probe.go:242 - alert dist: --F:--, type: NTFS
2022-03-11 23:52:34 - [INFO] - DESKTOP-6SSBCHN - E:/code/scs/server/config.go:175 - oldReplicate: 1test
2022-03-11 23:52:34 - [INFO] - DESKTOP-6SSBCHN - E:/code/scs/server/config.go:176 - newReplicate: 1test
```
