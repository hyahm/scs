//go:build windows
// +build windows

/*
 * @Author: your name
 * @Date: 2021-04-25 19:08:58
 * @LastEditTime: 2021-04-25 19:32:56
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /scs/public/rlimit.go
 */

package internal

import "fmt"

func setrlimit() {
	fmt.Println("windows can not set limit")
}
