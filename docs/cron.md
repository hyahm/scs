# 定时器

**隶属于script，详细位置请参考配置文件说明**

> 脚本从启动时间开始每一分钟循环一次

```
cron:
  loop: 60
```

> 脚本从未来某时刻开始每一分钟循环一次

```
cron:
  # 从今天下午13点0分34秒开始执行 **注意0不能省略**
  start: "13:00:04"  
  loop: 60
```

> 脚本从未来某时刻开始每一分钟循环一次

```
cron:
  # 从2021年9月12号13点0分34秒开始执行 **注意0不能省略**
  start: "2021-09-12 13:00:04"  
  loop: 60
```

> 每月1号执行一次

```
cron:
  # 每个月1号执行一起，2021年9月1号开始 **注意0不能省略**
  start: "2021-09-01 00:00:00"  
  loop: 1
  isMonth: true
```

修改完配置 `scsctl config reload`  生效