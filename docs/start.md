
# 快速开始

> 第一步： 这是我要启动的脚本
```python
# encoding=utf-8
import os
import sys
import time

def log(s):
    print(s)
    sys.stdout.flush()

    # do something
while True:
    log("end")
    time.sleep(1)
```


> 第二部 修改配置文件如下
```yaml
scripts:
  - name: test
    dir: /root
    command: python3 test.py
```

> 使用命令执行
```bash
scsctl config reload
```

> 查看结果
```
[root@node1 scs]# scsctl status
<node: localhost, url: https://127.0.0.1:11111>
--------------------------------------------------
PName   Name      Status     Pid     UpTime   Verion   CanNotStop  Failed    Disable   CPU       MEM(kb)   Command
test    test_0    Running    2345    2m51s              false       0         false     0.14      5784      cd /root && python3 test.py

PName        脚本的唯一标识
Name         副本的唯一标识 默认为脚本名 + "_<index>"
Status       脚本的状态
Pid          pid
UpTime       运行的时间
Verion       脚本版本
CanNotStop   是否不能被停止
Failed       失败重启的次数统计
Disable      脚本是否禁用
CPU          此脚本cpu使用率
MEM(kb)      此脚本内存的使用率
Command      启动的命令  文件夹+命令
```

### 程序异常退出自动拉起配置
`always: true`