@echo off

go build -o bin/scsd-3.7.1.exe  ./cmd/scsd/scsd.go;
go build -o bin/scsctl-3.7.1.exe  ./cmd/scsctl/scsctl.go;
$env:GOOS="linux";
go build -o bin/scsd_linux-3.7.1  ./cmd/scsd/scsd.go;
go build -o bin/scsctl_linux-3.7.1  ./cmd/scsctl/scsctl.go;
$env:GOOS="darwin";
go build -o bin/scsd_darwin-3.7.1  ./cmd/scsd/scsd.go;
go build -o bin/scsctl_darwin-3.7.1  ./cmd/scsctl/scsctl.go;
$env:GOOS="windows";