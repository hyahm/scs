## 准备一个报警器(根据示例来认识配置)

**创建一个自定义的报警器，下面是一个go语言的例子，也可以用其他语言写**
本地只要有go >= 1.11 即可， 创建一个 `main.go`文件
```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func callback(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("read body failed"))
		return
	}
	fmt.Println(string(b))
	w.Write([]byte("ok"))
}

func main() {
	http.HandleFunc("/", callback)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
```

直接 `go run main.go` 启动服务

> 测试POST请求
```
PS E:\code\scs> (curl -Method POST http://localhost:10000 -UseBasicParsing).Content
ok
// linux 使用这个就好了
// curl -X POST http://localhost:10000
// ok
```

> 将报警器添加到配置文件

```
alert:
  callback:
    # 接受请求的url
    urls:
      - http://localhost:10000
    headers:
      Content-Type:
        - application/json
```

> 让配置文件生效
```
scsctl config reload
```

## 先测试下服务崩溃的报警

> 打开配置文件，先将我们写的报警器配置上去，然后我们将脚本代码故意写错
```
scripts:
  - name: test
    # 将命令的第一个$去掉
    command: n=1;while($n -eq 1){ Get-Date;Start-Sleep -s 10}
```
> 加载配置
```
scsctl config reload
```
> 这时候配置已经生效， 但是不会影响运行的程序， 所以需要重启一下这个服务, 顺便又学习了一个客户端命令

```powershell
PS E:\code\scs> scsctl.exe restart test
{"code": 200, "msg": "waiting restart"}
```

> 回到`go run main.go`的命令行, 多出来下面的信息请求
```
PS C:\Users\cande\Desktop> go run .\main.go
{"Title":"service error stop","HostName":"DESKTOP-6SSBCHN","Pname":"test","Name":"test_0","BrokenTime":"2022-03-13 11:54:33.5085437 +0800 CST m=+5629.547435901","Reason":"exit status 1"}
```

## cpu报警

**为了避免报警信息错乱， 我们先把上面的错误脚本恢复（你已经是个老手了，步骤省略）**
> 默认的90通过特殊手段是可以超过的， 但是我们直接用最简单的, 直接改成1(即使用率超过1%就报警)
```yaml
cpu: 1
```
> 重载配置

> 不出意外的话 `go run main.go` 命令行终端那边的日志会多出来一条
```
{"Title":"cpu使用率超过 1.00%","HostName":"DESKTOP-6SSBCHN","UsePercent":32.19,"BrokenTime":"2022-03-13 12:52:06.5708559 +0800 CST m=+11.363103301","Top":"Process: Code.exe, Pid: 10520, Percent: 5.45%"}
```

**Title字段已经提示了报警的cpu报警设置的警报线，Top字段会显示使用率最高的进程**




## 磁盘报警

**为了避免报警信息错乱， 我们先把上面的警告恢复（你已经是个老手了，步骤省略）**
> 我们直接用最简单的, 直接改成1(即使用率超过1%就报警)， 注意的是，这是任意磁盘超过这个都会报警
```yaml
disk: 1
```
> 重载配置

> 不出意外的话 `go run main.go` 命令行终端那边的日志会多出来一条
```
{"Title":"硬盘使用率超过 1.00%, 当前使用率 16.16%","HostName":"DESKTOP-6SSBCHN","DiskPath":"C:","Use":47,"Total":292,"BrokenTime":"2022-03-13 12:59:05.5734208 +0800 CST m=+415.877292301"}
```

**注意看DiskPath，提示我c盘超过了，实际上我的D盘也是超过， 但是因为同是报错的报错，只会报警一个磁盘， 当把这个盘空间恢复后后面会报警提示D盘也超过了**


## 内存报警

> 我们直接用最简单的, 直接改成1(即使用率超过1%就报警) 

```yaml
mem: 1
``` 

> 重载配置

> 不出意外的话 `go run main.go` 命令行终端那边的日志会多出来一条
```
{"Title":"内存繁忙超过1.00%","HostName":"DESKTOP-6SSBCHN","Use":8,"Total":15,"BrokenTime":"2022-03-13 13:04:48.011744 +0800 CST m=+758.315615501","Top":"Process: Code.exe, Pid: 876, Percent: 5.74%"}
```

**Top字段会显示使用率最高的进程**


## scs服务报警

使用过程中，大家发现不管是系统内部还是挂载的脚本 都做到了全程自动监控  
聪明绝顶的你会想到另外一个问题：如果scs自身意外退出， 或者机器关闭了，那我该怎么办
其实里面对外提供了一个http的空接口， 通过200状态码来测试scs本身是否在线
为了安全，内部都是有token验证的， 也可以直接增加ip白名单来检测状态

- 为了测试，我准备一台linux虚拟机，同样安装了scs， 用我本地的scs服务来监控虚拟机的scs服务
为了简单， linux虚拟机 采用`一键脚本脚本`方式安装scs，启用的https. 配置文件如下
```
[root@localhost ~]# cat /etc/scs.yaml
# 监听端口
listen: :11111
# 服务日志配置
log:
    path: log
    day: true
    size: 0
# 请求头认证 Token： xxxx, `一键脚本脚本`方式安装scs 会随机生成一个
token: 'gSstfcoK&3,5LQbOCR|8k+8Ftp#TJ}'
...
```

- linux ip信息如下: 192.168.101.12

- windows 本地的ip： 192.168.101.11

- 这个时候有另外2项配置可以使用了 `monitor` ,`monitored`

  monitor(array): 就是监控机， 此次案例是我本地windows监控虚拟机的linux 所以我本地要配置 monitor  
  monitored(array): 就是被监控机， linux要配置 monitored  
- 填写配置， 因为被监控机器的协议可以是http也可以是https, 所以我们要填全,windows的配置如下添加  
	```
	monitor:
	  - https://192.168.101.12:11111
- 重载本地windows scs配置, 先不做下一步操作， 去报警器那边查看有如下日志， 这是因为linux服务器那边没有设置白名单
  ```powershell
  {"Title":"服务器或scs服务出现问题: https://192.168.101.12:11111","HostName":"DESKTOP-6SSBCHN","BrokenTime":"2022-03-13 13:34:47.4788516 +0800 CST m=+2557.782723101"}
  ```
- 然后我们需要去linux添加白名单windows ip
    ```
  	monitored:
      - 192.168.101.11
	```
- 重载一下linux的配置文件

     ```powershell
	{"Title":"服务器或scs服务恢复: https://192.168.101.12:11111","HostName":"DESKTOP-6SSBCHN","BrokenTime":"2022-03-13 13:52:10.1485054 +0800 CST m=+33.116418001","FixTime":"2022-03-13 14:02:51.2865763 +0800 CST"}
    ```

