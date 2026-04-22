#!/bin/sh
set -e  # 出错即退出

# --- 配置区 ---
WORKDIR="/data"
GO_VERSION="1.17.6"
SCS_CONF="/etc/scs.yaml"
INSTALL_OPT="/opt"
# 确保安装目录存在
sudo mkdir -p $INSTALL_OPT

# --- 1. Go 环境检查与安装 ---
# 优先检查现有环境
if ! command -v go &> /dev/null; then
    echo "未检测到 Go，准备安装版本: ${GO_VERSION}..."
    GO_TAR="go${GO_VERSION}.linux-amd64.tar.gz"

    sudo curl -k -s -L -o "${INSTALL_OPT}/${GO_TAR}" "https://dl.google.com/go/${GO_TAR}"
    sudo tar -C $INSTALL_OPT -xf "${INSTALL_OPT}/${GO_TAR}"

    # 写入环境变量并立即对当前进程生效
    echo "export PATH=\$PATH:${INSTALL_OPT}/go/bin" >> ~/.bashrc
    export PATH=$PATH:${INSTALL_OPT}/go/bin
fi

GOBIN=$(which go)
echo "使用 Go 路径: $GOBIN"

# --- 2. 源码下载逻辑 (智能测速) ---
sudo mkdir -p "$WORKDIR"
sudo chown -R "$USER:$USER" "$WORKDIR"
cd "$WORKDIR"

if [[ ! -d "scs" ]]; then
    echo "正在测试 Gitee 与 Github 延迟..."
    # 简化测速逻辑：取平均延迟数值
    LATENCY_GITEE=$(ping -c 2 -q gitee.com | awk -F'/' 'END {print ($5 ? $5 : 9999)}' | cut -d. -f1)
    LATENCY_GITHUB=$(ping -c 2 -q github.com | awk -F'/' 'END {print ($5 ? $5 : 9999)}' | cut -d. -f1)

    if [ "$LATENCY_GITEE" -lt "$LATENCY_GITHUB" ]; then
        echo "Gitee 较快，开始克隆..."
        git clone https://gitee.com/cander/scs.git
        $GOBIN env -w GOPROXY=https://goproxy.cn,direct
    else
        echo "Github 较快或延迟相当，开始克隆..."
        git clone https://github.com/hyahm/scs.git
    fi
    cd scs
else
    echo "目录已存在，执行更新..."
    cd scs && git pull
fi

# --- 3. 编译阶段 ---
export GOPROXY=https://goproxy.cn,direct
echo "正在编译 scsd..."
$GOBIN build -o scsd cmd/scsd/scsd.go

echo "正在编译 scsctl..."
# 编译并移动到系统路径，通常需要 sudo
$GOBIN build -o scsctl cmd/scsctl/scsctl.go
sudo mv scsctl /usr/local/bin/

# 生成 Token
TOKEN=$($GOBIN run cmd/random/random.go)

# --- 4. 配置文件生成 ---
if [[ ! -f "$SCS_CONF" ]]; then
    echo "生成配置文件: $SCS_CONF"
    sudo bash -c "cat > $SCS_CONF <<EOF
listen: :11111
log:
    path: $WORKDIR/scs/log
    day: true
    size: 0
token: '$TOKEN'
probe:
    mem: 60
    cpu: 90
    disk: 80
    interval: 10s
    continuityInterval: 1h
EOF"

    cat > ~/.scsctl.yaml <<EOF
nodes:
  localhost:
    url: http://127.0.0.1:11111
    token: '$TOKEN'
EOF
fi

# --- 5. Systemd 服务配置 ---
echo "配置 Systemd 服务..."
sudo bash -c "cat > /etc/systemd/system/scsd.service <<EOF
[Unit]
Description=Scs Service Control Script
After=network.target network-online.target
Wants=network-online.target

[Service]
Type=simple
User=$USER
LimitNOFILE=6553500
LimitNPROC=6553500
WorkingDirectory=$WORKDIR/scs
ExecStart=$WORKDIR/scs/scsd -f $SCS_CONF
ExecStop=/bin/kill -s QUIT \$MAINPID
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF"

# --- 6. 启动服务 ---
# 针对 Fedora/AlmaLinux 禁用 SELinux (临时)
if command -v setenforce &> /dev/null; then
    sudo setenforce 0 || true
fi

sudo systemctl daemon-reload
sudo systemctl enable --now scsd

echo "部署完成！服务状态："
sudo systemctl status scsd --no-pager