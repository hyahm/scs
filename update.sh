#!/bin/bash
git pull
export GOPROXY=https://goproxy.cn
go build -o scsd cmd/scs/main.go
go build -o /usr/local/bin/scsctl cmd/scsctl/main.go