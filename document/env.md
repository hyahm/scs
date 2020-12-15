# 环境变量

挂载上的服务每个里面必定存在下面4个环境变量， 请避免使用到

PNAME // 等于配置文件name的值
NAME  // 由系统自动生成
TOKEN  // 从配置文件中读取
PORT  // 等于配置文件port的值, 环境变量中需要PORT 也必须在port 设置， 环境变量中设置的PORT 无效

可以参考 test.py

### 显示某个脚本中的环境变量
```
scsctl env  name
```