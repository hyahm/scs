package main

import (
	"fmt"

	"github.com/hyahm/scs/pkg"
)

func main() {
	// 随机生成 30 - 50之间的随机数
	fmt.Print(pkg.RandomToken())
}
