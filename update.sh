#!/bin/bash
git pull
export GOPROXY=https://goproxy.cn
go build -o scsd cmd/scsd/scsd.go
go build -o /usr/local/bin/scsctl cmd/scsctl/scsctl.go