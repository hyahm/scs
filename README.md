# scs

service control service or script  
类似supervisor,但是更高级，支持所有系统   
自带监控及通知     
服务控制脚本能否停止 最大程度防止脚本数据丢失   
码云地址: https://gitee.com/cander/scs

# 适用适用的

场景一:  服务器需要监控报警cpu，内存, 磁盘，但是主要是要给我实时报警，以便提前避免不必要的事故
```
只要安装启动就已经监控的， 如需报警，需要添加报警器
下面是完整配置的参考 `/etc/scs.yaml` 详细用法参考文档
```
alert:
  email:
    host: smtp.qq.com
    port: 465
    username:  165464646@qq.com
    password: 123456
    to:
      - 727023460@qq.com
  rocket:
    server: https://chat.hyahm.com
    username:  test
    password: 123456
    to:
      - "#general"
  telegram:
    server: https://telegram.hyahm.com:8989
    to:
      - "-789789435"
  # https://work.weixin.qq.com/help?person_id=1&doc_id=13376#markdown%E7%B1%BB%E5%9E%8B 固定mark格式
  weixin:
    server: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=dd065367-b753-48fb-a974-bbfff0284c1c
  # 这个回调决定支持所有类型报警， 需要自己写
  callback:
    # 接受请求的url
    urls:
      - http://192.168.0.112:8080
    # 请求方式
    method:  POST
    headers:
      Content-Type:
        - application/json
```
```


场景二:  一个二进制文件想要执行，需要工具管理起来，而不用手动进入目录启动和找pid停止

```
# 最少配置， 其实只需要name和command， 相对于supervisor， 配置更简单, 执行`scsctl config reload`即可加载而不影响其他的配置
scripts:
- name: test
  dir: D:\scs
  command: python test.py
```

场景三:  想要一个定时器，但是系统自带的太麻烦，而且精确度不高(更多详细的配置请参考文档)

```
# 每3秒执行一次
scripts:
- name: test
  dir: D:\scs
  command: python test.py
  cron:
    loop: 3
```


场景四:  执行一段队列处理代码，但是需要保证数据处理完成后停止，也就是我在执行`stop` 信号后等处理完成后才会停止，防止队列中的这条数据丢失，并支持超时机制， 如果处理队列的一个请求因为位置原因导致卡住，如果不处理会降低分布式集群的效率，这时候需要自动重启服务，并通过参数传递来做响应的回滚操作

```
# 停止器
本身是通过代码http接口请求实现的，目前执行python和go的sdk 详细的请参考文档

```


场景五:  批量升级更新服务， 但是觉得ansible管理太麻烦, 可以在某一台控制机器直接执行 `update <panme> <name>`升级服务，
```
需要在客户端配置所有scsd的信息， 可以通过 -n 和 -g 来自定义管理的节点或组， 详细请参考文档
```

场景六:  降低运维开发的沟通成本， 运维给与开发的权限，方便开发远程调试，又保证服务器权限
```
请参考权限的文档
```




[文档地址](https://www.itjingtu.com/document/13)
[部分视频教程](https://www.bilibili.com/video/BV1bv411C7Qz/)

具体更新的内容请查看 [update.md](update.md)文件

QQ群:  346746477

