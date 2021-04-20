# 多副本配置

**隶属于 script**

> 启动10个脚本
```
replicate: 10
```

> 启动10个服务， 无法带了监听端口
```
# 设置port环境变量， 每个副本端口号自动+1， 如果遇到被使用的端口自动跳过+1
port: 8000
replicate: 10
# 设置可以通过${PORT}直接传递给启动脚本， 也可以代码通过获取环境变量 PORT来设置端口号
command:  python3 test.py --port ${PORT}
```