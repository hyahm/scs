# 安装

## 自动安装
linux(需要 git tar命令， 关闭selinux),mac, windows 请按照上面自行编译安装
```
/bin/bash -c "$(curl -fsSL http://download.hyahm.com/scs.sh)"
```

## 手动安装  

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

