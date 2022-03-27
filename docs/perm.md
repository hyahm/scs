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