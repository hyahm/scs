[安装](install.md)  
[客户端scsctl使用](scsctl.md)  
[报警](alert.md)  
[服务添加删除接口](script.md)  
[硬件监控配置说明](hardware.md)

# 报警服务

目前支持邮箱， rocket.chat,  telegram
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
    username:  test
    password: 123456
    to:
      - "-789789435"
```
### 自建发件邮箱服务器
参考： https://gitee.com/cander/maddy

### 自建rocket.chat
参考官网： https://rocket.chat/install/?gclid=undefined
> 注意点
- nginx代理的时候需要加上websocket支持
```
# the Meteor / Node.js app server
server {
  server_name yourdomain.com;

  access_log /etc/nginx/logs/yourapp.access;
  error_log /etc/nginx/logs/yourapp.error error;

  location / {
    proxy_pass http://localhost:3000;

    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header Host $host;  # pass the host header - http://wiki.nginx.org/HttpProxyModule#proxy_pass

    proxy_http_version 1.1;  # recommended with keepalive connections - http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_http_version

    # WebSocket proxying - from http://nginx.org/en/docs/http/websocket.html
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection $connection_upgrade;
  }

}
```
- mongodb的配置文件用yaml格式的， ini的配置格式是起不来rocket.chat 服务的

### telegram

如果你能下载的话， 准备一台国外服务器
telegram 有socks5 代理
```
go run cmd/proxy/socks5.go  # 默认监听:8080端口
```
然后增加一个发送消息的代理  
使用的是机器人， 将机器人拉到报警的群里就可以了  
```go run cmd/proxy/telegram.go -l :8989
-l  flag.String("l", "", "listen default :8080")
-u flag.String("u", "", "username") 与scs配置文件对应
-p flag.String("p", "", "password") 与scs配置文件对应
-i flag.String("i", "", "bot send message api // https://api.telegram.org/bot<token>/sendMessage")
) 与scs配置文件对应
```


### 报警接口
POST  /alert  
{
    "title": "haitenda",
    "pname": "asdgasdgasdgfasdf",
    "name": "01358072011",
    "reason": " asdfasdg  \n ashdfljsdf",
    "broken": true,  // 如果恢复了就设置成false, name 和panme 必须 其他的无视
    "interval": 100,
    "to": {
        "rocket": ["#general"]
    }
}

> curl 请求
```
curl -X POST  https://127.0.0.1:11111/set/alert -d '{"title":"test","pname":"alert","name":"test","reason":"test", "broken": false, "interval":100,"to":{"telegram":["-64168425"]}}' -k`
```
> crontab， 定时器操作时报警
```
<script>; if [[ $? -ne 0 ]]; then /usr/bin/curl -X POST  https://127.0.0.1:11111/set/alert -d '{"title":"test","pname":"alert","name":"test","reason":"test", "broken": true, "interval":100,"to":{"telegram":["-346345"]}}' -k; sleep 10; /usr/bin/curl -X POST  https://127.0.0.1:11111/set/alert -d '{"title":"test","pname":"alert","name":"test","reason":"test", "broken": false, "interval":100,"to":{"telegram":["-346345"]}}' -k; fi
```

