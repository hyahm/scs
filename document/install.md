[安装](install.md)  
[客户端scsctl使用](scsctl.md)  
[报警](alert.md)  
[服务添加删除接口](script.md)  
[硬件监控配置说明](hardware.md)

# Install(linux 为例，其他系统参考即可)

暂时没打包成二进制，需要自己编译
### 依赖
- git  敲下 git 如果有如初则安装完git
- go >= 1.12
从中文网下载对应系统的包，解压后将 go/bin 添加到环境变量， 敲下 go 如果有如初则安装完go  

### 下载源码
```
mkdir /data
cd /data
git clone https://github.com/hyahm/scs.git

```
### 打包成二进制文件
```
cd scs
export GOPROXY=https://goproxy.cn  # 国内需要加个代理
go build -o scsd cmd/scs/main.go  # 服务器端
go build -o /usr/local/bin/scsctl cmd/scsctl/main.go  # 服务器端
```
### 启动服务并自启
```
cp default.yaml /etc/scs.yaml  # 拷贝配置文件
cp systemd/scs.service /etc/systemd/system/scsd.service   # 拷贝启动脚本
systemctl start scsd  # 启动服务
systemctl enable scsd   # 开机自启
```

### 验证
```
# 等待1秒后， 查看状态，应该可以看到如下信息， 下面是一条测试的信息， 仅供参考  
scsctl status
<node: local, url: https://127.0.0.1:11111>
--------------------------------------------------
PName     Name      Status     Ppid      UpTime    Verion   CanNotStop     Failed    Command
test        test_0      Running    1469      3m49s               true           0         cd /data/scs && python3 test.py
```