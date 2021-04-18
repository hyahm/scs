# scs

service control service or script

# 主要功能
1.  为了保护数据不会在执行中被人工手动中断丢失， 可以让服务在某段时间内才能停止
2.  监控硬件信息， 主要是磁盘， cpu， 内存,  scs服务  
3.  服务之间可以相互控制增删改查  
4.  报警功能api  
5.  支持定时器功能执行命令或脚本
6.  客户端控制多台服务器
7.  通过配置文件安装服务
8.  可以将一些执行耗时的脚本托管给scs处理， 以便快速返回结果  


# 服务控制脚本
类似supervisor,但是更高级，支持所有系统 
自带监控及通知   
服务控制脚本能否停止 最大程度防止脚本数据丢失  
码云地址: https://gitee.com/cander/scs
# 功能
- [x] 跨平台  
- [x] 支持启动所有脚本  
- [x] 自带硬件检测  
- [x] 一键可启动(linux)  
- [x] 支持的报警渠道- 邮箱, rocket.chat报警, telegram, 企业微信

[版本更新](update.md)






# 客户端
1.1.0 版开始， 配置文件必须要配置值， 不然什么也不会出来
```yaml
nodes:
  localhost: 
    url: https://127.0.0.1:11111
    token:
group:
  local:
    - localhost
```

```

scsctl status 
scsctl status pname
scsctl status pname name
scsctl start 
scsctl start pname
scsctl start pname name
scsctl restart --all
scsctl restart pname 
scsctl restart pname name
scsctl kill --all
scsctl kill pname 
scsctl kill pname name
scsctl stop --all
scsctl stop pname 
scsctl stop pname name
scsctl update --all
scsctl update pname 
scsctl update pname name
scsctl remove --all
scsctl remove pname 
scsctl remove pname name
scsctl enable pname
scsctl disable pname
scsctl log  name[:update|log|lookPath] # 不区分大小写
# 加载配置文件
scsctl config reload
```

