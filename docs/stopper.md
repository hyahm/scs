# 停止器

假如下面一段代码，从输入11111开始无法停止 到输出 2333 后才能停止，    
常见做法应该是打印日志， 然后查看日志在打印can stop 的时候杀掉该服务保证上面的需求  
弊端： 需要一直盯着日志， 能停止的时间如果很短的话， 可能一不小心就错过，需要等待另一个可以停止的时间  


```python
import os
import sys
import time
import random
from pyscs import SCS

def log(s):
    print(s) 
    sys.stdout.flush()


    # do something
while True:
    for i in range(10):
        log("can stop in %ss" % (10-i))
        time.sleep(1)

    log(11111)
    time.sleep(random.randint(5, 8))
    log(2333)

```

使用停止器后的代码改进，  
只需要执行 `scsctl stop xxxx ` 即可达到同样的需求    

```python
import os
import sys
import time
import random
from pyscs import SCS

def log(s):
    print(s) 
    sys.stdout.flush()
    
# 注意这里的domain。 默认是 https://127.0.0.1:11111 如果使用http的话需要自己修改
scs = SCS(domain="https://127.0.0.1:11111")

    # do something
while True:
    time.sleep(1)
    scs.can_not_stop()
    log(11111)
    time.sleep(random.randint(5, 8))
    log(2333)
    scs.can_stop()

```

> 核心原理， 通过一个字段保存是否可以被停止的状态， 在可以停止的时候才能被停止
