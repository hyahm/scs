package main

import (
	"os"
)

func main() {
	_, err := os.Stat("/data/scs/cmd/proxy/telegram.go")
	if os.IsNotExist(err) {

	}
}
