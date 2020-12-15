// +build !windows

package public

import (
	"fmt"
	"os"
	"syscall"
)

func init() {
	var rlim syscall.Rlimit
	// s := []uint{
	// 	RLIMIT_AS,
	// 	RLIMIT_CORE,
	// 	RLIMIT_CPU,
	// 	RLIMIT_DATA,
	// 	RLIMIT_FSIZE,
	// 	RLIMIT_NOFILE,
	// 	RLIMIT_STACK,
	// }
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim)
	if err != nil {
		fmt.Println("get rlimit error: " + err.Error())
		os.Exit(1)
	}
	rlim.Cur = 6553500
	rlim.Max = 6553500
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlim)
	if err != nil {
		fmt.Println("set rlimit error: " + err.Error())
	}
}
