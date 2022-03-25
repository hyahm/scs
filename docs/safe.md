## 安全配置
```
# 为了支持客户端分布式操作， 接口一般都是对外的， 需要添加请求头验证
token: 
# 客户端免token认证
ignoreToken:
```


## 添加token 步骤

> 添加一个token

```yaml
token: 123456

```
> 重载配置文件, 第一次执行提示加载成功， 第二次执行直接token error
```powershell
PS E:\code\scs> scsctl.exe config reload                                   
{"code": 200, "msg": "config file reloaded"}
PS E:\code\scs> scsctl.exe config reload
{"code": 203, "msg": "token error"}
PS E:\code\scs> scsctl.exe config reload
```

> scsctl 其实也有一个连接服务器的配置， 但是会自动生成默认的配置所以没感觉到， 我们有2种方法解决这个问题
1. 第一种方法
    > 打开scsctl 配置文件， 文件在 \<$HOME\>\.scsctl.yaml, 打开后是这样的其中`nodes`, `local`, `group`在客户端操作那节在讲
    ```
    PS E:\code\scs> cat C:\Users\cande\.scsctl.yaml
    nodes:
    local: 
        url: "https://127.0.0.1:11111"
        token:  
    group: 
    ```
    > token 修改成服务器修改的
    ```
    PS E:\code\scs> cat C:\Users\cande\.scsctl.yaml
    nodes:
    local: 
        url: "https://127.0.0.1:11111"
        token:  123456
    group: 
    ```
    > 再次执行`scsctl.exe config reload`, 这个命令只是为了验证与服务端的连通性，
    ```
    PS E:\code\scs> scsctl.exe config reload
    {"code": 200, "msg": "config file reloaded"}
    ```
----------
2. 第二种方法
    直接在设置`token`的时候, 顺便就把`ignoreToken`给设置了， 这个是设置那些ip可以忽视token验证

    **回想一下，之前是不是也有配置ip忽略验证的， 没错就是监控scs服务的， 配置这个的话，也是可以的**

    > 这个时候 scs
    ```
    token: 123456
    ignoreToken:
    - 127.0.0.1
    ```
    > 重载一下配置文件, 本地scs配置没有修改也是可以的
    ```
    PS E:\code\scs> scsctl.exe config reload
    {"code": 200, "msg": "config file reloaded"}
    PS E:\code\scs> scsctl.exe config reload
    {"code": 200, "msg": "config file reloaded"}
    ```
    > scsctl的配置文件
    ```
    nodes:
    local: 
        url: "https://127.0.0.1:11111"
        token:  
    group: 
    ```

!> token 不设置是非常不安全的