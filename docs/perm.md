# 远程权限管理

如下配置, 我想把test分享给开发查看， 我们可以在对应的脚本下加一个token
```
token: aosdhf-90y890Y*(G0(&TG0*G9F)78fg87)
ignoreToken:
- 127.0.0.1

scripts:
  - name: remote
    command: while true;do date; sleep 5; done

  - name: test
    dir: E:\code\scs
    token: 123
    command: python test.py
    update: "git pull"
```

> 将此token告知对应开发的同事, scsctl 的配置
```
nodes:
  test:
    url: https://192.168.101.11:11111
    token: 123
```

> 可以看到日志
```
[root@master scs]# scsctl log test_0
2022-03-27 21:04:06 - [INFO] - DESKTOP-6SSBCHN - this is test token can stop
2022-03-27 21:04:06 - [INFO] - DESKTOP-6SSBCHN - this is test token can not stop
2022-03-27 21:04:11 - [INFO] - DESKTOP-6SSBCHN - this is test token can stop
2022-03-27 21:04:11 - [INFO] - DESKTOP-6SSBCHN - this is test token can not stop
2022-03-27 21:04:14 - [INFO] - DESKTOP-6SSBCHN - this is test token can stop
```


> 更新操作
```shell
[root@master scs]# scsctl update test
{"code":200,"msg":"waiting update","role":"look"}


[root@master scs]# scsctl log test_0
2022-03-27 21:08:17 - [INFO] - DESKTOP-6SSBCHN - this is test token can stop-----
2022-03-27 21:08:17 - [INFO] - DESKTOP-6SSBCHN - this is test token can not stop------
2022-03-27 21:08:20 - [INFO] - DESKTOP-6SSBCHN - this is test token can stop-----
```

> look可以操作的命令如下
```
  env         show env
  get         get info
  help        Help about any command
  info        show service
  log         script log
  restart     restart assign script
  start       start assign script
  status      Print assign script status
  stop        stop script
  update
```

