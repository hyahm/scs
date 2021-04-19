# 检测器


服务启动启动后默认就监控了 cpu, 内存, 硬盘
> 这里是配置文件关于硬件监控项配置
```
 # mem使用率, 默认90, -1 表示禁用
  mem: 60
  # cpu使用率, 默认90, -1 表示禁用
  cpu: 90
  # 硬盘使用率， 默认85, -1 表示禁用
  disk: 80
  # 排除的挂载点
  excludeDisk:
    - xxxx
  # 检测间隔， 默认10秒
  interval: 10s
  # 下次报警时间间隔， 如果恢复了就重置
  continuityInterval: 1h
```

### 如果修改了配置, 重载配置文件即可生效
```
scsctl config reload
```