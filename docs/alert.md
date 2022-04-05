
# 报警器
目前支持
- 邮箱
- rocket.chat
- telegram
- 企业微信
- 自定义接口(3.0后新增)

建议只需要选中一种即可， 根据公司情况来决定

下面是完整配置的参考
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

## 邮箱

只要邮箱授权了SMTP权限即可发送， 配置项为通用配置， 不做过多解释
不过因为邮箱的限制， 可能会失败


## rocket.chat 

```
这个配置更简单
server: 就是要客户端连接服务器填的url
// 账号密码就是登录rocket.chat 的账号密码
username:  test
password: 123456
to: 这个要注意不要忽略前面的符号， 群组前面一般都有个#
```

自建服务参考官网文档： https://rocket.chat/install


## telegram

因为某墙的原因， 无法直接连接， 但是因为特殊是可以用到的  

1. 必须要一台一台国外服务器  
2. 申请机器人, 得到名字 和 token  
***直接在聊天搜索: botFather, 根据步骤申请机器人开通即可***  
3. 创建一个群， 根据机器人的名字邀请进来  
4. 启动http接口转发的代理  

``` 
$env:GOOS="windows|linux|darwin"  // windows|linux|darwin 选择对应系统进行打包
go build cmd/telegram/telegram.go
// 会生成 telegram<.exe> 二进制文件, 拷贝至服务器, 通过下面的命令启动代理， 当然最好挂载 scs 上， 而不是直接这样启动
./telegram<.exe> -t <token> -l :8080

```
5. 配置还少了一个`to` 这里需要添加chat_id, 我们在群组里面@机器人 发送任何信息， 然后再国外服务器执行下面请求
```
curl https://api.telegram.org/bot<token>/sendMessage
{
    "ok": true,
    "result": [
        {
            "update_id": 16881432,
            "message": {
                ...
                "chat": {
                    "id": -122,
                    "title": "报警",
                    "type": "group",
                    "all_members_are_administrators": true
                },
                "date": 1647066975,
                "text": "xxxx @机器人",
                ...
            }
        },
    ]
}
```
注意`chat.title` 是群组名 `chat.id` 就是接收者的 `chat_id`



## 企业微信

跟telegram一样，需要配置机器人， 这个设置机器人就比较简单了 参考文档  
https://open.work.weixin.qq.com/help2/pc/14931?person_id=1&is_tencent=  


## 自定义回调接口(POST请求)

虽然支持的有好多种， 但是基本都是很少用到， 这时候可以自己写回调接口来进行处理

```
callback:
    # 接受请求的url
    urls:
      - http://192.168.0.112:8080
    headers:
      Content-Type:
        - application/json
```
