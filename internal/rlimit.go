package internal

import "runtime"

func Setrlimit() {
	if runtime.GOOS != "windows" {
		setrlimit()
	}
}
