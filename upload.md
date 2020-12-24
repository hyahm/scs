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
