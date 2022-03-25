//go:build !windows
// +build !windows

/*
 * @Author: your name
 * @Date: 2021-04-25 19:08:58
 * @LastEditTime: 2021-04-25 19:32:56
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /scs/public/rlimit.go
 */

package internal

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
)

func setrlimit() {
	var rlim syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim)
	if err != nil {
		fmt.Println("get rlimit error: " + err.Error())
		os.Exit(1)
	}
	var MaxRlimt uint64 = 1048576
	if runtime.GOOS == "darwin" {
		MaxRlimt = 10240
	}

	rlim.Cur = MaxRlimt
	rlim.Max = MaxRlimt
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlim)
	if err != nil {
		fmt.Println("set rlimit error: " + err.Error())
	}
}
