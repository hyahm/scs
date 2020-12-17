# /bin/bash
workdir=/data
sudo mkdir /opt
source ~/.bashrc
gobin=$(which go)
if [[ $? -ne 0 ]]; then
    sudo curl -k -s -L -o /opt/go1.15.2.linux-amd64.tar.gz https://studygolang.com/dl/golang/go1.15.2.linux-amd64.tar.gz
    cd /opt/
        tar -xf go1.15.2.linux-amd64.tar.gz
    echo "export PATH=$PATH:/opt/go/bin" >> ~/.bashrc
    source ~/.bashrc
    gobin=/opt/go/bin/go
fi
sudo mkdir $workdir
sudo chown -R $USER:$USER $workdir
cd $workdir
a=$(ping gitee.com -w 2 -i 3 | grep time= | awk -F'time=' '{print $2}' | awk -F'.' '{print $1}' | awk '{print $1}')
b=$(ping github.com -w 2 -i 3 | grep time= | awk -F'time=' '{print $2}' |  awk -F'.' '{print $1}' | awk '{print $1}')
if [[ ${a:-1000000} -lt ${b:-10000} ]];then
         sudo $gobin env -w GOPROXY=https://goproxy.cn
         $gobin env -w GOPROXY=https://goproxy.cn
         echo "git clone from https://gitee.com/cander/scs.git"
     git clone https://gitee.com/cander/scs.git
     if [[ $? -ne 0 ]]; then
            exit 1
     fi
else
        echo "git clone from https://gitee.com/cander/scs.git"
        git clone https://github.com/hyahm/scs.git
    if [[ $? -ne 0 ]]; then
            exit 1
    fi
fi
cd scs
export GOPROXY=https://goproxy.cn
$gobin build -o scsd cmd/scs/main.go
$gobin build -o /usr/local/bin/scsctl cmd/scsctl/main.go
if [[ $? -ne 0 ]]; then
        echo "build scsctl failed , you can run as root by yourself
        source ~/.bashrc
        cd /data/scs
        $gobin build -o /usr/local/bin/scsctl cmd/scsctl/main.go
        "
fi
token=$($gobin run cmd/random/random.go)

if [[ ! -f /etc/scs.yaml ]];then
sudo /bin/bash -c 'export TOKEN=$token && cat > /etc/scs.yaml << EOF
# 监听端口
listen: :11111
# 脚本保存控制台输出日志的行数
logCount: 500
# 服务日志配置
log:
    path: log
    day: true
    size: 0
# 请求头认证 Token： xxxx
token: "$TOKEN"
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
    server: https://telegram.hyahm.com
    username:  test
    password: 123456
    to:
      - -870968689
# 本地磁盘， cpu， 内存监控项， 确保elert存在才会有通知
probe:
  # mem使用率, 默认90
  mem: 60
  # cpu使用率, 默认90
  cpu: 90
  # 硬盘使用率， 默认85
  disk: 80
  # 排除的挂载点， 默认已经去掉了swap， 设备, 数组
  excludeDisk:
  # 检测间隔， 默认10秒
  interval: 10s
  # 下次报警时间间隔， 如果恢复了就重置
  continuityInterval: 1h
scripts:
  # - name: ls
  #   # 脚本执行的根目录
  #   dir: D:\\myproject\\scs
  #   env:
  #     GOPROXY: MMMM
  #   # 启动脚本
  #   command: "go"
  #   # 不写默认10分钟
  #   continuityInterval: 1h
  #   always: true
  #   # 禁用报警， 默认启动
  #   # disableAlert: true
  #   # replicate， 开启副本数
  #   replicate: 1
  #   killTime: 2s
  #   alert:
  #     email:
  #       - 727023885460@qq.com
  #
  #     rocket:
  #       - ""
EOF'
fi

sudo /bin/bash -c "cat > /etc/systemd/system/scs.service <<EOF
[Unit]
Description=Scs Service Control Script
After=network.target
After=network-online.target
Wants=network-online.target
[Service]
LimitNOFILE=6553500
LimitNPROC=6553500
WorkingDirectory=$workdir/scs
ExecStart=$workdir/scs/scsd -f /etc/scs.yaml
ExecStop=/bin/kill -s QUIT \$MAINPID
Type=simple
[Install]
WantedBy=multi-user.target
EOF"
sudo setenforce 0
systemctl daemon-reload
sudo systemctl start scs
cat > ~/scsctl.yaml <<EOF
nodes:
  localhost:
    url: https://127.0.0.1:11111
    token: "$token"
EOF