# v3.7.1
- 取消内部token
- 修复停止的服务pid不为0的问题
- status默认显示少了failed, cpu, mem, command, 如果不要显示加上 -v
- 修复replicate为0的时候错乱


# v3.7.0
- 增加脚本权限控制
- 修复查看日志卡住的问题


# v3.6.9
- 支持代理， 需要设置获取真实ip的请求头  
```
proxyHeader: X-Forwarded-For
``` 


# v3.6.8
- 修复remove的时候，脚本异常
- 修复删除script后，再次添加失败的问题
- 移除update unstable标签


# v3.6.7
- 修复disable后 start 可以继续启动的问题
- 移除客户端disable列
- disable的脚本不会显示在 status中


# v3.6.6
- 修复重载always的问题

# v3.6.5
- 修复重载的时候脚本控制丢失的问题
- 客户端状态显示服务器版本
- command 显示cd优化
- 增加脚本错误但是没报警器的日志

# v3.6.3(2022-3-19)
- `scsctl log <name>` 的时候，如果没有生成日志文件就会直接返回  
- 修复`scsctl status <pname> <name>`
- 修复crontab 修改时间无法生效的问题

# v3.6.2 beta(2022-3-17)
- 代码做了大量优化， 移除`scsctl remove --all`


# v3.6.0 (2022-3-13)
- 优化日志和log命令
- `cron` 增加`times` 来指定循环次数

# v3.2.0(2020-12-24)
- 里面的字段变量全部用go template 写法， 参考： https://cloud.tencent.com/developer/section/1145004
- 修改字段loopPath 为 preStart
- 

# v2.3.0(2020-12-24)
```
取消loop， 更改为cron, 基本满足所有crontab的要求
```

# v2.3.0(2020-12-24)
```
代码优化
```

# v2.2.2(2020-12-23)
```
新增loop 和 disable 选项
loop: 类似定时器， 多少秒执行一次
disable: 是否禁用脚本
```

# v1.2.4(2020-10-20)
```
移除command $PNAME $NAME $PORT的替换， 并添加到环境变量中， PNAME NAME PORT
新增
scsctl search xxxx
scsctl install xxxx
的基本结构实现
```

# v1.1.5(2020-10-01)
帮助显示错误的问题

# v1.1.4(2020-09-30)
修复单执行status成start


# v1.1.3(2020-09-30)
客户端增加超时配置， 默认3秒
```yaml
readTimeout
```

# v1.1.2(2020-09-30)
修复status start pname name无效

# v1.1.1(2020-09-30)
修复单节点操作代码异常

# v1.1.0(2020-09-30)
- 客户端新增集群方案（新增flag  -n -g）

- 客户端配置文件修改, 默认为家目录的scs.yaml 文件,下面为配置实例

```vim

所有节点
nodes:
  # 节点名关联参数(-n)
  me: 
    url: http://192.168.10.10:11111
    token: "al3455555555j0(&^(*jha67"
  novel:
    url: http://127.0.0.1:11111
    token: "o6666666666666664"
# 节点分类组
group:
  # 组名关联参数(-g)
  aa: 
    # 下面为关联的节点
    - me
    - novel
  bb: 
    - novel
```
- 显示配置文件的所有节点名
```
scsctl config show
```
