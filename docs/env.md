# 环境变量

`env`

挂载上的服务每个里面必定存在下面4个环境变量， 请避免使用到

PNAME // 等于配置文件name的值  
NAME  // 由系统自动生成, 副本名  
TOKEN  // 从配置文件中读取  
PORT  // 等于配置文件port的值, 环境变量中需要PORT 也必须在port 设置， 环境变量中设置的PORT 无效  
PROJECT_HOME  // 项目根目录  
OS  // scs系统变量  

**环境变量可以在很多地方使用go语言的 `text/template` 模块渲染使用, 详细说明参考配置章节说明**


> 我们添加一个环境变量

```yaml
scripts:
  # 这个是脚本的名字， 名字必须是字母数字或下划线组成，以后都是根据此名字来操作
  - name: test
    always: true
    env:
      TEST_ENV: test
    # 这个是执行的基础命令
    command: $n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
```
> 重载配置文件`scsctl config reload`
> 查看环境变量
```
PS E:\code\scs> scsctl.exe env test_0   
...
TEST_ENV: test
...
```

> 以`SCS_TPL` 为前缀的变量支持 模板渲染
```
# 类似这样， windows 和linux 生成的后缀名可能不一样，用此环境变量来获取对应系统的执行文件
 SCS_TPL_OUT: '{{ if eq .OS "windows" }}main.exe{{ else }} main{{ end}}'
```
