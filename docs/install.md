# 安装

## 自动安装
linux(需要 git tar命令， 关闭selinux),mac, windows 请按照上面自行编译安装
```
/bin/bash -c "$(curl -fsSL http://download.hyahm.com/scs.sh)"
```

## 手动安装 


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

> 编译二进制文件（go>=1.13）
```
 git clone https://github.com/hyahm/scs.git
 cd scs
 go env -w GOPROXY=https://goproxy.cn,direct # 国外机器不需要这个
 go build -o scsd cmd/scs/main.go
 go build -o /usr/local/bin/scsctl cmd/scsctl/main.go
 cp default.yaml scs.yaml
 ./scsd
 ```

